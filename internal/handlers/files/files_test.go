package files

import (
	"context"
	"io/ioutil"
	"net/http/httptest"
	"net/url"
	"os"
	"path"
	"testing"

	"github.com/mikolajb/servedir/internal/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHandler(t *testing.T) {
	tempFile, err := ioutil.TempFile("", "servedir")
	require.NoError(t, err)
	defer tempFile.Close()

	dir, file := path.Split(tempFile.Name())

	handler := New(dir)
	request := httptest.NewRequest("GET", "/"+file, nil)
	_, ctx := logger.NewTestLogger(t, context.Background())
	request = request.Clone(ctx)
	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, request)
}

func TestGetTargetHandlesStatError(t *testing.T) {
	_, ctx := logger.NewTestLogger(t, context.Background())
	_, err := getTargetHandler(ctx, url.Values{}, "", "")
	assert.EqualError(t, err, "stat : no such file or directory")
}

func TestGetTargetRecognizesSingleFile(t *testing.T) {
	_, ctx := logger.NewTestLogger(t, context.Background())
	tempFile, err := ioutil.TempFile("", "servedir")
	require.NoError(t, err)
	defer tempFile.Close()

	target, err := getTargetHandler(ctx, url.Values{}, "", tempFile.Name())
	require.NoError(t, err)
	_, ok := target.(*singleFile)
	assert.True(t, ok)
}

func TestGetTargetRecognizesDirectory(t *testing.T) {
	_, ctx := logger.NewTestLogger(t, context.Background())
	tempDir, err := ioutil.TempDir("", "servedir")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	target, err := getTargetHandler(ctx, url.Values{}, "", tempDir)
	require.NoError(t, err)
	_, ok := target.(*dirHandler)
	assert.True(t, ok)
}

func TestGetTargetRecognizesDirectoryZip(t *testing.T) {
	_, ctx := logger.NewTestLogger(t, context.Background())
	tempDir, err := ioutil.TempDir("", "servedir")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	values := url.Values{}
	values.Set("archive", "true")
	target, err := getTargetHandler(ctx, values, "", tempDir+".zip")
	require.NoError(t, err)
	_, ok := target.(*zipHandler)
	assert.True(t, ok)
}

func TestGetTargetRecognizesDirectoryZipSuffixError(t *testing.T) {
	_, ctx := logger.NewTestLogger(t, context.Background())
	tempDir, err := ioutil.TempDir("", "servedir")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	values := url.Values{}
	values.Set("archive", "true")
	target, err := getTargetHandler(ctx, values, "", tempDir)
	require.EqualError(t, err, ErrNoZipSuffix.Error())
	assert.Nil(t, target)
}
