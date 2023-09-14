package backend

import (
	"context"
	"database/sql"
	"fmt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	pb "github.com/TUM-Dev/Campus-Backend/server/api/tumdev"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc/metadata"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type UpdateNoteSuite struct {
	suite.Suite
	DB        *gorm.DB
	mock      sqlmock.Sqlmock
	deviceBuf *deviceBuffer
}

func (s *UpdateNoteSuite) SetupSuite() {
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

const ExpectedGetUpdateNoteQuery = "SELECT * FROM `update_note` WHERE `update_note`.`version_code` = ? ORDER BY `update_note`.`version_code` LIMIT 1"

func (s *UpdateNoteSuite) Test_GetUpdateNoteOne() {
	s.mock.ExpectQuery(regexp.QuoteMeta(ExpectedGetUpdateNoteQuery)).
		WillReturnRows(sqlmock.NewRows([]string{"source", "title", "url"}).
			AddRow(source1().Source, source1().Title, source1().URL, source1().FilesID, source1().Hook, source1().Files.File, source1().Files.Name, source1().Files.Path, source1().Files.Downloads, source1().Files.URL, source1().Files.Downloaded).
			AddRow(source2().Source, source2().Title, source2().URL, source2().FilesID, source2().Hook, source2().Files.File, source2().Files.Name, source2().Files.Path, source2().Files.Downloads, source2().Files.URL, source2().Files.Downloaded))

	meta := metadata.MD{}
	server := CampusServer{db: s.DB, deviceBuf: s.deviceBuf}
	response, err := server.GetUpdateNote(metadata.NewIncomingContext(context.Background(), meta), &pb.GetUpdateNoteRequest{Version: 1})
	require.NoError(s.T(), err)
	expectedResp := &pb.NewsSourceReply{
		Sources: []*pb.NewsSource{
			{Source: fmt.Sprintf("%d", source1().Source), Title: source1().Title, Icon: source1().Files.URL.String},
			{Source: fmt.Sprintf("%d", source2().Source), Title: source2().Title, Icon: source2().Files.URL.String},
		},
	}
	require.Equal(s.T(), expectedResp, response)
}

func (s *UpdateNoteSuite) Test_GetUpdateNoteNone() {
	s.mock.ExpectQuery(regexp.QuoteMeta(ExpectedGetSourceQuery)).
		WillReturnRows(sqlmock.NewRows([]string{"source", "title", "url"}))

	meta := metadata.MD{}
	server := CampusServer{db: s.DB, deviceBuf: s.deviceBuf}
	response, err := server.GetUpdateNote(metadata.NewIncomingContext(context.Background(), meta), &pb.GetUpdateNoteRequest{Version: 1})
	require.NoError(s.T(), err)
	expectedResp := &pb.NewsSourceReply{
		Sources: []*pb.NewsSource(nil),
	}
	require.Equal(s.T(), expectedResp, response)
}

func (s *UpdateNoteSuite) Test_GetUpdateNoteError() {
	s.mock.ExpectQuery(regexp.QuoteMeta(ExpectedGetTopNewsQuery)).WillReturnError(gorm.ErrInvalidDB)

	meta := metadata.MD{}
	server := CampusServer{db: s.DB, deviceBuf: s.deviceBuf}
	response, err := server.GetUpdateNote(metadata.NewIncomingContext(context.Background(), meta), &pb.GetUpdateNoteRequest{Version: 1})
	require.Equal(s.T(), status.Error(codes.Internal, "could not GetTopNews"), err)
	require.Nil(s.T(), response)
}

func (s *UpdateNoteSuite) AfterTest(_, _ string) {
	require.NoError(s.T(), s.mock.ExpectationsWereMet())
}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestUpdateNoteSuite(t *testing.T) {
	suite.Run(t, new(UpdateNoteSuite))
}
