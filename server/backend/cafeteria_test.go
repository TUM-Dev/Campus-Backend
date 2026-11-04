package backend

import (
	"context"
	"crypto/md5"
	"database/sql"
	"fmt"
	"os"
	"path"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	pb "github.com/TUM-Dev/Campus-Backend/server/api/tumdev"
	"github.com/TUM-Dev/Campus-Backend/server/backend/cron"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

const testImage = "./test_data/sampleimage.jpeg"

type CafeteriaSuite struct {
	suite.Suite
	DB             *gorm.DB
	mock           sqlmock.Sqlmock
	deviceBuf      *deviceBuffer
	prevStorageDir string
}

// SetupTest creates a fresh mock per test so that expectations of one test cannot leak into another
// and redirects image storage into a temporary directory.
func (s *CafeteriaSuite) SetupTest() {
	var (
		db  *sql.DB
		err error
	)

	db, s.mock, err = sqlmock.New()
	require.NoError(s.T(), err)

	dialector := mysql.New(mysql.Config{
		Conn:       db,
		DriverName: "mysql",
	})
	s.mock.ExpectQuery("SELECT VERSION()").
		WillReturnRows(sqlmock.NewRows([]string{"VERSION()"}).AddRow("8.4.0"))
	s.DB, err = gorm.Open(dialector, &gorm.Config{})
	require.NoError(s.T(), err)

	s.deviceBuf = newDeviceBuffer()

	s.prevStorageDir = cron.StorageDir
	cron.StorageDir = s.T().TempDir()
}

func (s *CafeteriaSuite) TearDownTest() {
	cron.StorageDir = s.prevStorageDir
}

func (s *CafeteriaSuite) AfterTest(_, _ string) {
	require.NoError(s.T(), s.mock.ExpectationsWereMet())
}

func (s *CafeteriaSuite) server() CampusServer {
	return CampusServer{db: s.DB, deviceBuf: s.deviceBuf}
}

func testCtx() context.Context {
	return metadata.NewIncomingContext(context.Background(), metadata.MD{})
}

// expectedImagePath mirrors the path imageWrapper/storeImage generate for an uploaded image
func expectedImagePath(img []byte, subDir string, id int64) string {
	return fmt.Sprintf("%s/%x.jpeg", path.Join(cron.StorageDir, "rating", subDir, fmt.Sprint(id)), md5.Sum(img))
}

func readTestImage(t *testing.T) []byte {
	img, err := os.ReadFile(testImage)
	require.NoError(t, err)
	return img
}

const (
	expectedHeadCountQuery    = "SELECT * FROM `canteen_head_counts` WHERE `canteen_head_counts`.`canteen_id` = ? ORDER BY `canteen_head_counts`.`canteen_id` LIMIT ?"
	expectedCanteenByName     = "SELECT * FROM `canteens` WHERE name LIKE ? ORDER BY `canteens`.`cafeteria` LIMIT ?"
	expectedInsertRating      = "INSERT INTO `canteen_ratings` (`points`,`comment`,`cafeteriaID`,`image`,`timestamp`) VALUES (?,?,?,?,?)"
	expectedCanteenTagExists  = "SELECT count(*) FROM `canteen_rating_tag_options` WHERE cafeteriaRatingTagOption LIKE ?"
	expectedInsertCanteenTag  = "INSERT INTO `canteen_rating_tags` (`correspondingRating`,`tagID`,`points`) VALUES (?,?,?)"
	expectedDishByName        = "SELECT * FROM `dishes` WHERE name LIKE ? AND cafeteriaID = ? ORDER BY `dishes`.`dish` LIMIT ?"
	expectedInsertDishRating  = "INSERT INTO `dish_ratings` (`points`,`dishID`,`comment`,`image`,`timestamp`) VALUES (?,?,?,?,?)"
	expectedDishNameTagIDs    = "SELECT `nameTagID` FROM `dish_to_dish_name_tags` WHERE dishID = ?"
	expectedInsertDishNameTag = "INSERT INTO `dish_name_tags` (`correspondingRating`,`tagNameID`,`points`) VALUES (?,?,?)"
	expectedDishTagExists     = "SELECT count(*) FROM `dish_rating_tag_options` WHERE dishRatingTagOption LIKE ?"
	expectedInsertDishTag     = "INSERT INTO `dish_rating_tags` (`parentRating`,`tagID`,`points`) VALUES (?,?,?)"
)

func (s *CafeteriaSuite) expectCanteenLookup(name string, id int64) {
	s.mock.ExpectQuery(regexp.QuoteMeta(expectedCanteenByName)).
		WithArgs(name, 1).
		WillReturnRows(sqlmock.NewRows([]string{"cafeteria", "name", "address", "latitude", "longitude"}).
			AddRow(id, name, "Boltzmannstraße 19, 85748 Garching", 48.2681, 11.6723))
}

func (s *CafeteriaSuite) Test_GetCanteenHeadCount() {
	ts := time.Date(2023, 10, 9, 11, 45, 22, 0, time.UTC)
	s.mock.ExpectQuery(regexp.QuoteMeta(expectedHeadCountQuery)).
		WithArgs("mensa-garching", 1).
		WillReturnRows(sqlmock.NewRows([]string{"canteen_id", "count", "max_count", "percent", "timestamp"}).
			AddRow("mensa-garching", 250, 1000, 25, ts))

	server := s.server()
	response, err := server.GetCanteenHeadCount(testCtx(), &pb.GetCanteenHeadCountRequest{CanteenId: "mensa-garching"})
	require.NoError(s.T(), err)
	require.Equal(s.T(), &pb.GetCanteenHeadCountReply{
		Count:     250,
		MaxCount:  1000,
		Percent:   25,
		Timestamp: timestamppb.New(ts),
	}, response)
}

func (s *CafeteriaSuite) Test_GetCanteenHeadCountNotFound() {
	s.mock.ExpectQuery(regexp.QuoteMeta(expectedHeadCountQuery)).
		WithArgs("mensa-unknown", 1).
		WillReturnRows(sqlmock.NewRows([]string{"canteen_id", "count", "max_count", "percent", "timestamp"}))

	server := s.server()
	response, err := server.GetCanteenHeadCount(testCtx(), &pb.GetCanteenHeadCountRequest{CanteenId: "mensa-unknown"})
	require.NoError(s.T(), err)
	require.Equal(s.T(), uint32(0), response.Count)
	require.Equal(s.T(), uint32(0), response.MaxCount)
	require.Equal(s.T(), float32(-1), response.Percent)
}

func (s *CafeteriaSuite) Test_CreateCanteenRating() {
	const (
		canteenName = "mensa-garching"
		comment     = "Everything perfect, 2 Stars"
		cafeteriaID = int64(2)
		ratingID    = int64(42)
		points      = int32(2)
	)
	img := readTestImage(s.T())
	imgPath := expectedImagePath(img, "cafeterias", cafeteriaID)

	s.expectCanteenLookup(canteenName, cafeteriaID)
	s.mock.ExpectBegin()
	s.mock.ExpectExec(regexp.QuoteMeta(expectedInsertRating)).
		WithArgs(points, comment, cafeteriaID, imgPath, sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(ratingID, 1))
	s.mock.ExpectCommit()
	for _, tag := range []struct{ id, points int64 }{{5, 1}, {7, 3}} {
		s.mock.ExpectQuery(regexp.QuoteMeta(expectedCanteenTagExists)).
			WithArgs(tag.id).
			WillReturnRows(sqlmock.NewRows([]string{"count(*)"}).AddRow(1))
		s.mock.ExpectBegin()
		s.mock.ExpectExec(regexp.QuoteMeta(expectedInsertCanteenTag)).
			WithArgs(ratingID, tag.id, tag.points).
			WillReturnResult(sqlmock.NewResult(1, 1))
		s.mock.ExpectCommit()
	}

	server := s.server()
	response, err := server.CreateCanteenRating(testCtx(), &pb.CreateCanteenRatingRequest{
		Points:    points,
		CanteenId: canteenName,
		Comment:   comment,
		Image:     img,
		RatingTags: []*pb.RatingTag{
			{TagId: 5, Points: 1},
			{TagId: 7, Points: 3},
		},
	})
	require.NoError(s.T(), err)
	require.Equal(s.T(), &pb.CreateCanteenRatingReply{}, response)
	require.FileExists(s.T(), imgPath)
}

func (s *CafeteriaSuite) Test_CreateCanteenRatingUnknownTag() {
	const (
		canteenName = "mensa-garching"
		cafeteriaID = int64(2)
		ratingID    = int64(42)
	)
	s.expectCanteenLookup(canteenName, cafeteriaID)
	s.mock.ExpectBegin()
	s.mock.ExpectExec(regexp.QuoteMeta(expectedInsertRating)).
		WithArgs(int32(4), "", cafeteriaID, "", sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(ratingID, 1))
	s.mock.ExpectCommit()
	s.mock.ExpectQuery(regexp.QuoteMeta(expectedCanteenTagExists)).
		WithArgs(int64(99)).
		WillReturnRows(sqlmock.NewRows([]string{"count(*)"}).AddRow(0))

	server := s.server()
	_, err := server.CreateCanteenRating(testCtx(), &pb.CreateCanteenRatingRequest{
		Points:     4,
		CanteenId:  canteenName,
		RatingTags: []*pb.RatingTag{{TagId: 99, Points: 1}},
	})
	require.Equal(s.T(), codes.InvalidArgument, status.Code(err))
}

func (s *CafeteriaSuite) Test_CreateCanteenRatingInvalidInput() {
	tests := []struct {
		name string
		req  *pb.CreateCanteenRatingRequest
	}{
		{"points too high", &pb.CreateCanteenRatingRequest{Points: 6, CanteenId: "mensa-garching"}},
		{"points negative", &pb.CreateCanteenRatingRequest{Points: -1, CanteenId: "mensa-garching"}},
		{"comment too long", &pb.CreateCanteenRatingRequest{Points: 3, CanteenId: "mensa-garching", Comment: string(make([]byte, 257))}},
		{"comment contains @", &pb.CreateCanteenRatingRequest{Points: 3, CanteenId: "mensa-garching", Comment: "mail me@example.com"}},
	}
	server := s.server()
	for _, tt := range tests {
		s.Run(tt.name, func() {
			response, err := server.CreateCanteenRating(testCtx(), tt.req)
			require.Nil(s.T(), response)
			require.Equal(s.T(), codes.InvalidArgument, status.Code(err))
		})
	}
}

func (s *CafeteriaSuite) Test_CreateCanteenRatingUnknownCanteen() {
	s.mock.ExpectQuery(regexp.QuoteMeta(expectedCanteenByName)).
		WithArgs("mensa-unknown", 1).
		WillReturnRows(sqlmock.NewRows([]string{"cafeteria", "name", "address", "latitude", "longitude"}))

	server := s.server()
	response, err := server.CreateCanteenRating(testCtx(), &pb.CreateCanteenRatingRequest{Points: 3, CanteenId: "mensa-unknown"})
	require.Nil(s.T(), response)
	require.Equal(s.T(), status.Error(codes.InvalidArgument, "Canteen does not exist. Rating has not been saved."), err)
}

func (s *CafeteriaSuite) Test_CreateDishRating() {
	const (
		canteenName = "MENSA_GARCHING"
		dishName    = "Vegane rote Grütze mit Soja-Vanillesauce"
		comment     = "Very tasty"
		cafeteriaID = int64(2)
		dishID      = int64(13)
		ratingID    = int64(42)
		nameTagID   = int64(3)
		points      = int32(5)
	)
	img := readTestImage(s.T())
	imgPath := expectedImagePath(img, "dishes", dishID)

	s.expectCanteenLookup(canteenName, cafeteriaID)
	s.mock.ExpectQuery(regexp.QuoteMeta(expectedDishByName)).
		WithArgs(dishName, cafeteriaID, 1).
		WillReturnRows(sqlmock.NewRows([]string{"dish", "name", "type", "cafeteriaID"}).
			AddRow(dishID, dishName, "Dessert", cafeteriaID))
	s.mock.ExpectBegin()
	s.mock.ExpectExec(regexp.QuoteMeta(expectedInsertDishRating)).
		WithArgs(points, dishID, comment, imgPath, sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(ratingID, 1))
	s.mock.ExpectCommit()
	s.mock.ExpectQuery(regexp.QuoteMeta(expectedDishNameTagIDs)).
		WithArgs(dishID).
		WillReturnRows(sqlmock.NewRows([]string{"nameTagID"}).AddRow(nameTagID))
	s.mock.ExpectBegin()
	s.mock.ExpectExec(regexp.QuoteMeta(expectedInsertDishNameTag)).
		WithArgs(ratingID, nameTagID, points).
		WillReturnResult(sqlmock.NewResult(1, 1))
	s.mock.ExpectCommit()
	s.mock.ExpectQuery(regexp.QuoteMeta(expectedDishTagExists)).
		WithArgs(int64(1)).
		WillReturnRows(sqlmock.NewRows([]string{"count(*)"}).AddRow(1))
	s.mock.ExpectBegin()
	s.mock.ExpectExec(regexp.QuoteMeta(expectedInsertDishTag)).
		WithArgs(ratingID, int64(1), int32(4)).
		WillReturnResult(sqlmock.NewResult(1, 1))
	s.mock.ExpectCommit()

	server := s.server()
	response, err := server.CreateDishRating(testCtx(), &pb.CreateDishRatingRequest{
		Points:     points,
		CanteenId:  canteenName,
		Dish:       dishName,
		Comment:    comment,
		Image:      img,
		RatingTags: []*pb.RatingTag{{TagId: 1, Points: 4}},
	})
	require.NoError(s.T(), err)
	require.Equal(s.T(), &pb.CreateDishRatingReply{}, response)
	require.FileExists(s.T(), imgPath)
}

func (s *CafeteriaSuite) Test_CreateDishRatingUnknownDish() {
	const cafeteriaID = int64(2)
	s.expectCanteenLookup("MENSA_GARCHING", cafeteriaID)
	s.mock.ExpectQuery(regexp.QuoteMeta(expectedDishByName)).
		WithArgs("Unknown dish", cafeteriaID, 1).
		WillReturnRows(sqlmock.NewRows([]string{"dish", "name", "type", "cafeteriaID"}))

	server := s.server()
	response, err := server.CreateDishRating(testCtx(), &pb.CreateDishRatingRequest{
		Points:    3,
		CanteenId: "MENSA_GARCHING",
		Dish:      "Unknown dish",
	})
	require.Nil(s.T(), response)
	require.Equal(s.T(), codes.InvalidArgument, status.Code(err))
}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestCafeteriaSuite(t *testing.T) {
	suite.Run(t, new(CafeteriaSuite))
}
