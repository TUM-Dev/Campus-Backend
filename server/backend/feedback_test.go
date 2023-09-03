package backend

import (
	"context"
	"database/sql"
	"github.com/DATA-DOG/go-sqlmock"
	pb "github.com/TUM-Dev/Campus-Backend/server/api"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/emptypb"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"regexp"
	"testing"
	"time"
)

type FeedbackSuite struct {
	suite.Suite
	DB        *gorm.DB
	mock      sqlmock.Sqlmock
	deviceBuf *deviceBuffer
}

func (s *FeedbackSuite) SetupSuite() {
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

func (s *FeedbackSuite) Test_SendFeedback() {
	req := &pb.SendFeedbackRequest{
		Topic:      "tca",
		Email:      "testing@example.com",
		EmailId:    "magic-id",
		Message:    "This is a Test",
		ImageCount: 1,
		Latitude:   0,
		Longitude:  0,
		OsVersion:  "Android 10.0",
		AppVersion: "TCA 10.2",
	}
	s.mock.ExpectBegin()
	s.mock.ExpectQuery(regexp.QuoteMeta(
		"INSERT INTO `feedback` (`image_count`,`email_id`,`receiver`,`reply_to`,`feedback`,`latitude`,`longitude`,`os_version`,`app_version`,`processed`) VALUES (?,?,?,?,?,?,?,?,?,?) RETURNING `timestamp`,`id`")).
		WithArgs(req.ImageCount, req.EmailId, "app@tum.de", req.Email, req.Message, req.Latitude, req.Longitude, req.OsVersion, req.AppVersion, false).
		WillReturnRows(sqlmock.NewRows([]string{"timestamp", "id"}).AddRow(time.Time{}, 1))
	s.mock.ExpectCommit()

	meta := metadata.MD{}
	server := CampusServer{db: s.DB, deviceBuf: s.deviceBuf}
	response, err := server.SendFeedback(metadata.NewIncomingContext(context.Background(), meta), req)
	require.NoError(s.T(), err)
	require.Equal(s.T(), &emptypb.Empty{}, response)
}
func (s *FeedbackSuite) AfterTest(_, _ string) {
	require.NoError(s.T(), s.mock.ExpectationsWereMet())
}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestExampleTestSuite(t *testing.T) {
	suite.Run(t, new(FeedbackSuite))
}
