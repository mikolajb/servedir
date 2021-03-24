package archivist

import (
	"archive/zip"
	"bytes"
	"errors"
	"io/ioutil"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const readerErrorText = "something bad happened with a reader"

func TestArchivist(t *testing.T) {
	var buffer bytes.Buffer
	filesToRead := map[string]string{
		"foo.txt": "foo foo",
		"bar.txt": "bar bar",
		"baz.txt": "baz baz",
	}

	arch := New(&buffer)

	for fileName, fileContent := range filesToRead {
		require.NoError(t, arch.Add(fileName, strings.NewReader(fileContent)))
	}
	require.NoError(t, arch.Close())

	reader, err := zip.NewReader(bytes.NewReader(buffer.Bytes()), int64(buffer.Len()))
	require.NoError(t, err)

	for _, f := range reader.File {
		content, ok := filesToRead[f.Name]
		if !ok {
			t.Errorf("file %s is not expected", f.Name)
			continue
		}

		openedFile, err := f.Open()
		require.NoError(t, err)
		readContent, err := ioutil.ReadAll(openedFile)
		require.NoError(t, err)

		assert.EqualValues(t, content, readContent)

		delete(filesToRead, f.Name)
	}

	if len(filesToRead) > 0 {
		t.Errorf("these files were not in the archive: %v", filesToRead)
	}
}

func TestArchivistHandlingEmptyPath(t *testing.T) {
	arch := New(nil)
	err := arch.Add("", nil)
	assert.NoError(t, err)
}

func TestArchivistHandlingReaderError(t *testing.T) {
	arch := New(nil)
	err := arch.Add("/something", &testReader{})
	assert.EqualError(t, err, readerErrorText)
}

func TestArchivistHandlingWtierError(t *testing.T) {
	arch := New(nil)
	var overflowedUint16 = 1 << 16
	err := arch.Add(strings.Repeat("a", overflowedUint16), nil)
	assert.EqualError(t, err, "zip: FileHeader.Name too long")
}

type testReader struct{}

func (*testReader) Read([]byte) (int, error) {
	return 0, errors.New(readerErrorText)
}
