package listing

import (
	"fmt"
	"html/template"
	"io"
)

const page = `<!DOCTYPE HTML><html><head><title>File downloader</title><style>body{font-size:2em;width:60%;margin-left:auto;margin-right:auto;margin-top:80px;}</style></head><body><a href="/files/">Browse</a> <a href="/upload/">Upload</a>
    <table style="width:100%;">
      <tbody>
       <tr><td><a href="{{ .Current }}">.</a></td><td></td><td></td></tr>
       <tr><td><a href="{{ .Up }}">..</a></td><td></td><td></td></tr>
{{ range .Files }}
       <tr><td>
            <a href="{{ .Path }}">{{ .Name }}</a>
{{ if .Directory }}
              ∷ <a href="{{ .Path }}.zip?archive=true">zip</a>
{{ end }}
          </td>
{{ if .Directory }}
            <td style="width:10%;text-align:right;"><i> dir </i></td>
{{ else }}
            <td style="width:10%;text-align:right;"><i>{{ .Size.Format }}</i></td>
{{ end }}
          <td style="width:35%;font-family:monospace;text-align:right;">{{ .Mod }}</td></tr>
{{ end }}</tbody></table></body></html>`

type FileServeFormData struct {
	Files    []*FileInfo
	Current  string
	Up       string
	template *template.Template
}

func New(currentPath, oneUpPath string) *FileServeFormData {
	return &FileServeFormData{
		Current:  currentPath,
		Up:       oneUpPath,
		Files:    make([]*FileInfo, 0),
		template: template.Must(template.New("page").Parse(page)),
	}
}

func (f *FileServeFormData) AddFile(file *FileInfo) {
	f.Files = append(f.Files, file)
}

func (f *FileServeFormData) WritePage(w io.Writer) error {
	return f.template.Execute(w, f)
}

type FileInfo struct {
	Name      string
	Path      string
	Size      FileSize
	Mod       string
	Directory bool
}

type FileSize int

func (x FileSize) Format() string {
	metric := []string{"", "K", "M", "G", "T"}

	for x > 1000 && len(metric) > 1 {
		x = x / 1000
		metric = metric[1:]
	}

	return fmt.Sprintf("%d%s", x, metric[0])
}
