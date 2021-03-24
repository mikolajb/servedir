package listing

import (
	"bytes"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	writerErrorText = "writer error"
	outputPage      = `<!DOCTYPE HTML><html><head><title>File downloader</title><style>body{font-size:2em;width:60%;margin-left:auto;margin-right:auto;margin-top:80px;}</style></head><body><a href="/files/">Browse</a> <a href="/upload/">Upload</a>
    <table style="width:100%;">
      <tbody>
       <tr><td><a href="current-path">.</a></td><td></td><td></td></tr>
       <tr><td><a href="up-path">..</a></td><td></td><td></td></tr>

       <tr><td>
            <a href="path">name.txt</a>

          </td>

            <td style="width:10%;text-align:right;"><i>123</i></td>

          <td style="width:35%;font-family:monospace;text-align:right;">ABC</td></tr>
</tbody></table></body></html>`
)

var (
	currentPath = "current-path"
	upPath      = "up-path"
)

func TestAddFile(t *testing.T) {
	lsting := New(currentPath, upPath)

	assert.Equal(t, currentPath, lsting.Current)
	assert.Equal(t, upPath, lsting.Up)
	assert.Len(t, lsting.Files, 0)

	lsting.AddFile(&FileInfo{})
	assert.Len(t, lsting.Files, 1)
}

func TestWritePage(t *testing.T) {
	lsting := New(currentPath, upPath)
	lsting.AddFile(&FileInfo{
		Name: "name.txt",
		Path: "path",
		Size: 123,
		Mod:  "ABC",
	})

	var buf bytes.Buffer
	lsting.WritePage(&buf)
	assert.Equal(t, outputPage, buf.String())
}

func TestWriterReturnsError(t *testing.T) {
	lsting := New("", "")
	err := lsting.WritePage(&faultyWriter{})
	assert.EqualError(t, err, writerErrorText)
}

func TestFormat(t *testing.T) {
	cases := []struct {
		size     FileSize
		expected string
	}{
		{
			size:     0,
			expected: "0",
		},
		{
			size:     1001,
			expected: "1K",
		},
		{
			size:     1234567,
			expected: "1M",
		},
		{
			size:     1234567000,
			expected: "1G",
		},
		{
			size:     1234567000541,
			expected: "1T",
		},
		{
			size:     1234567000541541,
			expected: "1234T",
		},
	}

	for _, c := range cases {
		t.Run(c.expected, func(t *testing.T) {
			assert.Equal(t, c.expected, c.size.Format())
		})
	}
}

type faultyWriter struct{}

func (*faultyWriter) Write([]byte) (int, error) {
	return 0, errors.New(writerErrorText)
}
