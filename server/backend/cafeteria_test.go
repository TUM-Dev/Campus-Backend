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
	StorageDir = "."
	canteenName := "mensa-garching"
	cafeteriaID := 0
	comment := "Everything perfect, 2 Stars"
	image_name_path := ".cafeterias/0/c73f6461a5ae03bb935b54a27f12e0a5.jpeg" // "testimage.txt"
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
			AddRow(0, ratingValue, comment, cafeteriaID, time.Date(2023, 10, 9, 11, 45, 22, 0, time.Local), image_name_path))

	/* todo current exec	s.mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `cafeteria_rating` (`points`,`comment`,`cafeteriaID`,`timestamp`,`image`) VALUES (?,?,?,?,?) RETURNING `cafeteriaRating`")).
	WithArgs(ratingValue, "Everything perfect, 2 Star", canteenName, sqlmock.AnyArg(), image_name_path).
	WillReturnResult(sqlmock.NewResult(1, 1))
	*/
	s.mock.ExpectCommit()

	meta := metadata.MD{}
	rating := generateCanteenRating(canteenName, ratingValue, s, comment)

	server := CampusServer{db: s.DB, deviceBuf: s.deviceBuf}
	response, err := server.CreateCanteenRating(metadata.NewIncomingContext(context.Background(), meta),
		&rating,
	)
	require.NoError(s.T(), err)
	if err != nil {
		log.WithError(err).Error("Canteen HeadCount data request failed.")
		//todo compare results	require.Equal(s.T(), status.Error(codes.NotFound, "No update note found"), err)
	} else {
		log.WithField("res", response).Info("Canteen HeadCount data request successful.")
	}
}

func generateCanteenRating(canteen string, rating int32, s *CafeteriaSuite, comment string) pb.CreateCanteenRatingRequest {
	//canteen_id := 1

	//dummyImage := createDummyImage(s.T(), 10, 10)
	//s.mock.ExpectBegin()
	/*s.mock.ExpectQuery(regexp.QuoteMeta("INSERT INTO `cafeteria_rating` (`cafeteriaRating`,`points`,`comment`,`cafeteriaID`,`timestamp`,`image`) VALUES (?,?,?,?,?,?) RETURNING `cafeteriaRating`,`points`,`comment`,`cafeteriaID`,`timestamp`,`image`")).
	WithArgs(1, 2, "custom comment", canteen_id,time.Date(2023, 10, 9, 11, 45, 22, 0, time.Local),image_name_path).
		WillReturnRows(sqlmock.NewRows([]string{"url", "file"}).AddRow(nil, 1))
	*/
	//	var help int32 = 3
	/*	s.mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `cafeteria_rating` (`points`,`comment`,`cafeteriaID`,`timestamp`,`image`) VALUES (?,?,?,?,?) RETURNING `cafeteriaRating`")).
		WithArgs(rating, "Everything perfect, 2 Star", canteen, sqlmock.AnyArg(), image_name_path).
		WillReturnResult(sqlmock.NewResult(1, 1))*/

	/*
		s.mock.ExpectQuery(regexp.QuoteMeta("INSERT INTO `files` (`name`,`path`,`downloads`,`downloaded`) VALUES (?,?,?,?) RETURNING `url`,`file`")).
		WithArgs("0.txt", sqlmock.AnyArg(), 1, true).
		WillReturnRows(sqlmock.NewRows([]string{"url", "file"}).AddRow(nil, 1))
	*/

	/*	s.mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `cafeteria_rating` (`points`,`comment`,`cafeteriaID`,`image`) VALUES (?,?,?,?) RETURNING `cafeteriaRating`")).
		WithArgs(rating, "Everything perfect, 2 Star", canteen, image_name_path)*/
	//WillReturnRows(sqlmock.NewRows([]string{"cafeteria"}).
	//	AddRow(2))
	//s.mock.ExpectCommit()

	y := make([]*pb.RatingTag, 2)
	y[0] = &pb.RatingTag{
		Points: float64(1 + rating),
		TagId:  1,
	}
	y[1] = &pb.RatingTag{
		Points: float64(2 + rating),
		TagId:  2,
	}

	return pb.CreateCanteenRatingRequest{
		Points:     rating,
		CanteenId:  canteen,
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
