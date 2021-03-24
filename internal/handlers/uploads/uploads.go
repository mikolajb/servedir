package uploads

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"

	"github.com/mikolajb/servedir/internal/upload"
	"github.com/rs/zerolog/hlog"
)

const (
	errorMessage  = "<!DOCTYPE HTML><html><head><title>File uploader</title><style>body{font-size:2em;width: 60%;margin-left:auto;margin-right:auto;margin-top:80px;}</style></head><body>Error</body></html>"
	confimMessage = "<!DOCTYPE HTML><html><head><title>File uploader</title><style>body{font-size:2em;width: 60%;margin-left:auto;margin-right:auto;margin-top:80px;}</style></head><body>Sent</body></html>"
)

type handler struct{}

func New() http.Handler {
	return &handler{}
}

func (*handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log := hlog.FromRequest(r)
	log.Info().Msg("rendering upload form")

	switch r.Method {
	case "POST":
		if err := r.ParseMultipartForm(32 << 20); err != nil {
			log.Err(err).Msg("error while parsing a form")
			return
		}

		for _, file := range r.MultipartForm.File["files"] {
			f, _ := file.Open()
			log.Debug().Msgf("accepting a file: %s", file.Filename)

			_, err := os.Stat(file.Filename)
			if os.IsExist(err) {
				log.Warn().Msgf("file %s already exists, skipping", file.Filename)
				continue
			}

			outputFile, err := os.OpenFile(file.Filename, os.O_CREATE|os.O_WRONLY, os.ModePerm)
			if err != nil {
				log.Err(err).Msgf("error while creating an output file: %s", file.Filename)
				return
			}
			written, err := io.Copy(outputFile, f)
			if err != nil {
				log.Err(err).Msg("error while writing to a file")
				return
			}
			log.Debug().Msgf("%d bytes wrote", written)
			if err := f.Close(); err != nil {
				log.Err(err).Msg("error while closing a file")
				return
			}
			if err := outputFile.Close(); err != nil {
				log.Err(err).Msg("error while closing an output file")
				return
			}
		}

		fmt.Fprint(w, confimMessage)
	case "GET":
		var files uint
		if len(r.URL.Query().Get("files")) > 0 {
			fileNumber, err := strconv.ParseUint(r.URL.Query().Get("files"), 10, 0)
			if err != nil {
				log.Err(err).Msg("error while processing URL query")
				http.Error(w, errorMessage, http.StatusBadRequest)
				return
			}
			files = uint(fileNumber)
		}
		uploadPage := upload.New(files)
		if err := uploadPage.WritePage(w); err != nil {
			http.Error(w, errorMessage, http.StatusInternalServerError)
			return
		}
	}
}
