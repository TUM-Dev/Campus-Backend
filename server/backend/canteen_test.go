package backend

import (
	"bytes"
	"context"
	"crypto/md5"
	"database/sql"
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"net/http"
	"net/http/httptest"
	"os"
	"path"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	pb "github.com/TUM-Dev/Campus-Backend/server/api/tumdev"
	"github.com/TUM-Dev/Campus-Backend/server/backend/cron"
	"github.com/TUM-Dev/Campus-Backend/server/model"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
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

type CanteenSuite struct {
	suite.Suite
	DB             *gorm.DB
	mock           sqlmock.Sqlmock
	deviceBuf      *deviceBuffer
	prevStorageDir string
}

// SetupTest creates a fresh mock per test so that expectations of one test cannot leak into another
// and redirects image storage into a temporary directory.
func (s *CanteenSuite) SetupTest() {
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
	s.T().Setenv(ratingImagesEnabledEnv, "true")
}

func (s *CanteenSuite) TearDownTest() {
	cron.StorageDir = s.prevStorageDir
}

func (s *CanteenSuite) AfterTest(_, _ string) {
	require.NoError(s.T(), s.mock.ExpectationsWereMet())
}

func (s *CanteenSuite) server() CampusServer {
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
	expectedCanteenByName     = "SELECT * FROM `canteens` WHERE name = ? ORDER BY `canteens`.`cafeteria` LIMIT ?"
	expectedInsertRating      = "INSERT INTO `canteen_ratings` (`points`,`comment`,`cafeteriaID`,`image`,`device_hash`,`status`,`updated_at`,`timestamp`) VALUES (?,?,?,?,?,?,?,?)"
	expectedCanteenTagsIn     = "SELECT `cafeteriaRatingTagOption` FROM `canteen_rating_tag_options` WHERE cafeteriaRatingTagOption IN (?,?)"
	expectedCanteenTagSingle  = "SELECT `cafeteriaRatingTagOption` FROM `canteen_rating_tag_options` WHERE cafeteriaRatingTagOption IN (?)"
	expectedInsertCanteenTag  = "INSERT INTO `canteen_rating_tags` (`correspondingRating`,`tagID`,`points`) VALUES (?,?,?)"
	expectedDishByName        = "SELECT * FROM `dishes` WHERE name = ? AND cafeteriaID = ? ORDER BY `dishes`.`dish` LIMIT ?"
	expectedInsertDishRating  = "INSERT INTO `dish_ratings` (`points`,`dishID`,`comment`,`image`,`device_hash`,`status`,`updated_at`,`timestamp`) VALUES (?,?,?,?,?,?,?,?)"
	expectedDishNameTagIDs    = "SELECT `nameTagID` FROM `dish_to_dish_name_tags` WHERE dishID = ?"
	expectedInsertDishNameTag = "INSERT INTO `dish_name_tags` (`correspondingRating`,`tagNameID`,`points`) VALUES (?,?,?)"
	expectedDishTagSingle     = "SELECT `dishRatingTagOption` FROM `dish_rating_tag_options` WHERE dishRatingTagOption IN (?)"
	expectedInsertDishTag     = "INSERT INTO `dish_rating_tags` (`parentRating`,`tagID`,`points`) VALUES (?,?,?)"
	expectedCanteenIDQuery    = "SELECT `cafeteria` FROM `canteens` WHERE name = ?"
	expectedDishIDQuery       = "SELECT `dish` FROM `dishes` WHERE name = ? AND cafeteriaID = ?"
	expectedCanteenStatsQuery = "SELECT * FROM `cafeteria_rating_statistics` WHERE cafeteriaId = ? ORDER BY `cafeteria_rating_statistics`.`cafeteriaID` LIMIT ?"
	expectedDishStatsQuery    = "SELECT * FROM `dish_rating_statistics` WHERE cafeteriaID = ? AND dishID = ? ORDER BY `dish_rating_statistics`.`cafeteriaID` LIMIT ?"
)

func (s *CanteenSuite) expectCanteenLookup(name string, id int64) {
	s.mock.ExpectQuery(regexp.QuoteMeta(expectedCanteenByName)).
		WithArgs(model.CanteenEnumName(name), 1).
		WillReturnRows(sqlmock.NewRows([]string{"cafeteria", "name", "address", "latitude", "longitude"}).
			AddRow(id, model.CanteenEnumName(name), "Boltzmannstraße 19, 85748 Garching", 48.2681, 11.6723))
}

func (s *CanteenSuite) Test_GetCanteenHeadCount() {
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

// Test_GetCanteenHeadCountEnumName checks that the eat-api enum_name form is accepted as well.
func (s *CanteenSuite) Test_GetCanteenHeadCountEnumName() {
	s.mock.ExpectQuery(regexp.QuoteMeta(expectedHeadCountQuery)).
		WithArgs("mensa-garching", 1).
		WillReturnRows(sqlmock.NewRows([]string{"canteen_id", "count", "max_count", "percent", "timestamp"}).
			AddRow("mensa-garching", 250, 1000, 25, time.Now()))

	server := s.server()
	response, err := server.GetCanteenHeadCount(testCtx(), &pb.GetCanteenHeadCountRequest{CanteenId: "MENSA_GARCHING"})
	require.NoError(s.T(), err)
	require.Equal(s.T(), uint32(250), response.Count)
}

func (s *CanteenSuite) Test_ListCanteens() {
	s.mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `canteens` ORDER BY name")).
		WillReturnRows(sqlmock.NewRows([]string{"cafeteria", "name", "address", "latitude", "longitude"}).
			AddRow(2, "MENSA_GARCHING", "Boltzmannstraße 19, Garching", 48.268132, 11.672263).
			AddRow(7, "STUCAFE_BOLTZMANNSTR", "Boltzmannstraße 15, Garching", 48.26554, 11.66862))

	server := s.server()
	response, err := server.ListCanteens(testCtx(), &pb.ListCanteensRequest{})
	require.NoError(s.T(), err)
	require.Equal(s.T(), &pb.ListCanteensReply{Canteen: []*pb.Canteen{
		{Id: "mensa-garching", Address: "Boltzmannstraße 19, Garching", Latitude: 48.268132, Longitude: 11.672263},
		{Id: "stucafe-boltzmannstr", Address: "Boltzmannstraße 15, Garching", Latitude: 48.26554, Longitude: 11.66862},
	}}, response)
}

func (s *CanteenSuite) Test_ListCanteensDBError() {
	s.mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `canteens` ORDER BY name")).
		WillReturnError(fmt.Errorf("connection refused"))

	server := s.server()
	response, err := server.ListCanteens(testCtx(), &pb.ListCanteensRequest{})
	require.Nil(s.T(), response)
	require.Equal(s.T(), codes.Internal, status.Code(err))
	require.NotContains(s.T(), err.Error(), "connection refused")
}

func (s *CanteenSuite) Test_ListAvailableCanteenTags() {
	s.mock.ExpectQuery(regexp.QuoteMeta("SELECT DE as de, EN as en, cafeteriaRatingTagOption as TagId FROM `canteen_rating_tag_options`")).
		WillReturnRows(sqlmock.NewRows([]string{"de", "en", "TagId"}).
			AddRow("Sauberkeit", "Cleanliness", 1).
			AddRow("Wartezeit", "Waiting time", 2))

	server := s.server()
	response, err := server.ListAvailableCanteenTags(testCtx(), &pb.ListAvailableCanteenTagsRequest{})
	require.NoError(s.T(), err)
	require.Equal(s.T(), &pb.ListAvailableCanteenTagsReply{RatingTags: []*pb.TagsOverview{
		{De: "Sauberkeit", En: "Cleanliness", TagId: 1},
		{De: "Wartezeit", En: "Waiting time", TagId: 2},
	}}, response)
}

func (s *CanteenSuite) Test_ListDishes() {
	s.mock.ExpectQuery(regexp.QuoteMeta("SELECT d.name FROM dishes_of_the_weeks weekly JOIN dishes d ON d.dish = weekly.dishID JOIN canteens c ON c.cafeteria = d.cafeteriaID WHERE (weekly.day = ? AND weekly.week = ? and weekly.year = ?) AND c.name = ?")).
		WithArgs(int32(1), int32(45), int32(2026), "MENSA_GARCHING").
		WillReturnRows(sqlmock.NewRows([]string{"name"}).AddRow("Pasta").AddRow("Curry"))

	server := s.server()
	response, err := server.ListDishes(testCtx(), &pb.ListDishesRequest{CanteenId: "mensa-garching", Year: 2026, Week: 45, Day: 1})
	require.NoError(s.T(), err)
	require.Equal(s.T(), []string{"Pasta", "Curry"}, response.Dish)
}

func (s *CanteenSuite) Test_ListDishesInvalidInput() {
	tests := []struct {
		name            string
		year, week, day int32
	}{
		{"year too low", 2021, 1, 0},
		{"week too low", 2026, 0, 0},
		{"week too high", 2026, 54, 0},
		{"week 53 in a 52-week ISO year", 2025, 53, 0},
		{"day too low", 2026, 1, -1},
		{"day too high", 2026, 1, 5},
	}
	server := s.server()
	for _, tt := range tests {
		s.Run(tt.name, func() {
			_, err := server.ListDishes(testCtx(), &pb.ListDishesRequest{
				CanteenId: "mensa-garching", Year: tt.year, Week: tt.week, Day: tt.day,
			})
			require.Equal(s.T(), codes.InvalidArgument, status.Code(err))
		})
	}
}

func (s *CanteenSuite) Test_ListDishesEmptyAndDBError() {
	server := s.server()
	for _, dbError := range []bool{false, true} {
		query := s.mock.ExpectQuery(regexp.QuoteMeta("SELECT d.name FROM dishes_of_the_weeks weekly JOIN dishes d ON d.dish = weekly.dishID JOIN canteens c ON c.cafeteria = d.cafeteriaID WHERE (weekly.day = ? AND weekly.week = ? and weekly.year = ?) AND c.name = ?")).
			WithArgs(int32(4), int32(52), int32(2022), "DOES_NOT_EXIST")
		if dbError {
			query.WillReturnError(fmt.Errorf("connection refused"))
		} else {
			query.WillReturnRows(sqlmock.NewRows([]string{"name"}))
		}
		response, err := server.ListDishes(testCtx(), &pb.ListDishesRequest{
			CanteenId: "DOES_NOT_EXIST", Year: 2022, Week: 52, Day: 4,
		})
		if dbError {
			require.Equal(s.T(), codes.Internal, status.Code(err))
			require.NotContains(s.T(), err.Error(), "connection refused")
		} else {
			require.NoError(s.T(), err)
			require.Empty(s.T(), response.Dish)
		}
	}
}

func (s *CanteenSuite) Test_ListCanteenRatingsEmptyCanteen() {
	server := s.server()
	response, err := server.ListCanteenRatings(testCtx(), &pb.ListCanteenRatingsRequest{})
	require.Nil(s.T(), response)
	require.Equal(s.T(), codes.InvalidArgument, status.Code(err))
}

func (s *CanteenSuite) Test_RatingLookupErrors() {
	tests := []struct {
		name    string
		dish    bool
		stage   string
		dbError bool
	}{
		{"unknown canteen", false, "canteen", false},
		{"canteen lookup error", false, "canteen", true},
		{"unrated canteen", false, "statistics", false},
		{"canteen statistics error", false, "statistics", true},
		{"dish unknown canteen", true, "canteen", false},
		{"dish canteen lookup error", true, "canteen", true},
		{"unknown dish", true, "dish", false},
		{"dish lookup error", true, "dish", true},
		{"unrated dish", true, "statistics", false},
		{"dish statistics error", true, "statistics", true},
	}
	server := s.server()
	for _, tt := range tests {
		s.Run(tt.name, func() {
			canteenRows := sqlmock.NewRows([]string{"cafeteria"})
			query := s.mock.ExpectQuery(regexp.QuoteMeta(expectedCanteenIDQuery)).
				WithArgs("MENSA_GARCHING")
			if tt.stage == "canteen" {
				if tt.dbError {
					query.WillReturnError(fmt.Errorf("connection refused"))
				} else {
					query.WillReturnRows(canteenRows)
				}
			} else {
				query.WillReturnRows(canteenRows.AddRow(2))
				if tt.dish {
					dishRows := sqlmock.NewRows([]string{"dish"})
					dishQuery := s.mock.ExpectQuery(regexp.QuoteMeta(expectedDishIDQuery)).
						WithArgs("Pasta", int32(2))
					if tt.stage == "dish" {
						if tt.dbError {
							dishQuery.WillReturnError(fmt.Errorf("connection refused"))
						} else {
							dishQuery.WillReturnRows(dishRows)
						}
					} else {
						dishQuery.WillReturnRows(dishRows.AddRow(13))
					}
				}
				if tt.stage == "statistics" {
					var statsQuery *sqlmock.ExpectedQuery
					if tt.dish {
						statsQuery = s.mock.ExpectQuery(regexp.QuoteMeta(expectedDishStatsQuery)).
							WithArgs(int32(2), int32(13), 1)
					} else {
						statsQuery = s.mock.ExpectQuery(regexp.QuoteMeta(expectedCanteenStatsQuery)).
							WithArgs(int32(2), 1)
					}
					if tt.dbError {
						statsQuery.WillReturnError(fmt.Errorf("connection refused"))
					} else {
						statsQuery.WillReturnRows(sqlmock.NewRows([]string{"cafeteriaID", "average"}))
					}
				}
			}
			var err error
			if tt.dish {
				var response *pb.GetDishRatingsReply
				response, err = server.GetDishRatings(testCtx(), &pb.GetDishRatingsRequest{CanteenId: "mensa-garching", Dish: "Pasta"})
				require.Nil(s.T(), response)
			} else {
				var response *pb.ListCanteenRatingsReply
				response, err = server.ListCanteenRatings(testCtx(), &pb.ListCanteenRatingsRequest{CanteenId: "mensa-garching"})
				require.Nil(s.T(), response)
			}
			expectedCode := codes.NotFound
			if tt.dbError {
				expectedCode = codes.Internal
				require.NotContains(s.T(), err.Error(), "connection refused")
			}
			require.Equal(s.T(), expectedCode, status.Code(err))
			require.NoError(s.T(), s.mock.ExpectationsWereMet())
		})
	}
}

func (s *CanteenSuite) Test_ValidationOverREST() {
	server := s.server()
	mux := runtime.NewServeMux()
	require.NoError(s.T(), pb.RegisterCampusHandlerServer(context.Background(), mux, &server))
	handler := http.StripPrefix("/v1", mux)
	tests := []struct {
		method, url, body string
		statusCode        int
	}{
		{http.MethodGet, "/v1/dishes?canteenId=mensa-garching&year=2026&week=99&day=0", "", http.StatusBadRequest},
		{http.MethodPost, "/v1/canteen/rating/get", `{}`, http.StatusBadRequest},
		{http.MethodPost, "/v1/canteen/rating/get", `{"canteenId":"DOES_NOT_EXIST"}`, http.StatusNotFound},
		{http.MethodPost, "/v1/dish/rating/get", `{}`, http.StatusNotFound},
	}
	s.mock.ExpectQuery(regexp.QuoteMeta(expectedCanteenIDQuery)).
		WithArgs("DOES_NOT_EXIST").
		WillReturnRows(sqlmock.NewRows([]string{"cafeteria"}))
	for _, tt := range tests {
		request := httptest.NewRequest(tt.method, tt.url, strings.NewReader(tt.body))
		request.Header.Set("Content-Type", "application/json")
		response := httptest.NewRecorder()
		handler.ServeHTTP(response, request)
		require.Equal(s.T(), tt.statusCode, response.Code, response.Body.String())
	}
}

func (s *CanteenSuite) Test_GetCanteenHeadCountNotFound() {
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

func (s *CanteenSuite) Test_CreateCanteenRating() {
	const (
		canteenName = "mensa-garching"
		comment     = "Everything perfect, 2 Stars"
		canteenID   = int64(2)
		ratingID    = int64(42)
		points      = int32(2)
	)
	img := readTestImage(s.T())
	imgPath := expectedImagePath(img, "cafeterias", canteenID)

	s.expectCanteenLookup(canteenName, canteenID)
	s.mock.ExpectQuery(regexp.QuoteMeta(expectedCanteenTagsIn)).
		WithArgs(int64(5), int64(7)).
		WillReturnRows(sqlmock.NewRows([]string{"cafeteriaRatingTagOption"}).AddRow(5).AddRow(7))
	s.mock.ExpectBegin()
	s.mock.ExpectExec(regexp.QuoteMeta(expectedInsertRating)).
		WithArgs(points, comment, canteenID, imgPath, sqlmock.AnyArg(), model.RatingStatusVisible, sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(ratingID, 1))
	for _, tag := range []struct{ id, points int64 }{{5, 1}, {7, 3}} {
		s.mock.ExpectExec(regexp.QuoteMeta(expectedInsertCanteenTag)).
			WithArgs(ratingID, tag.id, tag.points).
			WillReturnResult(sqlmock.NewResult(1, 1))
	}
	s.mock.ExpectCommit()

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
	require.Equal(s.T(), &pb.CreateCanteenRatingReply{RatingId: ratingID}, response)
	require.FileExists(s.T(), imgPath)
}

func (s *CanteenSuite) Test_CreateCanteenRatingUnknownTag() {
	const (
		canteenName = "mensa-garching"
		canteenID   = int64(2)
	)
	s.expectCanteenLookup(canteenName, canteenID)
	s.mock.ExpectQuery(regexp.QuoteMeta(expectedCanteenTagSingle)).
		WithArgs(int64(99)).
		WillReturnRows(sqlmock.NewRows([]string{"cafeteriaRatingTagOption"}))

	server := s.server()
	_, err := server.CreateCanteenRating(testCtx(), &pb.CreateCanteenRatingRequest{
		Points:     4,
		CanteenId:  canteenName,
		RatingTags: []*pb.RatingTag{{TagId: 99, Points: 1}},
	})
	require.Equal(s.T(), codes.InvalidArgument, status.Code(err))
}

func (s *CanteenSuite) Test_CreateCanteenRatingInvalidInput() {
	tests := []struct {
		name string
		req  *pb.CreateCanteenRatingRequest
	}{
		{"points too high", &pb.CreateCanteenRatingRequest{Points: 6, CanteenId: "mensa-garching"}},
		{"points negative", &pb.CreateCanteenRatingRequest{Points: -1, CanteenId: "mensa-garching"}},
		{"points zero", &pb.CreateCanteenRatingRequest{Points: 0, CanteenId: "mensa-garching"}},
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

func (s *CanteenSuite) Test_CreateCanteenRatingUnknownCanteen() {
	s.mock.ExpectQuery(regexp.QuoteMeta(expectedCanteenByName)).
		WithArgs("MENSA_UNKNOWN", 1).
		WillReturnRows(sqlmock.NewRows([]string{"cafeteria", "name", "address", "latitude", "longitude"}))

	server := s.server()
	response, err := server.CreateCanteenRating(testCtx(), &pb.CreateCanteenRatingRequest{Points: 3, CanteenId: "mensa-unknown"})
	require.Nil(s.T(), response)
	require.Equal(s.T(), status.Error(codes.InvalidArgument, "Canteen does not exist. Rating has not been saved."), err)
}

func (s *CanteenSuite) Test_CreateDishRating() {
	const (
		canteenName = "MENSA_GARCHING"
		dishName    = "Vegane rote Grütze mit Soja-Vanillesauce"
		comment     = "Very tasty"
		canteenID   = int64(2)
		dishID      = int64(13)
		ratingID    = int64(42)
		nameTagID   = int64(3)
		points      = int32(5)
	)
	img := readTestImage(s.T())
	imgPath := expectedImagePath(img, "dishes", dishID)

	s.expectCanteenLookup(canteenName, canteenID)
	s.mock.ExpectQuery(regexp.QuoteMeta(expectedDishByName)).
		WithArgs(dishName, canteenID, 1).
		WillReturnRows(sqlmock.NewRows([]string{"dish", "name", "type", "cafeteriaID"}).
			AddRow(dishID, dishName, "Dessert", canteenID))
	s.mock.ExpectQuery(regexp.QuoteMeta(expectedDishTagSingle)).
		WithArgs(int64(1)).
		WillReturnRows(sqlmock.NewRows([]string{"dishRatingTagOption"}).AddRow(1))
	s.mock.ExpectBegin()
	s.mock.ExpectExec(regexp.QuoteMeta(expectedInsertDishRating)).
		WithArgs(points, dishID, comment, imgPath, sqlmock.AnyArg(), model.RatingStatusVisible, sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(ratingID, 1))
	s.mock.ExpectQuery(regexp.QuoteMeta(expectedDishNameTagIDs)).
		WithArgs(dishID).
		WillReturnRows(sqlmock.NewRows([]string{"nameTagID"}).AddRow(nameTagID))
	s.mock.ExpectExec(regexp.QuoteMeta(expectedInsertDishNameTag)).
		WithArgs(ratingID, nameTagID, points).
		WillReturnResult(sqlmock.NewResult(1, 1))
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
	require.Equal(s.T(), &pb.CreateDishRatingReply{RatingId: ratingID}, response)
	require.FileExists(s.T(), imgPath)
}

func (s *CanteenSuite) Test_CreateDishRatingInvalidImageAndTagValidation() {
	server := s.server()

	tcases := []struct {
		name string
		tags []*pb.RatingTag
	}{
		{"tag points too low", []*pb.RatingTag{{TagId: 1, Points: 0}}},
		{"tag points too high", []*pb.RatingTag{{TagId: 1, Points: 6}}},
		{"duplicate tag id", []*pb.RatingTag{{TagId: 1, Points: 3}, {TagId: 1, Points: 4}}},
	}
	for _, tc := range tcases {
		s.Run(tc.name, func() {
			_, err := server.CreateDishRating(testCtx(), &pb.CreateDishRatingRequest{
				Points:     4,
				CanteenId:  "mensa-garching",
				Dish:       "Pasta",
				RatingTags: tc.tags,
			})
			require.Equal(s.T(), codes.InvalidArgument, status.Code(err))
		})
	}

	s.Run("corrupt image returns InvalidArgument without panic", func() {
		s.expectCanteenLookup("mensa-garching", 2)
		s.mock.ExpectQuery(regexp.QuoteMeta(expectedDishByName)).
			WithArgs("Pasta", int64(2), 1).
			WillReturnRows(sqlmock.NewRows([]string{"dish", "name", "type", "cafeteriaID"}).AddRow(13, "Pasta", "Pasta", 2))
		_, err := server.CreateDishRating(testCtx(), &pb.CreateDishRatingRequest{
			Points:    4,
			CanteenId: "mensa-garching",
			Dish:      "Pasta",
			Image:     []byte("not a valid image"),
		})
		require.Equal(s.T(), codes.InvalidArgument, status.Code(err))
	})

	s.Run("256 UTF-8 umlaut runes accepted while 257 rejected", func() {
		s.expectCanteenLookup("mensa-garching", 2)
		id, err := inputSanitizationForNewRatingElements(5, strings.Repeat("ü", 256), "mensa-garching", s.DB)
		require.NoError(s.T(), err)
		require.Equal(s.T(), int64(2), id)

		_, err = inputSanitizationForNewRatingElements(5, strings.Repeat("ü", 257), "mensa-garching", s.DB)
		require.Equal(s.T(), codes.InvalidArgument, status.Code(err))
	})
}

func (s *CanteenSuite) Test_CreateDishRatingUnknownDish() {
	const canteenID = int64(2)
	s.expectCanteenLookup("MENSA_GARCHING", canteenID)
	s.mock.ExpectQuery(regexp.QuoteMeta(expectedDishByName)).
		WithArgs("Unknown dish", canteenID, 1).
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

func (s *CanteenSuite) Test_CreateDishRatingInvalidPoints() {
	server := s.server()
	for _, points := range []int32{-1, 0, 6} {
		response, err := server.CreateDishRating(testCtx(), &pb.CreateDishRatingRequest{
			Points: points, CanteenId: "mensa-garching", Dish: "Pasta",
		})
		require.Nil(s.T(), response)
		require.Equal(s.T(), codes.InvalidArgument, status.Code(err))
	}
}

func (s *CanteenSuite) Test_RatingPointsBoundaries() {
	for _, points := range []int32{1, 5} {
		s.expectCanteenLookup("mensa-garching", 2)
		id, err := inputSanitizationForNewRatingElements(points, "", "mensa-garching", s.DB)
		require.NoError(s.T(), err)
		require.Equal(s.T(), int64(2), id)
	}
}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestCanteenSuite(t *testing.T) {
	suite.Run(t, new(CanteenSuite))
}

func TestStoreImageSupportsPNGAndScalesLargeImagesCleanly(t *testing.T) {
	t.Run("PNG upload decodes and stores valid JPEG", func(t *testing.T) {
		img := image.NewRGBA(image.Rect(0, 0, 64, 48))
		for y := 0; y < 48; y++ {
			for x := 0; x < 64; x++ {
				img.Set(x, y, color.RGBA{R: uint8(x * 4), G: uint8(y * 5), B: 128, A: 255})
			}
		}
		var buf bytes.Buffer
		require.NoError(t, png.Encode(&buf, img))

		dir := t.TempDir()
		storedPath, err := storeImage(dir, buf.Bytes())
		require.NoError(t, err)
		require.FileExists(t, storedPath)

		storedBytes, err := os.ReadFile(storedPath)
		require.NoError(t, err)
		decoded, err := jpeg.Decode(bytes.NewReader(storedBytes))
		require.NoError(t, err)
		require.Equal(t, 64, decoded.Bounds().Dx())
		require.Equal(t, 48, decoded.Bounds().Dy())
	})

	t.Run("large image over 512KiB resizes and compresses with clamped quality", func(t *testing.T) {
		// Uncompressed 1600x400 PNG exceeds targetMaxImageBytes (~2.5 MiB) and triggers both width resize and quality clamping.
		img := image.NewRGBA(image.Rect(0, 0, 1600, 400))
		for y := 0; y < 400; y++ {
			for x := 0; x < 1600; x++ {
				img.Set(x, y, color.RGBA{
					R: uint8((x*37 + y*17) & 0xff),
					G: uint8((x*13 ^ y*29) & 0xff),
					B: uint8((x*7 + y*41) & 0xff),
					A: 255,
				})
			}
		}
		var buf bytes.Buffer
		enc := png.Encoder{CompressionLevel: png.NoCompression}
		require.NoError(t, enc.Encode(&buf, img))
		require.Greater(t, buf.Len(), targetMaxImageBytes)

		dir := t.TempDir()
		storedPath, err := storeImage(dir, buf.Bytes())
		require.NoError(t, err)

		storedBytes, err := os.ReadFile(storedPath)
		require.NoError(t, err)
		decoded, err := jpeg.Decode(bytes.NewReader(storedBytes))
		require.NoError(t, err)
		require.Equal(t, maxRatingImageWidth, decoded.Bounds().Dx(), "wide image should be resized to maxRatingImageWidth")
		require.Greater(t, len(storedBytes), 20*1024, "compressed JPEG must not collapse to quality=1 artifact")
	})
}
