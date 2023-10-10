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

const ExpectedListNewsSourcesQuery = "SELECT `newsSource`.`source`,`newsSource`.`title`,`newsSource`.`url`,`newsSource`.`icon`,`newsSource`.`hook`,`File`.`file` AS `File__file`,`File`.`name` AS `File__name`,`File`.`path` AS `File__path`,`File`.`downloads` AS `File__downloads`,`File`.`url` AS `File__url`,`File`.`downloaded` AS `File__downloaded` FROM `newsSource` LEFT JOIN `files` `File` ON `newsSource`.`icon` = `File`.`file`"

func (s *NewsSuite) Test_ListNewsSourcesMultiple() {
	s.mock.ExpectQuery(regexp.QuoteMeta(ExpectedListNewsSourcesQuery)).
		WillReturnRows(sqlmock.NewRows([]string{"source", "title", "url", "icon", "hook", "File__file", "File__name", "File__path", "File__downloads", "File__url", "File__downloaded"}).
			AddRow(source1().Source, source1().Title, source1().URL, source1().FileID, source1().Hook, source1().File.File, source1().File.Name, source1().File.Path, source1().File.Downloads, source1().File.URL, source1().File.Downloaded).
			AddRow(source2().Source, source2().Title, source2().URL, source2().FileID, source2().Hook, source2().File.File, source2().File.Name, source2().File.Path, source2().File.Downloads, source2().File.URL, source2().File.Downloaded))

	meta := metadata.MD{}
	server := CampusServer{db: s.DB, deviceBuf: s.deviceBuf}
	response, err := server.ListNewsSources(metadata.NewIncomingContext(context.Background(), meta), nil)
	require.NoError(s.T(), err)
	expectedResp := &pb.ListNewsSourcesReply{
		Sources: []*pb.NewsSource{
			{Source: fmt.Sprintf("%d", source1().Source), Title: source1().Title, IconUrl: "https://api.tum.app/files/news/sources/src_1.png"},
			{Source: fmt.Sprintf("%d", source2().Source), Title: source2().Title, IconUrl: "https://api.tum.app/files/news/sources/src_2.png"},
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
	server := CampusServer{db: s.DB, deviceBuf: s.deviceBuf}
	response, err := server.ListNewsSources(metadata.NewIncomingContext(context.Background(), meta), nil)
	require.NoError(s.T(), err)
	expectedResp := &pb.ListNewsSourcesReply{
		Sources: []*pb.NewsSource(nil),
	}
	require.Equal(s.T(), expectedResp, response)
}

const ExpectedListNewsQuery = "SELECT `news`.`news`,`news`.`date`,`news`.`created`,`news`.`title`,`news`.`description`,`news`.`src`,`news`.`link`,`news`.`image`,`news`.`file`,`File`.`file` AS `File__file`,`File`.`name` AS `File__name`,`File`.`path` AS `File__path`,`File`.`downloads` AS `File__downloads`,`File`.`url` AS `File__url`,`File`.`downloaded` AS `File__downloaded` FROM `news` LEFT JOIN `files` `File` ON `news`.`file` = `File`.`file`"

func (s *NewsSuite) Test_ListNewsNone_withFilters() {
	s.mock.ExpectQuery(regexp.QuoteMeta(ExpectedListNewsQuery+" WHERE src = ? AND news > ?")).
		WithArgs(1, 2).
		WillReturnRows(sqlmock.NewRows([]string{"news", "date", "created", "title", "description", "src", "link", "image", "file", "File__file", "File__name", "File__path", "File__downloads", "File__url", "File__downloaded"}))

	meta := metadata.NewIncomingContext(context.Background(), metadata.MD{})
	server := CampusServer{db: s.DB, deviceBuf: s.deviceBuf}
	response, err := server.ListNews(meta, &pb.ListNewsRequest{NewsSource: 1, LastNewsId: 2})
	require.NoError(s.T(), err)
	expectedResp := &pb.ListNewsReply{
		News: []*pb.News{},
	}
	require.Equal(s.T(), expectedResp, response)
}
func (s *NewsSuite) Test_ListNewsNone() {
	s.mock.ExpectQuery(regexp.QuoteMeta(ExpectedListNewsQuery)).
		WillReturnRows(sqlmock.NewRows([]string{"news", "date", "created", "title", "description", "src", "link", "image", "file", "File__file", "File__name", "File__path", "File__downloads", "File__url", "File__downloaded"}))

	meta := metadata.NewIncomingContext(context.Background(), metadata.MD{})
	server := CampusServer{db: s.DB, deviceBuf: s.deviceBuf}
	response, err := server.ListNews(meta, &pb.ListNewsRequest{})
	require.NoError(s.T(), err)
	expectedResp := &pb.ListNewsReply{
		News: []*pb.News{},
	}
	require.Equal(s.T(), expectedResp, response)
}
func (s *NewsSuite) Test_ListNewsMultiple() {
	n1 := news1()
	n2 := news2()
	s.mock.ExpectQuery(regexp.QuoteMeta(" ")).
		WillReturnRows(sqlmock.NewRows([]string{"news", "date", "created", "title", "description", "src", "link", "image", "file", "File__file", "File__name", "File__path", "File__downloads", "File__url", "File__downloaded"}).
			AddRow(n1.News, n1.Date, n1.Created, n1.Title, n1.Description, n1.Src, n1.Link, n1.Image, n1.FileID, n1.File.File, n1.File.Name, n1.File.Path, n1.File.Downloads, n1.File.URL, n1.File.Downloaded).
			AddRow(n2.News, n2.Date, n2.Created, n2.Title, n2.Description, n2.Src, n2.Link, n2.Image, nil, nil, nil, nil, nil, nil, nil))

	meta := metadata.NewIncomingContext(context.Background(), metadata.MD{})
	server := CampusServer{db: s.DB, deviceBuf: s.deviceBuf}
	response, err := server.ListNews(meta, &pb.ListNewsRequest{})
	require.NoError(s.T(), err)
	expectedResp := &pb.ListNewsReply{
		News: []*pb.News{
			{Id: n1.News, Title: n1.Title, Text: n1.Description, Link: n1.Link, ImageUrl: "https://api.tum.app/files/news/sources/src_1.png", Source: fmt.Sprintf("%d", n1.Src), Created: timestamppb.New(n1.Created), Date: timestamppb.New(n1.Date)},
			{Id: n2.News, Title: n2.Title, Text: n2.Description, Link: n2.Link, ImageUrl: "", Source: fmt.Sprintf("%d", n2.Src), Created: timestamppb.New(n2.Created), Date: timestamppb.New(n2.Date)},
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

const ExpectedListNewsAlertsQuery = "SELECT `news_alert`.`news_alert`,`news_alert`.`file`,`news_alert`.`name`,`news_alert`.`link`,`news_alert`.`created`,`news_alert`.`from`,`news_alert`.`to`,`File`.`file` AS `File__file`,`File`.`name` AS `File__name`,`File`.`path` AS `File__path`,`File`.`downloads` AS `File__downloads`,`File`.`url` AS `File__url`,`File`.`downloaded` AS `File__downloaded` FROM `news_alert` LEFT JOIN `files` `File` ON `news_alert`.`file` = `File`.`file` WHERE news_alert.to >= NOW()"

func (s *NewsSuite) Test_ListNewsAlertsError() {
	s.mock.ExpectQuery(regexp.QuoteMeta(ExpectedListNewsAlertsQuery)).WillReturnError(gorm.ErrInvalidDB)

	meta := metadata.MD{}
	server := CampusServer{db: s.DB, deviceBuf: s.deviceBuf}
	response, err := server.ListNewsAlerts(metadata.NewIncomingContext(context.Background(), meta), &pb.ListNewsAlertsRequest{})
	require.Equal(s.T(), status.Error(codes.Internal, "could not ListNewsAlerts"), err)
	require.Nil(s.T(), response)
}
func (s *NewsSuite) Test_ListNewsAlertsNone_noFilter() {
	s.mock.ExpectQuery(regexp.QuoteMeta(ExpectedListNewsAlertsQuery)).WillReturnError(gorm.ErrRecordNotFound)

	server := CampusServer{db: s.DB, deviceBuf: s.deviceBuf}
	response, err := server.ListNewsAlerts(metadata.NewIncomingContext(context.Background(), metadata.MD{}), &pb.ListNewsAlertsRequest{})
	require.Equal(s.T(), status.Error(codes.NotFound, "no news alerts"), err)
	require.Nil(s.T(), response)
}
func (s *NewsSuite) Test_ListNewsAlertsNone_Filter() {
	s.mock.ExpectQuery(regexp.QuoteMeta(ExpectedListNewsAlertsQuery + " AND news_alert.alert > ?")).WithArgs(42).WillReturnError(gorm.ErrRecordNotFound)

	server := CampusServer{db: s.DB, deviceBuf: s.deviceBuf}
	response, err := server.ListNewsAlerts(metadata.NewIncomingContext(context.Background(), metadata.MD{}), &pb.ListNewsAlertsRequest{LastNewsAlertId: 42})
	require.Equal(s.T(), status.Error(codes.NotFound, "no news alerts"), err)
	require.Nil(s.T(), response)
}
func (s *NewsSuite) Test_ListNewsAlertsMultiple() {
	a1 := alert1()
	a2 := alert2()
	s.mock.ExpectQuery(regexp.QuoteMeta(ExpectedListNewsAlertsQuery)).
		WillReturnRows(sqlmock.NewRows([]string{"news_alert", "file", "name", "link", "created", "from", "to", "Files__file", "File__name", "File__path", "Files__downloads", "Files__url", "Files__downloaded"}).
			AddRow(a1.NewsAlert, a1.FileID, a1.Name, a1.Link, a1.Created, a1.From, a1.To, a1.File.File, a1.File.Name, a1.File.Path, a1.File.Downloads, a1.File.URL, a1.File.Downloaded).
			AddRow(a2.NewsAlert, a2.FileID, a2.Name, a2.Link, a2.Created, a2.From, a2.To, a2.File.File, a2.File.Name, a2.File.Path, a2.File.Downloads, a2.File.URL, a2.File.Downloaded))

	server := CampusServer{db: s.DB, deviceBuf: s.deviceBuf}
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
