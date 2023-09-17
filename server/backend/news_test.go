package backend

import (
	"context"
	"database/sql"
	"fmt"
	"regexp"
	"testing"

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

func newsFile(id int32) *model.Files {
	return &model.Files{
		File:       id,
		Name:       fmt.Sprintf("src_%d.png", id),
		Path:       "news/sources",
		Downloads:  1,
		URL:        sql.NullString{Valid: false},
		Downloaded: sql.NullBool{Bool: true, Valid: true},
	}
}

func source1() *model.NewsSource {
	return &model.NewsSource{
		Source:  1,
		Title:   "Amazing News 1",
		URL:     null.String{NullString: sql.NullString{String: "https://example.com/amazing1", Valid: true}},
		FilesID: newsFile(1).File,
		Files:   *newsFile(1),
		Hook:    null.String{NullString: sql.NullString{String: "", Valid: true}},
	}
}

func source2() *model.NewsSource {
	return &model.NewsSource{
		Source:  2,
		Title:   "Amazing News 2",
		URL:     null.String{NullString: sql.NullString{String: "https://example.com/amazing2", Valid: true}},
		FilesID: newsFile(2).File,
		Files:   *newsFile(2),
		Hook:    null.String{NullString: sql.NullString{String: "hook", Valid: true}},
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
	expectedResp := &pb.NewsSourceReply{
		Sources: []*pb.NewsSource{
			{Source: fmt.Sprintf("%d", source1().Source), Title: source1().Title, Icon: source1().Files.URL.String},
			{Source: fmt.Sprintf("%d", source2().Source), Title: source2().Title, Icon: source2().Files.URL.String},
		},
	}
	require.Equal(s.T(), expectedResp, response)
}

func (s *NewsSuite) Test_GetNewsSourcesNone() {
	s.mock.ExpectQuery(regexp.QuoteMeta(ExpectedGetSourceQuery)).
		WillReturnRows(sqlmock.NewRows([]string{"source", "title", "url", "icon", "hook", "Files__file", "Files__name", "Files__path", "Files__downloads", "Files__url", "Files__downloaded"}))

	meta := metadata.MD{}
	server := CampusServer{db: s.DB, deviceBuf: s.deviceBuf}
	response, err := server.GetNewsSources(metadata.NewIncomingContext(context.Background(), meta), nil)
	require.NoError(s.T(), err)
	expectedResp := &pb.NewsSourceReply{
		Sources: []*pb.NewsSource(nil),
	}
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
