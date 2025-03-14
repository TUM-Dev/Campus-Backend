package backend

import (
	"context"
	"database/sql"
	"fmt"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	pb "github.com/TUM-Dev/Campus-Backend/server/api/tumdev"
	"github.com/TUM-Dev/Campus-Backend/server/model"
	"github.com/guregu/null"
	"github.com/hashicorp/golang-lru/v2/expirable"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
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
		WillReturnRows(sqlmock.NewRows([]string{"VERSION()"}).AddRow("8.4.0"))
	s.DB, err = gorm.Open(dialector, &gorm.Config{})
	require.NoError(s.T(), err)

	s.deviceBuf = newDeviceBuffer()
}

func newsFile(id int64) *model.File {
	return &model.File{
		File:       id,
		Name:       fmt.Sprintf("src_%d.png", id),
		Path:       "news/sources/",
		Downloads:  1,
		URL:        null.String{},
		Downloaded: null.BoolFrom(true),
	}
}

func source1() *model.NewsSource {
	return &model.NewsSource{
		Source: 1,
		Title:  "Amazing News 1",
		URL:    null.StringFrom("https://example.com/amazing1"),
		FileID: newsFile(1).File,
		File:   *newsFile(1),
		Hook:   null.StringFrom(""),
	}
}

func source2() *model.NewsSource {
	return &model.NewsSource{
		Source: 2,
		Title:  "Amazing News 2",
		URL:    null.StringFrom("https://example.com/amazing2"),
		FileID: newsFile(2).File,
		File:   *newsFile(2),
		Hook:   null.StringFrom("hook"),
	}
}

const ExpectedListNewsSourcesQuery = "SELECT `news_sources`.`source`,`news_sources`.`title`,`news_sources`.`url`,`news_sources`.`icon`,`news_sources`.`hook`,`File`.`file` AS `File__file`,`File`.`name` AS `File__name`,`File`.`path` AS `File__path`,`File`.`downloads` AS `File__downloads`,`File`.`url` AS `File__url`,`File`.`downloaded` AS `File__downloaded` FROM `news_sources` LEFT JOIN `files` `File` ON `news_sources`.`icon` = `File`.`file`"

func (s *NewsSuite) Test_ListNewsSourcesMultiple() {
	s1, s2 := source1(), source2()
	s.mock.ExpectQuery(regexp.QuoteMeta(ExpectedListNewsSourcesQuery)).
		WillReturnRows(sqlmock.NewRows([]string{"source", "title", "url", "icon", "hook", "File__file", "File__name", "File__path", "File__downloads", "File__url", "File__downloaded"}).
			AddRow(s1.Source, s1.Title, s1.URL, s1.FileID, s1.Hook, s1.File.File, s1.File.Name, s1.File.Path, s1.File.Downloads, s1.File.URL, s1.File.Downloaded).
			AddRow(s2.Source, s2.Title, s2.URL, s2.FileID, s2.Hook, s2.File.File, s2.File.Name, s2.File.Path, s2.File.Downloads, s2.File.URL, s2.File.Downloaded))

	meta := metadata.MD{}
	server := s.getCampusTestServer()
	response, err := server.ListNewsSources(metadata.NewIncomingContext(context.Background(), meta), nil)
	require.NoError(s.T(), err)
	expectedResp := &pb.ListNewsSourcesReply{
		Sources: []*pb.NewsSource{
			{Source: fmt.Sprintf("%d", s1.Source), Title: s1.Title, IconUrl: "https://api.tum.app/files/news/sources/src_1.png"},
			{Source: fmt.Sprintf("%d", s2.Source), Title: s2.Title, IconUrl: "https://api.tum.app/files/news/sources/src_2.png"},
		},
	}
	require.Equal(s.T(), expectedResp, response)
}

func news1() *model.News {
	return &model.News{
		News:   1,
		Title:  "Amazing News 1",
		Link:   "https://example.com/amazing2",
		FileID: null.IntFrom(newsFile(1).File),
		File:   newsFile(1),
	}
}

func news2() *model.News {
	return &model.News{
		News:   2,
		Title:  "Amazing News 2",
		Link:   "https://example.com/amazing2",
		FileID: null.Int{},
		File:   nil,
	}
}

func (s *NewsSuite) Test_ListNewsSourcesNone() {
	s.mock.ExpectQuery(regexp.QuoteMeta(ExpectedListNewsSourcesQuery)).
		WillReturnRows(sqlmock.NewRows([]string{"source", "title", "url", "icon", "hook", "File__file", "File__name", "File__path", "File__downloads", "File__url", "File__downloaded"}))

	meta := metadata.MD{}
	server := s.getCampusTestServer()
	response, err := server.ListNewsSources(metadata.NewIncomingContext(context.Background(), meta), nil)
	require.NoError(s.T(), err)
	expectedResp := &pb.ListNewsSourcesReply{
		Sources: []*pb.NewsSource(nil),
	}
	require.Equal(s.T(), expectedResp, response)
}

const ExpectedListNewsQuery = "SELECT `news`.`news`,`news`.`date`,`news`.`created`,`news`.`title`,`news`.`description`,`news`.`src`,`news`.`link`,`news`.`image`,`news`.`file`,`File`.`file` AS `File__file`,`File`.`name` AS `File__name`,`File`.`path` AS `File__path`,`File`.`downloads` AS `File__downloads`,`File`.`url` AS `File__url`,`File`.`downloaded` AS `File__downloaded`,`NewsSource`.`source` AS `NewsSource__source`,`NewsSource`.`title` AS `NewsSource__title`,`NewsSource`.`url` AS `NewsSource__url`,`NewsSource`.`icon` AS `NewsSource__icon`,`NewsSource`.`hook` AS `NewsSource__hook`,`NewsSource__File`.`file` AS `NewsSource__File__file`,`NewsSource__File`.`name` AS `NewsSource__File__name`,`NewsSource__File`.`path` AS `NewsSource__File__path`,`NewsSource__File`.`downloads` AS `NewsSource__File__downloads`,`NewsSource__File`.`url` AS `NewsSource__File__url`,`NewsSource__File`.`downloaded` AS `NewsSource__File__downloaded` FROM `news` LEFT JOIN `files` `File` ON `news`.`file` = `File`.`file` LEFT JOIN `news_sources` `NewsSource` ON `news`.`src` = `NewsSource`.`source` LEFT JOIN `files` `NewsSource__File` ON `NewsSource`.`icon` = `NewsSource__File`.`file`"

func (s *NewsSuite) Test_ListNewsNone_withFilters() {
	s.mock.ExpectQuery(regexp.QuoteMeta(ExpectedListNewsQuery+" WHERE src = ? AND news > ?")).
		WithArgs(1, 2).
		WillReturnRows(sqlmock.NewRows([]string{"news", "date", "created", "title", "description", "src", "link", "image", "file", "File__file", "File__name", "File__path", "File__downloads", "File__url", "File__downloaded", "source", "NewsSource__source", "NewsSource__title", "NewsSource__url", "NewsSource__icon", "NewsSource__hook", "NewsSource__File__file", "NewsSource__File__name", "NewsSource__File__path", "NewsSource__File__downloads", "NewsSource__File__url", "NewsSource__File__downloaded"}))

	meta := metadata.NewIncomingContext(context.Background(), metadata.MD{})
	server := s.getCampusTestServer()
	response, err := server.ListNews(meta, &pb.ListNewsRequest{NewsSource: 1, LastNewsId: 2})
	require.NoError(s.T(), err)
	require.Equal(s.T(), &pb.ListNewsReply{News: nil}, response)
}
func (s *NewsSuite) Test_ListNewsNone() {
	s.mock.ExpectQuery(regexp.QuoteMeta(ExpectedListNewsQuery)).
		WillReturnRows(sqlmock.NewRows([]string{"news", "date", "created", "title", "description", "src", "link", "image", "file", "File__file", "File__name", "File__path", "File__downloads", "File__url", "File__downloaded", "source", "NewsSource__source", "NewsSource__title", "NewsSource__url", "NewsSource__icon", "NewsSource__hook", "NewsSource__File__file", "NewsSource__File__name", "NewsSource__File__path", "NewsSource__File__downloads", "NewsSource__File__url", "NewsSource__File__downloaded"}))

	meta := metadata.NewIncomingContext(context.Background(), metadata.MD{})
	server := s.getCampusTestServer()
	response, err := server.ListNews(meta, &pb.ListNewsRequest{})
	require.NoError(s.T(), err)
	require.Equal(s.T(), &pb.ListNewsReply{News: nil}, response)
}
func (s *NewsSuite) Test_ListNewsMultiple() {
	n1 := news1()
	n2 := news2()
	s.mock.ExpectQuery(regexp.QuoteMeta(" ")).
		WillReturnRows(sqlmock.NewRows([]string{"news", "date", "created", "title", "description", "src", "link", "image", "file", "File__file", "File__name", "File__path", "File__downloads", "File__url", "File__downloaded"}).
			AddRow(n1.News, n1.Date, n1.Created, n1.Title, n1.Description, n1.NewsSourceID, n1.Link, n1.Image, n1.FileID, n1.File.File, n1.File.Name, n1.File.Path, n1.File.Downloads, n1.File.URL, n1.File.Downloaded).
			AddRow(n2.News, n2.Date, n2.Created, n2.Title, n2.Description, n2.NewsSourceID, n2.Link, n2.Image, nil, nil, nil, nil, nil, nil, nil))

	meta := metadata.NewIncomingContext(context.Background(), metadata.MD{})
	server := s.getCampusTestServer()
	response, err := server.ListNews(meta, &pb.ListNewsRequest{})
	require.NoError(s.T(), err)
	expectedResp := &pb.ListNewsReply{
		News: []*pb.News{
			{Id: n1.News, Title: n1.Title, Text: n1.Description, Link: n1.Link, ImageUrl: "https://api.tum.app/files/news/sources/src_1.png", SourceId: fmt.Sprintf("%d", n1.NewsSourceID), SourceIconUrl: n1.NewsSource.File.FullExternalUrl(), SourceTitle: n1.NewsSource.Title, Created: timestamppb.New(n1.Created), Date: timestamppb.New(n1.Date)},
			{Id: n2.News, Title: n2.Title, Text: n2.Description, Link: n2.Link, ImageUrl: "", SourceId: fmt.Sprintf("%d", n2.NewsSourceID), SourceIconUrl: n2.NewsSource.File.FullExternalUrl(), SourceTitle: n2.NewsSource.Title, Created: timestamppb.New(n2.Created), Date: timestamppb.New(n2.Date)},
		},
	}
	require.Equal(s.T(), expectedResp, response)
}

func newsAlertFile(id int64) *model.File {
	return &model.File{
		File:       id,
		Name:       fmt.Sprintf("src_%d.png", id),
		Path:       "news/sources/",
		Downloads:  1,
		URL:        null.String{},
		Downloaded: null.BoolFrom(true),
	}
}
func alert1() *model.NewsAlert {
	return &model.NewsAlert{
		NewsAlert: 1,
		FileID:    newsAlertFile(1).File,
		File:      *newsAlertFile(1),
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
		FileID:    newsAlertFile(2).File,
		File:      *newsAlertFile(2),
		Name:      null.String{},
		Link:      null.String{},
		Created:   time.Time.Add(time.Now(), time.Hour),
		From:      time.Time.Add(time.Now(), time.Hour*2),
		To:        time.Time.Add(time.Now(), time.Hour*3),
	}
}

const ExpectedListNewsAlertsQuery = "SELECT `news_alerts`.`news_alert`,`news_alerts`.`file`,`news_alerts`.`name`,`news_alerts`.`link`,`news_alerts`.`created`,`news_alerts`.`from`,`news_alerts`.`to`,`File`.`file` AS `File__file`,`File`.`name` AS `File__name`,`File`.`path` AS `File__path`,`File`.`downloads` AS `File__downloads`,`File`.`url` AS `File__url`,`File`.`downloaded` AS `File__downloaded` FROM `news_alerts` LEFT JOIN `files` `File` ON `news_alerts`.`file` = `File`.`file` WHERE news_alerts.to >= NOW()"

func (s *NewsSuite) Test_ListNewsAlertsError() {
	s.mock.ExpectQuery(regexp.QuoteMeta(ExpectedListNewsAlertsQuery)).WillReturnError(gorm.ErrInvalidDB)

	meta := metadata.MD{}
	server := s.getCampusTestServer()
	response, err := server.ListNewsAlerts(metadata.NewIncomingContext(context.Background(), meta), &pb.ListNewsAlertsRequest{})
	require.Equal(s.T(), status.Error(codes.Internal, "could not ListNewsAlerts"), err)
	require.Nil(s.T(), response)
}
func (s *NewsSuite) Test_ListNewsAlertsNone_noFilter() {
	s.mock.ExpectQuery(regexp.QuoteMeta(ExpectedListNewsAlertsQuery)).WillReturnError(gorm.ErrRecordNotFound)

	server := s.getCampusTestServer()
	response, err := server.ListNewsAlerts(metadata.NewIncomingContext(context.Background(), metadata.MD{}), &pb.ListNewsAlertsRequest{})
	require.Equal(s.T(), status.Error(codes.NotFound, "no news alerts"), err)
	require.Nil(s.T(), response)
}
func (s *NewsSuite) Test_ListNewsAlertsNone_Filter() {
	s.mock.ExpectQuery(regexp.QuoteMeta(ExpectedListNewsAlertsQuery + " AND news_alerts.news_alert > ?")).WithArgs(42).WillReturnError(gorm.ErrRecordNotFound)

	server := s.getCampusTestServer()
	response, err := server.ListNewsAlerts(metadata.NewIncomingContext(context.Background(), metadata.MD{}), &pb.ListNewsAlertsRequest{LastNewsAlertId: 42})
	require.Equal(s.T(), status.Error(codes.NotFound, "no news alerts"), err)
	require.Nil(s.T(), response)
}
func (s *NewsSuite) Test_ListNewsAlertsMultiple() {
	a1, a2 := alert1(), alert2()
	s.mock.ExpectQuery(regexp.QuoteMeta(ExpectedListNewsAlertsQuery)).
		WillReturnRows(sqlmock.NewRows([]string{"news_alert", "file", "name", "link", "created", "from", "to", "File__file", "File__name", "File__path", "File__downloads", "File__url", "File__downloaded"}).
			AddRow(a1.NewsAlert, a1.FileID, a1.Name, a1.Link, a1.Created, a1.From, a1.To, a1.File.File, a1.File.Name, a1.File.Path, a1.File.Downloads, a1.File.URL, a1.File.Downloaded).
			AddRow(a2.NewsAlert, a2.FileID, a2.Name, a2.Link, a2.Created, a2.From, a2.To, a2.File.File, a2.File.Name, a2.File.Path, a2.File.Downloads, a2.File.URL, a2.File.Downloaded))

	server := s.getCampusTestServer()
	response, err := server.ListNewsAlerts(metadata.NewIncomingContext(context.Background(), metadata.MD{}), &pb.ListNewsAlertsRequest{})
	require.NoError(s.T(), err)
	expectedResp := &pb.ListNewsAlertsReply{
		Alerts: []*pb.NewsAlert{
			{ImageUrl: "https://api.tum.app/files/news/sources/src_1.png", Link: a1.Link.String, Created: timestamppb.New(a1.Created), From: timestamppb.New(a1.From), To: timestamppb.New(a1.To)},
			{ImageUrl: "https://api.tum.app/files/news/sources/src_2.png", Link: a2.Link.String, Created: timestamppb.New(a2.Created), From: timestamppb.New(a2.From), To: timestamppb.New(a2.To)},
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

func (s *NewsSuite) getCampusTestServer() *CampusServer {
	return &CampusServer{
		db:              s.DB,
		deviceBuf:       s.deviceBuf,
		newsSourceCache: expirable.NewLRU[string, []model.NewsSource](1, nil, time.Hour*6),
		newsCache:       expirable.NewLRU[string, []model.News](1024, nil, time.Minute*30),
	}
}
