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

//"SELECT * FROM `update_note` WHERE `update_note`.`version_code` = ? ORDER BY `update_note`.`version_code` LIMIT 1"

func (s *CafeteriaSuite) Test_GetCafeteriaHeadCount() {
	//inputString := "2023-10-09 11:45:22"

	// Define a custom layout that matches the input string format
	//layout := "2006-01-02 15:04:05"
	test := time.Date(1970, 1, 1, 0, 0, 0, 0, time.Local)

	// Parse the input string into a time.Time object
	//t, err := time.Parse(layout, inputString)
	/*timeWrapper := timestamp.Timestamp{
		Seconds: t.Unix(),              // The number of seconds since January 1, 1970
		Nanos:   int32(t.Nanosecond()), // The nanoseconds part
	}*/
	s.mock.ExpectQuery(regexp.QuoteMeta(ExpectedGetCafeteriaQuery)).
		WithArgs("mensa-garching").
		WillReturnRows(sqlmock.NewRows([]string{"canteen_id", "count", "max_count", "percent", "timestamp"}).
			AddRow("mensa-garching", 0, 1000, 0, test))
	/*
		all expectations were already fulfilled, call to Query 'SELECT * FROM `canteen_head_count` WHERE `canteen_head_count`.`canteen_id` = ? ORDER BY `canteen_head_count`.`canteen_id` LIMIT 1'
	*/
	/*
		meta := metadata.MD{}
		server := CampusServer{db: s.DB, deviceBuf: s.deviceBuf}
		response, err := server.GetCafeteria(metadata.NewIncomingContext(context.Background(), meta), &pb.GetCafeteriaRequest{Version: 1})
		require.NoError(s.T(), err)
		expectedResp := &pb.GetCafeteriaReply{
			Message:     "Test Message",
			VersionName: "1.0.0",
		}*/
	meta := metadata.MD{}
	server := CampusServer{db: s.DB, deviceBuf: s.deviceBuf}
	response, err := server.GetCanteenHeadCount(metadata.NewIncomingContext(context.Background(), meta),
		&pb.GetCanteenHeadCountRequest{
			CanteenId: "mensa-garching",
		},
	)
	require.NoError(s.T(), err)
	//	require.NoError(s.T(), err)
	if err != nil {
		log.WithError(err).Error("Canteen HeadCount data request failed.")
	} else {
		log.WithField("res", response).Info("Canteen HeadCount data request successful.")
	}
	//	require.Equal(s.T(), expectedResp, response)
}

/*
meta := metadata.MD{}

		server := CampusServer{db: s.DB, deviceBuf: s.deviceBuf}
		res, err := server.GetCanteenHeadCount(ctx, &pb.GetCanteenHeadCountRequest{
			CanteenId: "mensa-garching",
		})

		if err != nil {
			log.WithError(err).Error("Canteen HeadCount data request failed.")
		} else {
			log.WithField("res", res).Info("Canteen HeadCount data request successful.")
		}
		require.Equal(s.T(), expectedResp, response)
	}

/*

	func (s *CafeteriaSuite) Test_GetCafeteriaNone() {
		s.mock.ExpectQuery(regexp.QuoteMeta(ExpectedGetCafeteriaQuery)).
			WillReturnRows(sqlmock.NewRows([]string{"version_code", "version_name", "message"}))

		meta := metadata.MD{}
		server := CampusServer{db: s.DB, deviceBuf: s.deviceBuf}
		response, err := server.GetCafeteria(metadata.NewIncomingContext(context.Background(), meta), &pb.GetCafeteriaRequest{Version: 1})
		require.Equal(s.T(), status.Error(codes.NotFound, "No update note found"), err)
		require.Nil(s.T(), response)
	}

	func (s *CafeteriaSuite) Test_GetCafeteriaError() {
		s.mock.ExpectQuery(regexp.QuoteMeta(ExpectedGetCafeteriaQuery)).WillReturnError(gorm.ErrInvalidDB)

		meta := metadata.MD{}
		server := CampusServer{db: s.DB, deviceBuf: s.deviceBuf}
		response, err := server.GetCafeteria(metadata.NewIncomingContext(context.Background(), meta), &pb.GetCafeteriaRequest{Version: 1})
		require.Equal(s.T(), status.Error(codes.Internal, "Internal server error"), err)
		require.Nil(s.T(), response)
	}
*/
func (s *CafeteriaSuite) AfterTest(_, _ string) {
	require.NoError(s.T(), s.mock.ExpectationsWereMet())
}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestCafeteriaSuite(t *testing.T) {
	suite.Run(t, new(CafeteriaSuite))
}
