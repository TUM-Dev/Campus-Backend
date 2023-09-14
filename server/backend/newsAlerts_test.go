package backend

import (
	"context"
	"database/sql"
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
	"regexp"
	"testing"
	"time"
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

const ExpectedGetTopNewsQuery = "SELECT `news_alert`.`news_alert`,`news_alert`.`file`,`news_alert`.`name`,`news_alert`.`link`,`news_alert`.`created`,`news_alert`.`from`,`news_alert`.`to`,`Files`.`file` AS `Files__file`,`Files`.`name` AS `Files__name`,`Files`.`path` AS `Files__path`,`Files`.`downloads` AS `Files__downloads`,`Files`.`url` AS `Files__url`,`Files`.`downloaded` AS `Files__downloaded` FROM `news_alert` LEFT JOIN `files` `Files` ON `news_alert`.`file` = `Files`.`file` WHERE NOW() between `from` and `to` ORDER BY `news_alert`.`news_alert` LIMIT 1"

func (s *NewsAlertSuite) Test_GetTopNewsOne() {
	expectedAlert := model.NewsAlert{
		NewsAlert: 1,
		FilesID:   3001,
		Files: model.Files{
			File:       3001,
			Name:       "Tournament_app_02-02.png",
			Path:       "newsalerts/",
			Downloads:  0,
			URL:        sql.NullString{Valid: false},
			Downloaded: sql.NullBool{Bool: true, Valid: true},
		},
		Name:    null.String{NullString: sql.NullString{String: "Exzellenzuniversität", Valid: true}},
		Link:    null.String{NullString: sql.NullString{String: "https://tum.de", Valid: true}},
		Created: time.Time.Add(time.Now(), time.Hour*-4),
		From:    time.Time.Add(time.Now(), time.Hour*-2),
		To:      time.Time.Add(time.Now(), time.Hour*2),
	}
	s.mock.ExpectQuery(regexp.QuoteMeta(ExpectedGetTopNewsQuery)).
		WillReturnRows(sqlmock.NewRows([]string{"news_alert", "file", "name", "link", "created", "from", "to", "Files__file", "Files__name", "Files__path", "Files__downloads", "Files__url", "Files__downloaded"}).
			AddRow(expectedAlert.NewsAlert, expectedAlert.FilesID, expectedAlert.Name, expectedAlert.Link, expectedAlert.Created, expectedAlert.From, expectedAlert.To, expectedAlert.Files.File, expectedAlert.Files.Name, expectedAlert.Files.Path, expectedAlert.Files.Downloads, expectedAlert.Files.URL, expectedAlert.Files.Downloaded))

	meta := metadata.MD{}
	server := CampusServer{db: s.DB, deviceBuf: s.deviceBuf}
	response, err := server.GetTopNews(metadata.NewIncomingContext(context.Background(), meta), nil)
	require.NoError(s.T(), err)
	require.Equal(s.T(), &pb.GetTopNewsReply{
		ImageUrl: expectedAlert.Files.URL.String,
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
