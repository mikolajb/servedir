package files

import (
	"context"
	"errors"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/hlog"
)

const errorMessage = "<!DOCTYPE HTML><html><head><title>File uploader</title><style>body{font-size:2em;width: 60%;margin-left:auto;margin-right:auto;margin-top:80px;}</style></head><body>Error</body></html>"

var ErrNoZipSuffix = errors.New("path doesn't have '.zip' suffix")

func New(path string) http.Handler {
	return &handler{
		path: path,
	}
}

type handler struct {
	path string
}

func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	currentOSPath := filepath.Join(
		h.path,
		filepath.Join(
			strings.Split(
				strings.TrimPrefix(r.URL.Path, RootURLPath),
				"/",
			)...,
		),
	)
	log := hlog.FromRequest(r)

	handler, err := getTargetHandler(ctx, r.URL.Query(), r.URL.Path, currentOSPath)
	if err != nil {
		log.Err(err).Msg("error while creating a target handler")
		http.Error(w, errorMessage, http.StatusInternalServerError)
		return
	}
	if err := handler.Handle(ctx, w); err != nil {
		http.Error(w, errorMessage, http.StatusInternalServerError)
		return
	}
}

func getTargetHandler(ctx context.Context, queryValues url.Values, urlPath, osPath string) (targetHandler, error) {
	log := zerolog.Ctx(ctx).With().Str("urlPath", urlPath).Str("osPath", osPath).Logger()
	log.Debug().Msg("obtaining target for a request")

	handleArchive := queryValues.Get("archive") == "true"

	if handleArchive {
		if !strings.HasSuffix(osPath, ".zip") {
			log.Warn().Msg("path doesn't have '.zip' suffix")
			return nil, ErrNoZipSuffix
		}
		osPath = osPath[0 : len(osPath)-4]
	}

	fileInfo, err := os.Stat(osPath)
	if err != nil {
		return nil, err
	}

	if fileInfo.IsDir() {
		if handleArchive {
			return newZipHandler(osPath), nil
		} else {
			return newDirHandler(osPath, urlPath), nil
		}
	}
	return newSingleFile(osPath), nil
}
