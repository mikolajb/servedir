package upload

import (
	"html/template"
	"io"
)

const page = `<!DOCTYPE HTML><html><head><title>File uploader</title><style>body{font-size:2em;width: 60%;margin-left:auto;margin-right:auto;margin-top:80px;}</style></head><body><a href="/files/">Browse</a> <a href="/upload/">Upload</a>
<a href="/upload/?files={{ .OneMore }}">One more</a>
<a href="/upload/?files={{ .OneLess }}">One less</a>
<form enctype="multipart/form-data" action="/upload/" method="post">
{{ range .Items }}
      <input type="file" name="files" style="width:100%">
{{ end }}
      <input style="width:9%" type="submit" value="Upload">
</form></body></html>
`

type UploadFormData struct {
	OneMore  uint
	OneLess  uint
	Items    []uint
	template *template.Template
}

func New(files uint) *UploadFormData {
	if files == 0 {
		files = 1
	}

	if files > 20 {
		files = 20
	}

	return &UploadFormData{
		OneMore:  uint(files + 1),
		OneLess:  uint(files - 1),
		Items:    make([]uint, files),
		template: template.Must(template.New("page").Parse(page)),
	}
}

func (f *UploadFormData) WritePage(w io.Writer) error {
	return f.template.Execute(w, f)
}
