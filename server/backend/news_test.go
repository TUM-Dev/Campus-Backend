package backend

import (
	"context"
	"database/sql"
	"fmt"
	"regexp"
	"testing"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/DATA-DOG/go-sqlmock"
	pb "github.com/TUM-Dev/Campus-Backend/server/api/tumdev"
	"github.com/TUM-Dev/Campus-Backend/server/model"
	"github.com/guregu/null"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc/metadata"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type NewsSuite struct {
	suite.Suite
	DB        *gorm.DB
	mock      sqlmock.Sqlmock
	deviceBuf *deviceBuffer
}

func (s *NewsSuite) SetupSuite() {
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

func newsFile(id int64) *model.Files {
	return &model.Files{
		File:       id,
		Name:       fmt.Sprintf("src_%d.png", id),
		Path:       "news/sources",
		Downloads:  1,
		URL:        null.String{},
		Downloaded: null.BoolFrom(true),
	}
}

func source1() *model.NewsSource {
	return &model.NewsSource{
		Source:  1,
		Title:   "Amazing News 1",
		URL:     null.StringFrom("https://example.com/amazing1"),
		FilesID: newsFile(1).File,
		Files:   *newsFile(1),
		Hook:    null.StringFrom(""),
	}
}

func source2() *model.NewsSource {
	return &model.NewsSource{
		Source:  2,
		Title:   "Amazing News 2",
		URL:     null.StringFrom("https://example.com/amazing2"),
		FilesID: newsFile(2).File,
		Files:   *newsFile(2),
		Hook:    null.StringFrom("hook"),
	}
}

const ExpectedGetSourceQuery = "SELECT `newsSource`.`source`,`newsSource`.`title`,`newsSource`.`url`,`newsSource`.`icon`,`newsSource`.`hook`,`Files`.`file` AS `Files__file`,`Files`.`name` AS `Files__name`,`Files`.`path` AS `Files__path`,`Files`.`downloads` AS `Files__downloads`,`Files`.`url` AS `Files__url`,`Files`.`downloaded` AS `Files__downloaded` FROM `newsSource` LEFT JOIN `files` `Files` ON `newsSource`.`icon` = `Files`.`file`"

func (s *NewsSuite) Test_GetNewsSourcesMultiple() {
	s.mock.ExpectQuery(regexp.QuoteMeta(ExpectedGetSourceQuery)).
		WillReturnRows(sqlmock.NewRows([]string{"source", "title", "url", "icon", "hook", "Files__file", "Files__name", "Files__path", "Files__downloads", "Files__url", "Files__downloaded"}).
			AddRow(source1().Source, source1().Title, source1().URL, source1().FilesID, source1().Hook, source1().Files.File, source1().Files.Name, source1().Files.Path, source1().Files.Downloads, source1().Files.URL, source1().Files.Downloaded).
			AddRow(source2().Source, source2().Title, source2().URL, source2().FilesID, source2().Hook, source2().Files.File, source2().Files.Name, source2().Files.Path, source2().Files.Downloads, source2().Files.URL, source2().Files.Downloaded))

	meta := metadata.MD{}
	server := CampusServer{db: s.DB, deviceBuf: s.deviceBuf}
	response, err := server.GetNewsSources(metadata.NewIncomingContext(context.Background(), meta), nil)
	require.NoError(s.T(), err)
	expectedResp := &pb.GetNewsSourcesReply{
		Sources: []*pb.NewsSource{
			{Source: fmt.Sprintf("%d", source1().Source), Title: source1().Title, Icon: source1().Files.URL.String},
			{Source: fmt.Sprintf("%d", source2().Source), Title: source2().Title, Icon: source2().Files.URL.String},
		},
	}
	require.Equal(s.T(), expectedResp, response)
}

func news1() *model.News {
	return &model.News{
		News:    1,
		Title:   "Amazing News 1",
		Link:    "https://example.com/amazing2",
		FilesID: null.IntFrom(newsFile(1).File),
		Files:   newsFile(1),
	}
}

func news2() *model.News {
	return &model.News{
		News:    2,
		Title:   "Amazing News 2",
		Link:    "https://example.com/amazing2",
		FilesID: null.Int{},
		Files:   nil,
	}
}

func (s *NewsSuite) Test_GetNewsSourcesNone() {
	s.mock.ExpectQuery(regexp.QuoteMeta(ExpectedGetSourceQuery)).
		WillReturnRows(sqlmock.NewRows([]string{"source", "title", "url", "icon", "hook", "Files__file", "Files__name", "Files__path", "Files__downloads", "Files__url", "Files__downloaded"}))

	meta := metadata.MD{}
	server := CampusServer{db: s.DB, deviceBuf: s.deviceBuf}
	response, err := server.GetNewsSources(metadata.NewIncomingContext(context.Background(), meta), nil)
	require.NoError(s.T(), err)
	expectedResp := &pb.GetNewsSourcesReply{
		Sources: []*pb.NewsSource(nil),
	}
	require.Equal(s.T(), expectedResp, response)
}

const ExpectedGetNewsQuery = "SELECT `news`.`news`,`news`.`date`,`news`.`created`,`news`.`title`,`news`.`description`,`news`.`src`,`news`.`link`,`news`.`image`,`news`.`file`,`Files`.`file` AS `Files__file`,`Files`.`name` AS `Files__name`,`Files`.`path` AS `Files__path`,`Files`.`downloads` AS `Files__downloads`,`Files`.`url` AS `Files__url`,`Files`.`downloaded` AS `Files__downloaded` FROM `news` LEFT JOIN `files` `Files` ON `news`.`file` = `Files`.`file`"

func (s *NewsSuite) Test_GetNewsNone_withFilters() {
	s.mock.ExpectQuery(regexp.QuoteMeta(ExpectedGetNewsQuery+" WHERE src = ? AND news > ?")).
		WithArgs(1, 2).
		WillReturnRows(sqlmock.NewRows([]string{"news", "date", "created", "title", "description", "src", "link", "image", "file", "Files__file", "Files__name", "Files__path", "Files__downloads", "Files__url", "Files__downloaded"}))

	meta := metadata.NewIncomingContext(context.Background(), metadata.MD{})
	server := CampusServer{db: s.DB, deviceBuf: s.deviceBuf}
	response, err := server.GetNews(meta, &pb.GetNewsRequest{NewsSource: 1, LastNewsId: 2})
	require.NoError(s.T(), err)
	expectedResp := &pb.GetNewsReply{
		News: []*pb.NewsItem{},
	}
	require.Equal(s.T(), expectedResp, response)
}
func (s *NewsSuite) Test_GetNewsNone() {
	s.mock.ExpectQuery(regexp.QuoteMeta(ExpectedGetNewsQuery)).
		WillReturnRows(sqlmock.NewRows([]string{"news", "date", "created", "title", "description", "src", "link", "image", "file", "Files__file", "Files__name", "Files__path", "Files__downloads", "Files__url", "Files__downloaded"}))

	meta := metadata.NewIncomingContext(context.Background(), metadata.MD{})
	server := CampusServer{db: s.DB, deviceBuf: s.deviceBuf}
	response, err := server.GetNews(meta, &pb.GetNewsRequest{})
	require.NoError(s.T(), err)
	expectedResp := &pb.GetNewsReply{
		News: []*pb.NewsItem{},
	}
	require.Equal(s.T(), expectedResp, response)
}
func (s *NewsSuite) Test_GetNewsMultiple() {
	n1 := news1()
	n2 := news2()
	s.mock.ExpectQuery(regexp.QuoteMeta(" ")).
		WillReturnRows(sqlmock.NewRows([]string{"news", "date", "created", "title", "description", "src", "link", "image", "file", "Files__file", "Files__name", "Files__path", "Files__downloads", "Files__url", "Files__downloaded"}).
			AddRow(n1.News, n1.Date, n1.Created, n1.Title, n1.Description, n1.Src, n1.Link, n1.Image, n1.FilesID, n1.Files.File, n1.Files.Name, n1.Files.Path, n1.Files.Downloads, n1.Files.URL, n1.Files.Downloaded).
			AddRow(n2.News, n2.Date, n2.Created, n2.Title, n2.Description, n2.Src, n2.Link, n2.Image, nil, nil, nil, nil, nil, nil, nil))

	meta := metadata.NewIncomingContext(context.Background(), metadata.MD{})
	server := CampusServer{db: s.DB, deviceBuf: s.deviceBuf}
	response, err := server.GetNews(meta, &pb.GetNewsRequest{})
	require.NoError(s.T(), err)
	expectedResp := &pb.GetNewsReply{
		News: []*pb.NewsItem{
			{Id: n1.News, Title: n1.Title, Text: n1.Description, Link: n1.Link, ImageUrl: n1.Image.String, Source: fmt.Sprintf("%d", n1.Src), Created: timestamppb.New(n1.Created)},
			{Id: n2.News, Title: n2.Title, Text: n2.Description, Link: n2.Link, ImageUrl: n2.Image.String, Source: fmt.Sprintf("%d", n2.Src), Created: timestamppb.New(n2.Created)},
		},
	}
	require.Equal(s.T(), expectedResp, response)
}

func newsAlertFile(id int32) *model.Files {
	return &model.Files{
		File:       id,
		Name:       fmt.Sprintf("src_%d.png", id),
		Path:       "news/sources",
		Downloads:  1,
		URL:        null.String{},
		Downloaded: null.BoolFrom(true),
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

const ExpectedGetNewsAlertsQuery = "SELECT `news_alert`.`news_alert`,`news_alert`.`file`,`news_alert`.`name`,`news_alert`.`link`,`news_alert`.`created`,`news_alert`.`from`,`news_alert`.`to`,`Files`.`file` AS `Files__file`,`Files`.`name` AS `Files__name`,`Files`.`path` AS `Files__path`,`Files`.`downloads` AS `Files__downloads`,`Files`.`url` AS `Files__url`,`Files`.`downloaded` AS `Files__downloaded` FROM `news_alert` LEFT JOIN `files` `Files` ON `news_alert`.`file` = `Files`.`file` WHERE news_alert.to >= NOW()"

func (s *NewsSuite) Test_GetNewsAlertsError() {
	s.mock.ExpectQuery(regexp.QuoteMeta(ExpectedGetNewsAlertsQuery)).WillReturnError(gorm.ErrInvalidDB)

	meta := metadata.MD{}
	server := CampusServer{db: s.DB, deviceBuf: s.deviceBuf}
	response, err := server.GetNewsAlerts(metadata.NewIncomingContext(context.Background(), meta), &pb.GetNewsAlertsRequest{})
	require.Equal(s.T(), status.Error(codes.Internal, "could not GetNewsAlerts"), err)
	require.Nil(s.T(), response)
}
func (s *NewsSuite) Test_GetNewsAlertsNone_noFilter() {
	s.mock.ExpectQuery(regexp.QuoteMeta(ExpectedGetNewsAlertsQuery)).WillReturnError(gorm.ErrRecordNotFound)

	server := CampusServer{db: s.DB, deviceBuf: s.deviceBuf}
	response, err := server.GetNewsAlerts(metadata.NewIncomingContext(context.Background(), metadata.MD{}), &pb.GetNewsAlertsRequest{})
	require.Equal(s.T(), status.Error(codes.NotFound, "no news alerts"), err)
	require.Nil(s.T(), response)
}
func (s *NewsSuite) Test_GetNewsAlertsNone_Filter() {
	s.mock.ExpectQuery(regexp.QuoteMeta(ExpectedGetNewsAlertsQuery + " AND news_alert.alert > ?")).WithArgs(42).WillReturnError(gorm.ErrRecordNotFound)

	server := CampusServer{db: s.DB, deviceBuf: s.deviceBuf}
	response, err := server.GetNewsAlerts(metadata.NewIncomingContext(context.Background(), metadata.MD{}), &pb.GetNewsAlertsRequest{LastNewsAlertId: 42})
	require.Equal(s.T(), status.Error(codes.NotFound, "no news alerts"), err)
	require.Nil(s.T(), response)
}
func (s *NewsSuite) Test_GetNewsAlertsMultiple() {
	a1 := alert1()
	a2 := alert2()
	s.mock.ExpectQuery(regexp.QuoteMeta(ExpectedGetNewsAlertsQuery)).
		WillReturnRows(sqlmock.NewRows([]string{"news_alert", "file", "name", "link", "created", "from", "to", "Files__file", "Files__name", "Files__path", "Files__downloads", "Files__url", "Files__downloaded"}).
			AddRow(a1.NewsAlert, a1.FilesID, a1.Name, a1.Link, a1.Created, a1.From, a1.To, a1.Files.File, a1.Files.Name, a1.Files.Path, a1.Files.Downloads, a1.Files.URL, a1.Files.Downloaded).
			AddRow(a2.NewsAlert, a2.FilesID, a2.Name, a2.Link, a2.Created, a2.From, a2.To, a2.Files.File, a2.Files.Name, a2.Files.Path, a2.Files.Downloads, a2.Files.URL, a2.Files.Downloaded))

	server := CampusServer{db: s.DB, deviceBuf: s.deviceBuf}
	response, err := server.GetNewsAlerts(metadata.NewIncomingContext(context.Background(), metadata.MD{}), &pb.GetNewsAlertsRequest{})
	require.NoError(s.T(), err)
	expectedResp := &pb.GetNewsAlertsReply{
		Alerts: []*pb.NewsAlert{
			{ImageUrl: a1.Files.URL.String, Link: a1.Link.String, Created: timestamppb.New(a1.Created), From: timestamppb.New(a1.From), To: timestamppb.New(a1.To)},
			{ImageUrl: a2.Files.URL.String, Link: a2.Link.String, Created: timestamppb.New(a2.Created), From: timestamppb.New(a2.From), To: timestamppb.New(a2.To)},
		}}
	require.Equal(s.T(), expectedResp, response)
}

func (s *NewsSuite) AfterTest(_, _ string) {
	require.NoError(s.T(), s.mock.ExpectationsWereMet())
}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestNewsSuite(t *testing.T) {
	suite.Run(t, new(NewsSuite))
}
