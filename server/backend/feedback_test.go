package backend

import (
	"bytes"
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"image"
	"image/png"
	"io"
	"os"
	"path"
	"sync"
	"testing"
	"time"

	pb "github.com/TUM-Dev/Campus-Backend/server/api/tumdev"
	"github.com/TUM-Dev/Campus-Backend/server/backend/cron"
	"github.com/TUM-Dev/Campus-Backend/server/model"
	"github.com/TUM-Dev/Campus-Backend/server/utils"
	"github.com/guregu/null"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
)

type mockedFeedbackStream struct {
	grpc.ServerStream
	recived []*pb.CreateFeedbackRequest
	reply   *pb.CreateFeedbackReply
	T       *testing.T
}

func (f mockedFeedbackStream) SendAndClose(reply *pb.CreateFeedbackReply) error {
	require.Equal(f.T, f.reply, reply)
	return nil
}

// because of the way the mocked stream works, we need to keep track of the index
// we however can't track this inside of mockedFeedbackStream, as we cannot mutate the struct
// => tracking this as  a global variable is the only way
var index = uint(0)

func (f mockedFeedbackStream) Recv() (*pb.CreateFeedbackRequest, error) {
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

func Test_CreateFeedback_TwoFiles(t *testing.T) {
	// this is not parallelism because cron.StorageDir is SHARED STATE
	// => needs to be NOT A RACE CONDITION to be faster
	// t.Parallel()
	ctx := context.Background()
	db := utils.SetupTestContainer(ctx, t)
	// -- setup above
	dir, err := os.MkdirTemp("", "two_files")
	require.NoError(t, err)
	defer require.NoError(t, os.RemoveAll(dir))
	cron.StorageDir = dir

	server := CampusServer{db: db, feedbackEmailLastReuestAt: &sync.Map{}}
	returnedTime := time.Now()

	// execute call
	dummyImage := createDummyImage(t, 10, 10)
	dummyText := []byte("Dummy Text")
	stream := mockedFeedbackStream{
		T: t,
		recived: []*pb.CreateFeedbackRequest{
			{Recipient: pb.CreateFeedbackRequest_TUM_DEV, FromEmail: "testing@example.com", Message: "Hello with image", Attachment: dummyText},
			{Attachment: dummyImage},
		},
		reply: &pb.CreateFeedbackReply{},
	}
	require.NoError(t, server.CreateFeedback(stream))

	// check that the correct operations happened to the file system
	fsFiles := extractUploadedFiles(t, cron.StorageDir, 2)
	expectFileMatches(t, (*fsFiles)[0], "0.txt", returnedTime, dummyText)
	expectFileMatches(t, (*fsFiles)[1], "1.png", returnedTime, dummyImage)

	// should have inserted feedback
	var feeedbacks []model.Feedback
	require.NoError(t, db.WithContext(ctx).Find(&feeedbacks).Error)
	require.Len(t, feeedbacks, 1)
	actual := feeedbacks[0]
	require.Equal(t, "app@tum.de", actual.Recipient)
	require.Equal(t, null.StringFrom("testing@example.com"), actual.ReplyToEmail)
	require.Equal(t, "Hello with image", actual.Feedback)

	// should have created files
	var dbFiles []model.File
	require.NoError(t, db.WithContext(ctx).Find(&dbFiles).Error)
	require.Len(t, dbFiles, 2)
	actualFile := dbFiles[0]
	require.Equal(t, "0.txt", actualFile.Name)
	require.Equal(t, int32(1), actualFile.Downloads)
	require.Equal(t, null.BoolFrom(true), actualFile.Downloaded)
	actualFile = dbFiles[1]
	require.Equal(t, "1.png", actualFile.Name)
	require.Equal(t, int32(1), actualFile.Downloads)
	require.Equal(t, null.BoolFrom(true), actualFile.Downloaded)

	// test if re-submitting feedback is blocked
	stream2 := mockedFeedbackStream{
		T: t,
		recived: []*pb.CreateFeedbackRequest{
			{Recipient: pb.CreateFeedbackRequest_TUM_DEV, FromEmail: "testing@example.com", Message: "Hello with image", Attachment: dummyText},
			{Attachment: dummyImage},
		},
		reply: &pb.CreateFeedbackReply{},
	}
	require.Error(t, server.CreateFeedback(stream2), status.Error(codes.ResourceExhausted, "You have already send a feedback recently. Please wait 300 seconds"))

	// the db did not change
	var feeedbacks2 []model.Feedback
	require.NoError(t, db.WithContext(ctx).Find(&feeedbacks2).Error)
	require.Equal(t, feeedbacks, feeedbacks2)
	var dbFiles2 []model.File
	err = db.WithContext(ctx).Find(&dbFiles2).Error
	require.NoError(t, err)
	require.Equal(t, dbFiles, dbFiles2)
	// all files that were added are cleaned up correctly
	parentDir, err := os.ReadDir(path.Join(cron.StorageDir, "feedback"))
	require.NoError(t, err)
	require.Len(t, parentDir, 1)
}

func expectFileMatches(t *testing.T, file os.DirEntry, name string, returnedTime time.Time, content []byte) {
	require.Equal(t, name, file.Name())
	require.True(t, file.Type().IsRegular())
	info, err := file.Info()
	require.NoError(t, err)
	require.LessOrEqual(t, returnedTime.Unix(), info.ModTime().Unix())
	require.Len(t, content, int(info.Size()))
}

func Test_CreateFeedback_NoImage(t *testing.T) {
	// this is not paralelisable because cron.StorageDir is SHARED STATE
	// => needs to be NOT A RACE CONDITION to be faster
	// t.Parallel()
	ctx := context.Background()
	db := utils.SetupTestContainer(ctx, t)
	// -- setup above
	dir, err := os.MkdirTemp("", "no_files")
	require.NoError(t, err)
	defer require.NoError(t, os.RemoveAll(dir))

	server := CampusServer{db: db, feedbackEmailLastReuestAt: &sync.Map{}}
	stream := mockedFeedbackStream{
		T: t,
		recived: []*pb.CreateFeedbackRequest{
			{Recipient: pb.CreateFeedbackRequest_TUM_DEV, FromEmail: "testing@example.com", Message: "Hello without image"},
			{}, // empty images should be ignored
		},
		reply: &pb.CreateFeedbackReply{},
	}
	require.NoError(t, server.CreateFeedback(stream))

	// no image should be uploaded to the file system
	extractUploadedFiles(t, cron.StorageDir, 0)

	// should have inserted feedback
	var feeedbacks []model.Feedback
	err = db.WithContext(ctx).Find(&feeedbacks).Error
	require.NoError(t, err)
	require.Len(t, feeedbacks, 1)
	actual := feeedbacks[0]
	require.Equal(t, "app@tum.de", actual.Recipient)
	require.Equal(t, null.StringFrom("testing@example.com"), actual.ReplyToEmail)
	require.Equal(t, "Hello without image", actual.Feedback)

	// should have created no files
	var dbFiles []model.File
	err = db.WithContext(ctx).Find(&dbFiles).Error
	require.NoError(t, err)
	require.Equal(t, dbFiles, []model.File{})
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
