package main

import (
	"net/http"
	"fmt"
	"io"
	"io/ioutil"
	"time"
	"html/template"
	"flag"
	"path/filepath"
	"strings"
	"os"
	"archive/zip"
	"strconv"
)

func archivist(w io.Writer) (func(string, io.Reader), func()) {
	zipWriter := zip.NewWriter(w)

	return func(path string, data io.Reader){
		fmt.Println("compressing: ", path)

		if len(path) == 0 {
			return
		}

		f, _ := zipWriter.Create(path)
		_, err := io.Copy(f, data)

		if err != nil {
			fmt.Println(err.Error())
		}

		return
	}, func() {
		err := zipWriter.Close()

		if err != nil {
			fmt.Println(err.Error())
			fmt.Fprint(w, ErrorMessage)
		}

		return
	}
}

var flagPath string
var flagPort uint

type FileInfo struct {
	Name string
	Path string
	Size FileSize
	Mod string
	Directory bool
}

type FileSize int
func (x FileSize) Format() string {
	metric := []string{"", "K", "M", "G", "T"}

	for x > 1000 {
		x = x / 1000
		metric = metric[1:]
	}

	return fmt.Sprintf("%d%s", x, metric[0])
}

type FileServeFormData struct {
	Files []FileInfo
	Current string
	Up string
}

type FileServeHandler string

const ErrorMessage = "<!DOCTYPE HTML><html><head><title>File uploader</title><style>body{font-size:2em;width: 60%;margin-left:auto;margin-right:auto;margin-top:80px;}</style></head><body>Error</body></html>"

func (fileServeHandler FileServeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	splitted := strings.Split(r.URL.Path, "/")
	relpath := filepath.Join(string(fileServeHandler),
		strings.Join(splitted,
			string(filepath.Separator)))
	fmt.Println(r.URL.Path, relpath)

	if r.URL.Query().Get("archive") == "true" {
		if !strings.HasSuffix(relpath, ".zip") {
			fmt.Println("path doesn't have '.zip' suffix")
			fmt.Fprint(w, ErrorMessage)
			return
		}

		fmt.Println(relpath)
		relpath = relpath[0:len(relpath) - 4]
		fmt.Println(relpath)
	}

	fileInfo, err := os.Stat(relpath)

	if err != nil {
		fmt.Println(err.Error())
		fmt.Fprint(w, ErrorMessage)
		return
	}

	if !fileInfo.IsDir() {
		file, err := os.Open(relpath)

		if err != nil {
			fmt.Println(err.Error())
			fmt.Fprint(w, ErrorMessage)
			return
		}

		io.Copy(w, file)
		return
	}

	if r.URL.Query().Get("archive") == "true" {
		add, result := archivist(w)

		err = filepath.Walk(relpath,
			func (path string, info os.FileInfo, err error) error {
				archRelpath, err := filepath.Rel(relpath, path)

				if info.IsDir() {
					return nil
				}

				file, err := os.Open(path)
				add(archRelpath, file)
				return nil
			})
		result()
		return
	}

	formData := FileServeFormData{}
	formData.Current = fmt.Sprintf("/files/%s", r.URL.Path)
	formData.Up = "/files/"

	if len(splitted) > 1 {
		formData.Up = fmt.Sprintf("%s%s",
			formData.Up,
			strings.Join(splitted[:len(splitted)-1], "/"))
	}

	res, err := ioutil.ReadDir(relpath)

	if err != nil {
		fmt.Println(err)
		fmt.Fprint(w, ErrorMessage)
	}

	for _, i := range(res) {
		newFileInfo := FileInfo{i.Name(),
			filepath.Join(formData.Current, i.Name()),
			FileSize(0),
			i.ModTime().Format(time.Stamp),
			true}

		if !i.IsDir() {
			newFileInfo.Size = FileSize(i.Size())
			newFileInfo.Directory = false
		}

		formData.Files = append(formData.Files, newFileInfo)
	}


	tmpl, err := template.New("page").Parse(`<!DOCTYPE HTML><html><head><title>File downloader</title><style>body{font-size:2em;width:60%;margin-left:auto;margin-right:auto;margin-top:80px;}</style></head><body><a href="/files/">Browse</a> <a href="/upload/">Upload</a>
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
{{ end }}</tbody></table></body></html>`)

	if err != nil { panic(err) }
	err = tmpl.Execute(w, formData)
	if err != nil { panic(err) }
}

type UploadFormData struct {
	OneMore uint64
	OneLess uint64
	Items []byte
}

func main() {
	flag.StringVar(&flagPath, "path", ".", "Path")
	flag.UintVar(&flagPort, "port", 8080, "Port")
	flag.Parse()

	fmt.Println("Serving directory:", flagPath, "on port:", flagPort)

	http.Handle("/files/", http.StripPrefix("/files/", FileServeHandler(flagPath)))
	http.HandleFunc("/upload/", func(w http.ResponseWriter, r *http.Request) {

		if r.Method == "POST" {
			r.ParseMultipartForm(32 << 20)

			for _, file := range r.MultipartForm.File["files"] {
				f, _ := file.Open()
				fmt.Println(file.Filename)
				_, err := os.Stat(file.Filename)

				if os.IsExist(err) {
					fmt.Println("File '", file.Filename, "' alerady exist.")
					continue
				}

				data, _ := ioutil.ReadAll(f)
				ioutil.WriteFile(file.Filename, data, 0777)
				f.Close()
			}

			fmt.Fprint(w, "<!DOCTYPE HTML><html><head><title>File uploader</title><style>body{font-size:2em;width: 60%;margin-left:auto;margin-right:auto;margin-top:80px;}</style></head><body>Sent</body></html>")
			return
		}

		tmpl, err := template.New("upload").Parse(`<!DOCTYPE HTML><html><head><title>File uploader</title><style>body{font-size:2em;width: 60%;margin-left:auto;margin-right:auto;margin-top:80px;}</style></head><body><a href="/files/">Browse</a> <a href="/upload/">Upload</a>
<a href="/upload/?files={{ .OneMore }}">One more</a>
<a href="/upload/?files={{ .OneLess }}">One less</a>
<form enctype="multipart/form-data" action="/upload/" method="post">
      {{ range .Items }}
      <input type="file" name="files" style="width:100%">
      {{ end }}
      <input style="width:9%" type="submit" value="Upload">
</form></body></html>
`)
		files := uint64(1)

		if len(r.URL.Query().Get("files")) > 0 {
			files, err = strconv.ParseUint(r.URL.Query().Get("files"), 10, 0)

			if err != nil {
				fmt.Println(err)
				files = uint64(1)
			}
		}

		if files > 20 {	files = 20 }
		if files == 0 { files = 1 }

		formData := make([]byte, files)
		err = tmpl.Execute(w, UploadFormData{files + 1, files -1, formData})
		if err != nil { panic(err) }
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
				fmt.Fprint(w, `<!DOCTYPE HTML><html><head><title>File uploader</title><style>body{font-size:2em;width:60%;margin-left:auto;margin-right:auto;margin-top:80px;}</style></head><body><a href="/files/">Browse</a> <a href="/upload/">Upload</a></body></html>`)})

	err := http.ListenAndServe(fmt.Sprintf(":%d", flagPort), nil)

	if err != nil { panic(err) }
}
