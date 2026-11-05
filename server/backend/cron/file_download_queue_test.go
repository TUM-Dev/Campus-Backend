package cron

import (
	"bytes"
	"errors"
	"image"
	"image/png"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/TUM-Dev/Campus-Backend/server/internal/webfetch"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func fileDownloadTestDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock) {
	t.Helper()
	conn, mock, err := sqlmock.New()
	require.NoError(t, err)
	t.Cleanup(func() {
		require.NoError(t, mock.ExpectationsWereMet())
		_ = conn.Close()
	})
	db, err := gorm.Open(mysql.New(mysql.Config{Conn: conn, SkipInitializeWithVersion: true}), &gorm.Config{})
	require.NoError(t, err)
	return db, mock
}

const pendingFileQuery = "SELECT * FROM `files` WHERE downloaded = 0 AND url IS NOT NULL ORDER BY downloads ASC,file ASC LIMIT ?"

func expectFileUpdate(mock sqlmock.Sqlmock, column string, fileID int64) {
	mock.ExpectBegin()
	if column == "downloads" {
		mock.ExpectExec(regexp.QuoteMeta("UPDATE `files` SET `downloads`=downloads + 1 WHERE `file` = ?")).
			WithArgs(fileID).WillReturnResult(sqlmock.NewResult(0, 1))
	} else {
		mock.ExpectExec(regexp.QuoteMeta("UPDATE `files` SET `downloaded`=? WHERE `file` = ?")).
			WithArgs(true, fileID).WillReturnResult(sqlmock.NewResult(0, 1))
	}
	mock.ExpectCommit()
}

func TestFileDownloadQueueCommitsProgressAndContinuesAfterFailure(t *testing.T) {
	db, mock := fileDownloadTestDB(t)
	oldStorage, oldClient := StorageDir, webfetch.Client
	StorageDir = t.TempDir()
	t.Cleanup(func() { StorageDir, webfetch.Client = oldStorage, oldClient })

	var poster bytes.Buffer
	require.NoError(t, png.Encode(&poster, image.NewRGBA(image.Rect(0, 0, 2, 2))))
	var requested []string
	webfetch.Client = &http.Client{Transport: movieTestTransport(func(r *http.Request) (*http.Response, error) {
		requested = append(requested, r.URL.Path)
		if r.URL.Path == "/missing" {
			// The previous file must already be committed when this request starts.
			require.NoError(t, mock.ExpectationsWereMet())
			expectFileUpdate(mock, "downloads", int64(3))
			expectFileUpdate(mock, "downloaded", int64(3))
			return &http.Response{StatusCode: 404, Status: "404 Not Found", Body: io.NopCloser(strings.NewReader("missing"))}, nil
		}
		return &http.Response{StatusCode: 200, Body: io.NopCloser(bytes.NewReader(poster.Bytes()))}, nil
	})}

	mock.ExpectQuery(regexp.QuoteMeta(pendingFileQuery)).WithArgs(fileDownloadBatchSize).
		WillReturnRows(sqlmock.NewRows([]string{"file", "name", "path", "url"}).
			AddRow(1, "first.jpg", "movie/", "https://example.com/first").
			AddRow(2, "missing.jpg", "movie/", "https://example.com/missing").
			AddRow(3, "last.jpg", "movie/", "https://example.com/last"))
	expectFileUpdate(mock, "downloads", int64(1))
	expectFileUpdate(mock, "downloaded", int64(1))
	expectFileUpdate(mock, "downloads", int64(2))

	err := (&CronService{db: db}).fileDownloadCron()
	require.ErrorContains(t, err, "file 2: upstream returned 404 Not Found")
	require.Equal(t, []string{"/first", "/missing", "/last"}, requested)
	for _, name := range []string{"first.jpg", "last.jpg"} {
		_, err := os.Stat(filepath.Join(StorageDir, "movie", name))
		require.NoError(t, err)
	}
}

func TestFileDownloadQueueReturnsQueryFailure(t *testing.T) {
	db, mock := fileDownloadTestDB(t)
	want := errors.New("database unavailable")
	mock.ExpectQuery(regexp.QuoteMeta(pendingFileQuery)).WithArgs(fileDownloadBatchSize).WillReturnError(want)
	require.ErrorIs(t, (&CronService{db: db}).fileDownloadCron(), want)
}

func TestFileDownloadQueueWithNoPendingFiles(t *testing.T) {
	db, mock := fileDownloadTestDB(t)
	mock.ExpectQuery(regexp.QuoteMeta(pendingFileQuery)).WithArgs(fileDownloadBatchSize).
		WillReturnRows(sqlmock.NewRows([]string{"file"}))
	require.NoError(t, (&CronService{db: db}).fileDownloadCron())
}
