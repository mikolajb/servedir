package upload

import (
	"bytes"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	writerErrorText = "writer error"
	outputPage      = `<!DOCTYPE HTML><html><head><title>File uploader</title><style>body{font-size:2em;width: 60%;margin-left:auto;margin-right:auto;margin-top:80px;}</style></head><body><a href="/files/">Browse</a> <a href="/upload/">Upload</a>
<a href="/upload/?files=4">One more</a>
<a href="/upload/?files=2">One less</a>
<form enctype="multipart/form-data" action="/upload/" method="post">

      <input type="file" name="files" style="width:100%">

      <input type="file" name="files" style="width:100%">

      <input type="file" name="files" style="width:100%">

      <input style="width:9%" type="submit" value="Upload">
</form></body></html>
`
)

func TestUploadInitialization(t *testing.T) {
	upld := New(0)
	assert.EqualValues(t, 0, upld.OneLess)
	assert.EqualValues(t, 2, upld.OneMore)

	upld = New(5)
	assert.EqualValues(t, 4, upld.OneLess)
	assert.EqualValues(t, 6, upld.OneMore)

	upld = New(21)
	assert.EqualValues(t, 19, upld.OneLess)
	assert.EqualValues(t, 21, upld.OneMore)
}

func TestWriterReturnsErrors(t *testing.T) {
	upld := New(3)

	err := upld.WritePage(&faultyWriter{})
	assert.EqualError(t, err, writerErrorText)
}

func TestWritePage(t *testing.T) {
	upld := New(3)

	var buf bytes.Buffer
	upld.WritePage(&buf)
	assert.Equal(t, outputPage, buf.String())
}

type faultyWriter struct{}

func (*faultyWriter) Write([]byte) (int, error) {
	return 0, errors.New(writerErrorText)
}
