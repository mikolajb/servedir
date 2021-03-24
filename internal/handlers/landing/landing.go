package landing

import (
	"fmt"
	"net/http"

	"github.com/rs/zerolog/hlog"
)

type handler struct{}

func New() http.Handler {
	return &handler{}
}

func (*handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log := hlog.FromRequest(r)
	log.Info().Msg("serving landing page")
	fmt.Fprint(w, `<!DOCTYPE HTML><html><head><title>File uploader</title><style>body{font-size:2em;width:60%;margin-left:auto;margin-right:auto;margin-top:80px;}</style></head><body><a href="/files/">Browse</a> <a href="/upload/">Upload</a></body></html>`)
}
