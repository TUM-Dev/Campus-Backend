package cron

import (
	"image"
	"os"
	"testing"

	"github.com/disintegration/imaging"
	"github.com/stretchr/testify/require"
)

func TestMaybeResizeImage(t *testing.T) {
	t.Run("Resize Image", func(t *testing.T) {
		dstPath := "test_image.jpg"
		require.NoError(t, createDummyImage(dstPath, 2000, 1000))
		defer os.Remove(dstPath)
		require.NoError(t, resizeImage(dstPath))
		img, err := imaging.Open(dstPath)
		require.NoError(t, err)
		require.Equal(t, 1280, img.Bounds().Dx())
		require.Equal(t, 640, img.Bounds().Dy())
	})
	t.Run("Do not Resize smaller Image", func(t *testing.T) {
		dstPath := "test_image.jpg"
		require.NoError(t, createDummyImage(dstPath, 1000, 2000))
		defer os.Remove(dstPath)
		require.NoError(t, resizeImage(dstPath))
		img, err := imaging.Open(dstPath)
		require.NoError(t, err)
		require.Equal(t, 1280, img.Bounds().Dx())
		require.Equal(t, 2560, img.Bounds().Dy())
	})
}

func TestEnsureFileDoesNotExist(t *testing.T) {
	tmpFilePath := "test_dir/test_file.txt"
	defer func() { _ = os.RemoveAll("test_dir") }()

	t.Run("FileDoesNotExist", func(t *testing.T) {
		require.NoError(t, ensureFileDoesNotExist(tmpFilePath))

		_, dirErr := os.Stat("test_dir")
		require.NoError(t, dirErr)

		_, fileErr := os.Stat(tmpFilePath)
		require.True(t, os.IsNotExist(fileErr))
	})

	t.Run("FileExists", func(t *testing.T) {
		_, createErr := os.Create(tmpFilePath)
		require.NoError(t, createErr)

		require.NoError(t, ensureFileDoesNotExist(tmpFilePath))

		_, dirErr := os.Stat("test_dir")
		require.NoError(t, dirErr)

		_, fileErr := os.Stat(tmpFilePath)
		require.True(t, os.IsNotExist(fileErr))
	})
}

// createDummyImage creates a dummy image file with the specified dimensions
func createDummyImage(filePath string, width, height int) error {
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	return imaging.Save(img, filePath, imaging.JPEGQuality(75))
}

// createDummyFile creates a dummy non-image file
func createDummyFile(filePath string, content []byte) error {
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()
	if _, err := file.Write(content); err != nil {
		return err
	}
	return nil
}
