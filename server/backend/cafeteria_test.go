package backend

import (
	"context"
	"database/sql"
	pb "github.com/TUM-Dev/Campus-Backend/server/api/tumdev"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc/metadata"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

const (
	testImage = "./test_data/sampleimage.jpeg"
)

type CafeteriaSuite struct {
	suite.Suite
	DB        *gorm.DB
	mock      sqlmock.Sqlmock
	deviceBuf *deviceBuffer
}

func (s *CafeteriaSuite) SetupSuite() {
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
		WillReturnRows(sqlmock.NewRows([]string{"VERSION()"}).AddRow("10.11.4-MariaDB"))
	s.DB, err = gorm.Open(dialector, &gorm.Config{})
	require.NoError(s.T(), err)

	s.deviceBuf = newDeviceBuffer()
}

const ExpectedGetCafeteriaQuery = "SELECT * FROM `canteen_head_count` WHERE `canteen_head_count`.`canteen_id` = ? ORDER BY `canteen_head_count`.`canteen_id` LIMIT 1"

func (s *CafeteriaSuite) Test_GetCafeteriaHeadCount() {
	s.mock.ExpectQuery(regexp.QuoteMeta(ExpectedGetCafeteriaQuery)).
		WithArgs("mensa-garching").
		WillReturnRows(sqlmock.NewRows([]string{"canteen_id", "count", "max_count", "percent", "timestamp"}).
			AddRow("mensa-garching", 0, 1000, 0, time.Date(2023, 10, 9, 11, 45, 22, 0, time.Local)))

	meta := metadata.MD{}
	server := CampusServer{db: s.DB, deviceBuf: s.deviceBuf}
	response, err := server.GetCanteenHeadCount(metadata.NewIncomingContext(context.Background(), meta),
		&pb.GetCanteenHeadCountRequest{
			CanteenId: "mensa-garching",
		},
	)
	require.NoError(s.T(), err)
	if err != nil {
		log.WithError(err).Error("Canteen HeadCount data request failed.")
		//todo compare results	require.Equal(s.T(), status.Error(codes.NotFound, "No update note found"), err)
	} else {
		log.WithField("res", response).Info("Canteen HeadCount data request successful.")
	}
}

const ExpectedGetCafeteriaByName = "SELECT * FROM `cafeteria` WHERE name LIKE ? ORDER BY `cafeteria`.`cafeteria` LIMIT 1"

func (s *CafeteriaSuite) Test_CreateCanteenRating() {
	StorageDir = "." // Override directory used in the production environment
	canteenName := "mensa-garching"
	cafeteriaID := 0
	comment := "Everything perfect, 2 Stars"
	ratingID := 0
	image_name_path := ".cafeterias/0/c73f6461a5ae03bb935b54a27f12e0a5.jpeg"
	var ratingValue int32 = 2

	//	s.mock.ExpectBegin()
	s.mock.ExpectQuery(regexp.QuoteMeta(ExpectedGetCafeteriaByName)).
		WithArgs(canteenName).
		WillReturnRows(sqlmock.NewRows([]string{"cafeteria", "name", "address", "latitude", "longitude"}).
			AddRow(cafeteriaID, "mensa-garching", "Boltzmannstraße 19, Garching", 48.2681, 11.6723))

	s.mock.ExpectBegin()
	s.mock.ExpectQuery(regexp.QuoteMeta("INSERT INTO `cafeteria_rating` (`points`,`comment`,`cafeteriaID`,`timestamp`,`image`) VALUES (?,?,?,?,?) RETURNING `cafeteriaRating`")).
		WithArgs(ratingValue, comment, cafeteriaID, sqlmock.AnyArg(), image_name_path).
		WillReturnRows(sqlmock.NewRows([]string{"cafeteriaRating", "points", "comment", "cafeteriaID", "timestamp", "image"}).
			AddRow(ratingID, ratingValue, comment, cafeteriaID, time.Date(2023, 10, 9, 11, 45, 22, 0, time.Local), image_name_path))
	s.mock.ExpectCommit()

	meta := metadata.MD{}
	rating := generateCanteenRating(canteenName, ratingValue, s, comment)

	server := CampusServer{db: s.DB, deviceBuf: s.deviceBuf}
	response, err := server.CreateCanteenRating(metadata.NewIncomingContext(context.Background(), meta),
		&rating,
	)
	// todo test the response
	require.NoError(s.T(), err)
	if err != nil {
		log.WithError(err).Error("Canteen HeadCount data request failed.")
		//todo compare results	require.Equal(s.T(), status.Error(codes.NotFound, "No update note found"), err)
	} else {
		log.WithField("res", response).Info("Canteen HeadCount data request successful.")
	}
}

func generateCanteenRating(canteen string, rating int32, s *CafeteriaSuite, comment string) pb.CreateCanteenRatingRequest {
	y := make([]*pb.RatingTag, 2)
	var myRating = prepareTagRating(s, 1, 0, 1, 5)
	y[0] = &myRating
	var myRatingSecond = prepareTagRating(s, 2, 0, 3, 7)
	y[1] = &myRatingSecond

	return pb.CreateCanteenRatingRequest{
		Points:     rating,
		CanteenId:  canteen,
		Comment:    comment,
		RatingTags: y,
		Image:      getImageToBytes(testImage),
	}
}

func prepareTagRating(s *CafeteriaSuite, tagRatingID, ratingID int, points int, tagID int) pb.RatingTag {
	// Query to check whether the tag exists
	const ExpectedCafeteriaRatingTagOption = "SELECT count(*) FROM `cafeteria_rating_tag_option` WHERE cafeteriaRatingTagOption LIKE ?"
	s.mock.ExpectQuery(regexp.QuoteMeta(ExpectedCafeteriaRatingTagOption)).
		WithArgs(tagID).WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))

	s.mock.ExpectBegin()
	s.mock.ExpectQuery(regexp.QuoteMeta("INSERT INTO `cafeteria_rating_tag` (`correspondingRating`,`points`,`tagID`) VALUES (?,?,?) RETURNING `CafeteriaRatingTag`")).
		WithArgs(ratingID, points, tagID).
		WillReturnRows(sqlmock.NewRows([]string{"CafeteriaRatingTag", "correspondingRating", "points", "tagID"}).
			AddRow(tagRatingID, ratingID, points, tagID))
	s.mock.ExpectCommit()

	return pb.RatingTag{
		Points: float64(points),
		TagId:  int64(tagID),
	}
}

const ExpectedGetDishByName = "SELECT * FROM `dish` WHERE name LIKE ? ORDER BY `dish`.`dish` LIMIT 1"

// Test if dish ratings are correctly created
func (s *CafeteriaSuite) Test_CreateDishRating() {
	StorageDir = "." // Override directory used in the production environment
	canteenName := "MENSA_GARCHING"
	dishName := "Vegane rote Grütze mit Soja-Vanillesauce"
	dishId := 0
	dishType := "Pasta"
	cafeteriaID := 0
	comment := "Everything perfect, 2 Stars"
	dishRatingID := 0
	image_name_path := ".cafeterias/0/c73f6461a5ae03bb935b54a27f12e0a5.jpeg"
	var ratingValue int32 = 2

	//	s.mock.ExpectBegin()
	s.mock.ExpectQuery(regexp.QuoteMeta(ExpectedGetDishByName)).
		WithArgs(dishName).
		WillReturnRows(sqlmock.NewRows([]string{"dish", "name", "type", "cafeteriaID"}).
			AddRow(dishId, dishName, dishType, cafeteriaID))

	s.mock.ExpectBegin()
	s.mock.ExpectQuery(regexp.QuoteMeta("INSERT INTO `dish_rating` (`dishRating`,`points`,`cafeteriaID`,`dishID`,`comment`,`timestamp`,`image`) VALUES (?,?,?,?,?,?) RETURNING `dishRating`")).
		WithArgs(ratingValue, comment, cafeteriaID, sqlmock.AnyArg(), image_name_path).
		WillReturnRows(sqlmock.NewRows([]string{"dishRating", "points", "cafeteriaID", "dishID", "comment", "timestamp", "image"}).
			AddRow(dishRatingID, ratingValue, cafeteriaID, dishId, comment, time.Date(2023, 10, 9, 11, 45, 22, 0, time.Local), image_name_path))
	s.mock.ExpectCommit()

	meta := metadata.MD{}
	rating := generateDishRating(canteenName, ratingValue, s, comment, dishName)

	server := CampusServer{db: s.DB, deviceBuf: s.deviceBuf}
	response, err := server.CreateDishRating(metadata.NewIncomingContext(context.Background(), meta),
		&rating,
	)
	// todo test the response
	require.NoError(s.T(), err)
	if err != nil {
		log.WithError(err).Error("Canteen HeadCount data request failed.")
		//todo compare results	require.Equal(s.T(), status.Error(codes.NotFound, "No update note found"), err)
	} else {
		log.WithField("res", response).Info("Canteen HeadCount data request successful.")
	}
}

func generateDishRating(canteen string, rating int32, s *CafeteriaSuite, comment string, dishName string) pb.CreateDishRatingRequest {
	y := make([]*pb.RatingTag, 2)
	// todo cahnge to dish tags
	// todo add the dish specific tags
	var myRating = prepareTagRating(s, 1, 0, 1, 5)
	y[0] = &myRating
	var myRatingSecond = prepareTagRating(s, 2, 0, 3, 7)
	y[1] = &myRatingSecond

	return pb.CreateDishRatingRequest{
		Points:     rating,
		CanteenId:  canteen,
		Dish:       dishName,
		Comment:    comment,
		RatingTags: y,
		Image:      getImageToBytes(testImage),
	}
}

func (s *CafeteriaSuite) AfterTest(_, _ string) {
	require.NoError(s.T(), s.mock.ExpectationsWereMet())
}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestCafeteriaSuite(t *testing.T) {
	suite.Run(t, new(CafeteriaSuite))
}
