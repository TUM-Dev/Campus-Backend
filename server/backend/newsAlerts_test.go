package backend

import (
	"context"
	"database/sql"
	"fmt"
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

const ExpectedGetTopNewsAlertQuery = "SELECT `news_alert`.`news_alert`,`news_alert`.`file`,`news_alert`.`name`,`news_alert`.`link`,`news_alert`.`created`,`news_alert`.`from`,`news_alert`.`to`,`Files`.`file` AS `Files__file`,`Files`.`name` AS `Files__name`,`Files`.`path` AS `Files__path`,`Files`.`downloads` AS `Files__downloads`,`Files`.`url` AS `Files__url`,`Files`.`downloaded` AS `Files__downloaded` FROM `news_alert` LEFT JOIN `files` `Files` ON `news_alert`.`file` = `Files`.`file` WHERE NOW() between `from` and `to` ORDER BY `news_alert`.`news_alert` LIMIT 1"

func (s *NewsAlertSuite) Test_GetTopNewsAlertOne() {
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
	s.mock.ExpectQuery(regexp.QuoteMeta(ExpectedGetTopNewsAlertQuery)).
		WillReturnRows(sqlmock.NewRows([]string{"news_alert", "file", "name", "link", "created", "from", "to", "Files__file", "Files__name", "Files__path", "Files__downloads", "Files__url", "Files__downloaded"}).
			AddRow(expectedAlert.NewsAlert, expectedAlert.FilesID, expectedAlert.Name, expectedAlert.Link, expectedAlert.Created, expectedAlert.From, expectedAlert.To, expectedAlert.Files.File, expectedAlert.Files.Name, expectedAlert.Files.Path, expectedAlert.Files.Downloads, expectedAlert.Files.URL, expectedAlert.Files.Downloaded))

	server := CampusServer{db: s.DB, deviceBuf: s.deviceBuf}
	response, err := server.GetTopNewsAlert(metadata.NewIncomingContext(context.Background(), metadata.MD{}), nil)
	require.NoError(s.T(), err)
	require.Equal(s.T(), &pb.GetTopNewsAlertReply{Alert: &pb.NewsAlert{
		ImageUrl: expectedAlert.Files.URL.String,
		Link:     expectedAlert.Link.String,
		Created:  timestamppb.New(expectedAlert.Created),
		From:     timestamppb.New(expectedAlert.From),
		To:       timestamppb.New(expectedAlert.To)}}, response)
}
func (s *NewsAlertSuite) Test_GetTopNewsAlertNone() {
	s.mock.ExpectQuery(regexp.QuoteMeta(ExpectedGetTopNewsAlertQuery)).WillReturnError(gorm.ErrRecordNotFound)

	server := CampusServer{db: s.DB, deviceBuf: s.deviceBuf}
	response, err := server.GetTopNewsAlert(metadata.NewIncomingContext(context.Background(), metadata.MD{}), nil)
	require.Equal(s.T(), status.Error(codes.NotFound, "no current active top news"), err)
	require.Nil(s.T(), response)
}
func (s *NewsAlertSuite) Test_GetTopNewsAlertError() {
	s.mock.ExpectQuery(regexp.QuoteMeta(ExpectedGetTopNewsAlertQuery)).WillReturnError(gorm.ErrInvalidDB)

	server := CampusServer{db: s.DB, deviceBuf: s.deviceBuf}
	response, err := server.GetTopNewsAlert(metadata.NewIncomingContext(context.Background(), metadata.MD{}), nil)
	require.Equal(s.T(), status.Error(codes.Internal, "could not GetTopNewsAlert"), err)
	require.Nil(s.T(), response)
}

func newsAlertFile(id int32) *model.Files {
	return &model.Files{
		File:       id,
		Name:       fmt.Sprintf("src_%d.png", id),
		Path:       "news/sources",
		Downloads:  1,
		URL:        sql.NullString{Valid: false},
		Downloaded: sql.NullBool{Bool: true, Valid: true},
	}
}
func alert1() *model.NewsAlert {
	return &model.NewsAlert{
		NewsAlert: 1,
		FilesID:   newsAlertFile(1).File,
		Files:     *newsAlertFile(1),
		Name:      null.String{},
		Link:      null.String{},
		Created:   time.Time.Add(time.Now(), time.Hour*-4),
		From:      time.Time.Add(time.Now(), time.Hour*-2),
		To:        time.Time.Add(time.Now(), time.Hour*-2),
	}
}

func alert2() *model.NewsAlert {
	return &model.NewsAlert{
		NewsAlert: 2,
		FilesID:   newsAlertFile(1).File,
		Files:     *newsAlertFile(1),
		Name:      null.String{},
		Link:      null.String{},
		Created:   time.Time.Add(time.Now(), time.Hour),
		From:      time.Time.Add(time.Now(), time.Hour*2),
		To:        time.Time.Add(time.Now(), time.Hour*3),
	}
}

const ExpectedGetNewsAlertsQuery = "SELECT `news_alert`.`news_alert`,`news_alert`.`file`,`news_alert`.`name`,`news_alert`.`link`,`news_alert`.`created`,`news_alert`.`from`,`news_alert`.`to`,`Files`.`file` AS `Files__file`,`Files`.`name` AS `Files__name`,`Files`.`path` AS `Files__path`,`Files`.`downloads` AS `Files__downloads`,`Files`.`url` AS `Files__url`,`Files`.`downloaded` AS `Files__downloaded` FROM `news_alert` LEFT JOIN `files` `Files` ON `news_alert`.`file` = `Files`.`file`"

func (s *NewsAlertSuite) Test_GetNewsAlertsError() {
	s.mock.ExpectQuery(regexp.QuoteMeta(ExpectedGetNewsAlertsQuery)).WillReturnError(gorm.ErrInvalidDB)

	meta := metadata.MD{}
	server := CampusServer{db: s.DB, deviceBuf: s.deviceBuf}
	response, err := server.GetNewsAlerts(metadata.NewIncomingContext(context.Background(), meta), nil)
	require.Equal(s.T(), status.Error(codes.Internal, "could not GetNewsAlerts"), err)
	require.Nil(s.T(), response)
}
func (s *NewsAlertSuite) Test_GetNewsAlertsNone() {
	s.mock.ExpectQuery(regexp.QuoteMeta(ExpectedGetNewsAlertsQuery)).WillReturnError(gorm.ErrRecordNotFound)

	server := CampusServer{db: s.DB, deviceBuf: s.deviceBuf}
	response, err := server.GetNewsAlerts(metadata.NewIncomingContext(context.Background(), metadata.MD{}), nil)
	require.Equal(s.T(), status.Error(codes.NotFound, "no news alerts"), err)
	require.Nil(s.T(), response)
}
func (s *NewsAlertSuite) Test_GetNewsAlertsMultiple() {
	a1 := alert1()
	a2 := alert2()
	s.mock.ExpectQuery(regexp.QuoteMeta(ExpectedGetNewsAlertsQuery)).
		WillReturnRows(sqlmock.NewRows([]string{"news_alert", "file", "name", "link", "created", "from", "to", "Files__file", "Files__name", "Files__path", "Files__downloads", "Files__url", "Files__downloaded"}).
			AddRow(a1.NewsAlert, a1.FilesID, a1.Name, a1.Link, a1.Created, a1.From, a1.To, a1.Files.File, a1.Files.Name, a1.Files.Path, a1.Files.Downloads, a1.Files.URL, a1.Files.Downloaded).
			AddRow(a2.NewsAlert, a2.FilesID, a2.Name, a2.Link, a2.Created, a2.From, a2.To, a2.Files.File, a2.Files.Name, a2.Files.Path, a2.Files.Downloads, a2.Files.URL, a2.Files.Downloaded))

	server := CampusServer{db: s.DB, deviceBuf: s.deviceBuf}
	response, err := server.GetNewsAlerts(metadata.NewIncomingContext(context.Background(), metadata.MD{}), nil)
	require.NoError(s.T(), err)
	expectedResp := &pb.GetNewsAlertsReply{
		Alerts: []*pb.NewsAlert{
			{ImageUrl: a1.Files.URL.String, Link: a1.Link.String, Created: timestamppb.New(a1.Created), From: timestamppb.New(a1.From), To: timestamppb.New(a1.To)},
			{ImageUrl: a2.Files.URL.String, Link: a2.Link.String, Created: timestamppb.New(a2.Created), From: timestamppb.New(a2.From), To: timestamppb.New(a2.To)},
		}}
	require.Equal(s.T(), expectedResp, response)
}

const ExpectedGetNewsAlertQuery = "SELECT `news_alert`.`news_alert`,`news_alert`.`file`,`news_alert`.`name`,`news_alert`.`link`,`news_alert`.`created`,`news_alert`.`from`,`news_alert`.`to`,`Files`.`file` AS `Files__file`,`Files`.`name` AS `Files__name`,`Files`.`path` AS `Files__path`,`Files`.`downloads` AS `Files__downloads`,`Files`.`url` AS `Files__url`,`Files`.`downloaded` AS `Files__downloaded` FROM `news_alert` LEFT JOIN `files` `Files` ON `news_alert`.`file` = `Files`.`file` WHERE `news_alert`.`news_alert` = ? ORDER BY `news_alert`.`news_alert` LIMIT 1"

func (s *NewsAlertSuite) Test_GetNewsAlertError() {
	s.mock.ExpectQuery(regexp.QuoteMeta(ExpectedGetNewsAlertQuery)).WillReturnError(gorm.ErrInvalidDB)

	meta := metadata.MD{}
	server := CampusServer{db: s.DB, deviceBuf: s.deviceBuf}
	response, err := server.GetNewsAlert(metadata.NewIncomingContext(context.Background(), meta), &pb.GetNewsAlertRequest{Id: 42})
	require.Equal(s.T(), status.Error(codes.Internal, "could not GetNewsAlert"), err)
	require.Nil(s.T(), response)
}
func (s *NewsAlertSuite) Test_GetNewsAlertNone() {
	s.mock.ExpectQuery(regexp.QuoteMeta(ExpectedGetNewsAlertQuery)).WillReturnError(gorm.ErrRecordNotFound)

	server := CampusServer{db: s.DB, deviceBuf: s.deviceBuf}
	response, err := server.GetNewsAlert(metadata.NewIncomingContext(context.Background(), metadata.MD{}), &pb.GetNewsAlertRequest{Id: 42})
	require.Equal(s.T(), status.Error(codes.NotFound, "no news alert"), err)
	require.Nil(s.T(), response)
}
func (s *NewsAlertSuite) Test_GetNewsAlertOne() {
	a1 := alert1()
	s.mock.ExpectQuery(regexp.QuoteMeta(ExpectedGetNewsAlertQuery)).
		WillReturnRows(sqlmock.NewRows([]string{"news_alert", "file", "name", "link", "created", "from", "to", "Files__file", "Files__name", "Files__path", "Files__downloads", "Files__url", "Files__downloaded"}).
			AddRow(a1.NewsAlert, a1.FilesID, a1.Name, a1.Link, a1.Created, a1.From, a1.To, a1.Files.File, a1.Files.Name, a1.Files.Path, a1.Files.Downloads, a1.Files.URL, a1.Files.Downloaded))

	server := CampusServer{db: s.DB, deviceBuf: s.deviceBuf}
	response, err := server.GetNewsAlert(metadata.NewIncomingContext(context.Background(), metadata.MD{}), &pb.GetNewsAlertRequest{Id: 1})
	require.NoError(s.T(), err)
	expectedResp := &pb.GetNewsAlertReply{
		Alert: &pb.NewsAlert{ImageUrl: a1.Files.URL.String, Link: a1.Link.String, Created: timestamppb.New(a1.Created), From: timestamppb.New(a1.From), To: timestamppb.New(a1.To)},
	}
	require.Equal(s.T(), expectedResp, response)
}

func (s *NewsAlertSuite) AfterTest(_, _ string) {
	require.NoError(s.T(), s.mock.ExpectationsWereMet())
}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestNewsAlertSuite(t *testing.T) {
	suite.Run(t, new(NewsAlertSuite))
}
