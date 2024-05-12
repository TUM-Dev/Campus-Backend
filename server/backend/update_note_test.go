package backend

import (
	"context"
	"database/sql"
	"regexp"
	"testing"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

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
		WillReturnRows(sqlmock.NewRows([]string{"VERSION()"}).AddRow("8.4.0"))
	s.DB, err = gorm.Open(dialector, &gorm.Config{})
	require.NoError(s.T(), err)

	s.deviceBuf = newDeviceBuffer()
}

const ExpectedGetUpdateNoteQuery = "SELECT * FROM `update_note` WHERE `update_note`.`version_code` = ? ORDER BY `update_note`.`version_code` LIMIT ?"

func (s *UpdateNoteSuite) Test_GetUpdateNoteOne() {
	s.mock.ExpectQuery(regexp.QuoteMeta(ExpectedGetUpdateNoteQuery)).
		WithArgs(1, 1).
		WillReturnRows(sqlmock.NewRows([]string{"version_code", "version_name", "message"}).
			AddRow(1, "1.0.0", "Test Message"))

	meta := metadata.MD{}
	server := CampusServer{db: s.DB, deviceBuf: s.deviceBuf}
	response, err := server.GetUpdateNote(metadata.NewIncomingContext(context.Background(), meta), &pb.GetUpdateNoteRequest{Version: 1})
	require.NoError(s.T(), err)
	expectedResp := &pb.GetUpdateNoteReply{
		Message:     "Test Message",
		VersionName: "1.0.0",
	}
	require.Equal(s.T(), expectedResp, response)
}

func (s *UpdateNoteSuite) Test_GetUpdateNoteNone() {
	s.mock.ExpectQuery(regexp.QuoteMeta(ExpectedGetUpdateNoteQuery)).
		WithArgs(1, 1).
		WillReturnRows(sqlmock.NewRows([]string{"version_code", "version_name", "message"}))

	meta := metadata.MD{}
	server := CampusServer{db: s.DB, deviceBuf: s.deviceBuf}
	response, err := server.GetUpdateNote(metadata.NewIncomingContext(context.Background(), meta), &pb.GetUpdateNoteRequest{Version: 1})
	require.Equal(s.T(), status.Error(codes.NotFound, "No update note found"), err)
	require.Nil(s.T(), response)
}

func (s *UpdateNoteSuite) Test_GetUpdateNoteError() {
	s.mock.ExpectQuery(regexp.QuoteMeta(ExpectedGetUpdateNoteQuery)).
		WithArgs(1, 1).
		WillReturnError(gorm.ErrInvalidDB)

	meta := metadata.MD{}
	server := CampusServer{db: s.DB, deviceBuf: s.deviceBuf}
	response, err := server.GetUpdateNote(metadata.NewIncomingContext(context.Background(), meta), &pb.GetUpdateNoteRequest{Version: 1})
	require.Equal(s.T(), status.Error(codes.Internal, "Internal server error"), err)
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
