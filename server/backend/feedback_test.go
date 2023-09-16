package backend

import (
	"bytes"
	"context"
	"database/sql"
	"image"
	"image/png"
	"io"
	"os"
	"path"
	"regexp"
	"testing"
	"time"

	"github.com/TUM-Dev/Campus-Backend/server/backend/cron"

	"github.com/DATA-DOG/go-sqlmock"
	pb "github.com/TUM-Dev/Campus-Backend/server/api/tumdev"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type FeedbackSuite struct {
	suite.Suite
	DB   *gorm.DB
	mock sqlmock.Sqlmock
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
}

type mockedFeedbackStream struct {
	grpc.ServerStream
	recived []*pb.SendFeedbackRequest
	reply   *pb.SendFeedbackReply
	T       *testing.T
}

func (f mockedFeedbackStream) SendAndClose(reply *pb.SendFeedbackReply) error {
	require.Equal(f.T, f.reply, reply)
	return nil
}

// because of the way the mocked stream works, we need to keep track of the index
// we however can't track this inside of mockedFeedbackStream, as we cannot mutate the struct
// => tracking this as  a global variable is the only way
var index = uint(0)

func (f mockedFeedbackStream) Recv() (*pb.SendFeedbackRequest, error) {
	if int(index) >= len(f.recived) {
		return nil, io.EOF
	}
	index++
	return f.recived[index-1], nil
}

func (f mockedFeedbackStream) Context() context.Context {
	// reset index, as this function is called before Recv()
	// This is a hacky solution, but it works
	index = 0
	return context.Background()
}

// createDummyImage creates a dummy image with the specified dimensions
func createDummyImage(t *testing.T, width, height int) []byte {
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	// encode img to buffer
	buf := new(bytes.Buffer)
	require.NoError(t, png.Encode(buf, img))
	return buf.Bytes()
}

func (s *FeedbackSuite) Test_SendFeedback_OneImage() {
	cron.StorageDir = "test_one_image/"
	defer func(path string) {
		err := os.RemoveAll(path)
		if err != nil {
			s.T().Fatal(err)
		}
	}(cron.StorageDir)

	server := CampusServer{db: s.DB}
	s.mock.ExpectBegin()
	returnedTime := time.Now()
	s.mock.ExpectQuery(regexp.QuoteMeta("INSERT INTO `files` (`name`,`path`,`downloads`,`downloaded`) VALUES (?,?,?,?) RETURNING `url`,`file`")).
		WithArgs("0.txt", sqlmock.AnyArg(), 1, true).
		WillReturnRows(sqlmock.NewRows([]string{"url", "file"}).AddRow(nil, 1))
	s.mock.ExpectQuery(regexp.QuoteMeta("INSERT INTO `files` (`name`,`path`,`downloads`,`downloaded`) VALUES (?,?,?,?) RETURNING `url`,`file`")).
		WithArgs("1.png", sqlmock.AnyArg(), 1, true).
		WillReturnRows(sqlmock.NewRows([]string{"url", "file"}).AddRow(nil, 1))
	s.mock.ExpectQuery(regexp.QuoteMeta("INSERT INTO `feedback` (`image_count`,`email_id`,`receiver`,`reply_to`,`feedback`,`latitude`,`longitude`,`os_version`,`app_version`,`processed`) VALUES (?,?,?,?,?,?,?,?,?,?) RETURNING `timestamp`,`id`")).
		WithArgs(2, sqlmock.AnyArg(), "app@tum.de", "testing@example.com", "Hello World", nil, nil, nil, nil, false).
		WillReturnRows(sqlmock.NewRows([]string{"timestamp", "id"}).AddRow(returnedTime, 1))
	s.mock.ExpectCommit()

	dummyImage := createDummyImage(s.T(), 10, 10)
	dummyText := []byte("Dummy Text")
	stream := mockedFeedbackStream{
		T: s.T(),
		recived: []*pb.SendFeedbackRequest{
			{Recipient: pb.SendFeedbackRequest_TUM_DEV, FromEmail: "testing@example.com", Message: "Hello World", Metadata: &pb.SendFeedbackRequest_Metadata{}, Attachment: dummyText},
			{Attachment: dummyImage},
		},
		reply: &pb.SendFeedbackReply{},
	}
	require.NoError(s.T(), server.SendFeedback(stream))

	// check that the correct operations happened to the file system
	files := extractUploadedFiles(s.T(), cron.StorageDir, 2)
	expectFileMatches(s.T(), (*files)[0], "0.txt", returnedTime, dummyText)
	expectFileMatches(s.T(), (*files)[1], "1.png", returnedTime, dummyImage)
}

func expectFileMatches(t *testing.T, file os.DirEntry, name string, returnedTime time.Time, content []byte) {
	require.Equal(t, name, file.Name())
	require.True(t, file.Type().IsRegular())
	info, err := file.Info()
	require.NoError(t, err)
	require.LessOrEqual(t, returnedTime.Unix(), info.ModTime().Unix())
	require.Len(t, content, int(info.Size()))
}

func (s *FeedbackSuite) Test_SendFeedback_NoImage() {
	cron.StorageDir = "test_no_image/"

	server := CampusServer{db: s.DB}
	s.mock.ExpectBegin()
	s.mock.ExpectQuery(regexp.QuoteMeta("INSERT INTO `feedback` (`image_count`,`email_id`,`receiver`,`reply_to`,`feedback`,`latitude`,`longitude`,`os_version`,`app_version`,`processed`) VALUES (?,?,?,?,?,?,?,?,?,?) RETURNING `timestamp`,`id`")).
		WithArgs(0, sqlmock.AnyArg(), "app@tum.de", "testing@example.com", "Hello World", nil, nil, nil, nil, false).
		WillReturnRows(sqlmock.NewRows([]string{"timestamp", "id"}).AddRow(time.Now(), 1))
	s.mock.ExpectCommit()

	stream := mockedFeedbackStream{
		T: s.T(),
		recived: []*pb.SendFeedbackRequest{
			{Recipient: pb.SendFeedbackRequest_TUM_DEV, FromEmail: "testing@example.com", Message: "Hello World", Metadata: &pb.SendFeedbackRequest_Metadata{}},
			{}, // empty images should be ignored
		},
		reply: &pb.SendFeedbackReply{},
	}
	require.NoError(s.T(), server.SendFeedback(stream))

	// no image should be uploaded to the file system
	extractUploadedFiles(s.T(), cron.StorageDir, 0)
}

func extractUploadedFiles(t *testing.T, storageRoot string, expected int) *[]os.DirEntry {
	parentDir, err := os.ReadDir(path.Join(storageRoot, "feedback"))
	if expected == 0 {
		require.Error(t, err, os.ErrNotExist.Error())
		require.Empty(t, parentDir)
		return nil
	}

	require.NoError(t, err)
	require.Len(t, parentDir, 1)
	dir, err := os.ReadDir(path.Join(storageRoot, "feedback", parentDir[0].Name()))
	require.NoError(t, err)
	require.Len(t, dir, expected)
	return &dir
}

func (s *FeedbackSuite) AfterTest(_, _ string) {
	require.NoError(s.T(), s.mock.ExpectationsWereMet())
}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestExampleTestSuite(t *testing.T) {
	suite.Run(t, new(FeedbackSuite))
}
