package backend

import (
	"context"
	"database/sql"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	pb "github.com/TUM-Dev/Campus-Backend/server/api/tumdev"
	"github.com/TUM-Dev/Campus-Backend/server/model"
	"github.com/guregu/null"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type NewsAlertSuite struct {
	suite.Suite
	DB        *gorm.DB
	mock      sqlmock.Sqlmock
	deviceBuf *deviceBuffer
}

func (s *NewsAlertSuite) SetupSuite() {
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

const ExpectedGetTopNewsQuery = "SELECT `news_alert`.`news_alert`,`news_alert`.`file`,`news_alert`.`name`,`news_alert`.`link`,`news_alert`.`created`,`news_alert`.`from`,`news_alert`.`to`,`File`.`file` AS `File__file`,`File`.`name` AS `File__name`,`File`.`path` AS `File__path`,`File`.`downloads` AS `File__downloads`,`File`.`url` AS `File__url`,`File`.`downloaded` AS `File__downloaded` FROM `news_alert` LEFT JOIN `files` `File` ON `news_alert`.`file` = `File`.`file` WHERE NOW() between `from` and `to` ORDER BY `news_alert`.`news_alert` LIMIT 1"

func (s *NewsAlertSuite) Test_GetTopNewsOne() {
	expectedAlert := model.NewsAlert{
		NewsAlert: 1,
		FileID:    3001,
		File: model.File{
			File:       3001,
			Name:       "Tournament_app_02-02.png",
			Path:       "newsalerts/",
			Downloads:  0,
			URL:        null.String{},
			Downloaded: null.Bool{},
		},
		Name:    null.StringFrom("Exzellenzuniversität"),
		Link:    null.StringFrom("https://tum.de"),
		Created: time.Time.Add(time.Now(), time.Hour*-4),
		From:    time.Time.Add(time.Now(), time.Hour*-2),
		To:      time.Time.Add(time.Now(), time.Hour*2),
	}
	s.mock.ExpectQuery(regexp.QuoteMeta(ExpectedGetTopNewsQuery)).
		WillReturnRows(sqlmock.NewRows([]string{"news_alert", "file", "name", "link", "created", "from", "to", "File__file", "File__name", "File__path", "File__downloads", "File__url", "File__downloaded"}).
			AddRow(expectedAlert.NewsAlert, expectedAlert.FileID, expectedAlert.Name, expectedAlert.Link, expectedAlert.Created, expectedAlert.From, expectedAlert.To, expectedAlert.File.File, expectedAlert.File.Name, expectedAlert.File.Path, expectedAlert.File.Downloads, expectedAlert.File.URL, expectedAlert.File.Downloaded))

	meta := metadata.MD{}
	server := CampusServer{db: s.DB, deviceBuf: s.deviceBuf}
	response, err := server.GetTopNews(metadata.NewIncomingContext(context.Background(), meta), nil)
	require.NoError(s.T(), err)
	require.Equal(s.T(), &pb.GetTopNewsReply{
		ImageUrl: expectedAlert.File.URL.String,
		Link:     expectedAlert.Link.String,
		Created:  timestamppb.New(expectedAlert.Created),
		From:     timestamppb.New(expectedAlert.From),
		To:       timestamppb.New(expectedAlert.To),
	}, response)
}
func (s *NewsAlertSuite) Test_GetTopNewsNone() {
	s.mock.ExpectQuery(regexp.QuoteMeta(ExpectedGetTopNewsQuery)).WillReturnError(gorm.ErrRecordNotFound)

	meta := metadata.MD{}
	server := CampusServer{db: s.DB, deviceBuf: s.deviceBuf}
	response, err := server.GetTopNews(metadata.NewIncomingContext(context.Background(), meta), nil)
	require.Equal(s.T(), status.Error(codes.NotFound, "no current active top news"), err)
	require.Nil(s.T(), response)
}
func (s *NewsAlertSuite) Test_GetTopNewsError() {
	s.mock.ExpectQuery(regexp.QuoteMeta(ExpectedGetTopNewsQuery)).WillReturnError(gorm.ErrInvalidDB)

	meta := metadata.MD{}
	server := CampusServer{db: s.DB, deviceBuf: s.deviceBuf}
	response, err := server.GetTopNews(metadata.NewIncomingContext(context.Background(), meta), nil)
	require.Equal(s.T(), status.Error(codes.Internal, "could not GetTopNews"), err)
	require.Nil(s.T(), response)
}

func (s *NewsAlertSuite) AfterTest(_, _ string) {
	require.NoError(s.T(), s.mock.ExpectationsWereMet())
}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestNewsAlertSuite(t *testing.T) {
	suite.Run(t, new(NewsAlertSuite))
}
