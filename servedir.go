package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"

	"github.com/mikolajb/servedir/internal/handlers/files"
	"github.com/mikolajb/servedir/internal/handlers/landing"
	"github.com/mikolajb/servedir/internal/handlers/uploads"
	"github.com/mikolajb/servedir/internal/logger"
)

var flagPath string
var flagPort uint

func main() {
	flag.StringVar(&flagPath, "path", ".", "Path")
	flag.UintVar(&flagPort, "port", 8080, "Port")
	flag.Parse()

	log, ctx := logger.NewLogger(context.Background())
	log.Info().Uint("port", flagPort).Str("os root path", flagPath).Msgf("serving directory %s on port %d", flagPath, flagPort)

	http.Handle(files.RootURLPath,
		logger.NewLoggerForHTTPHandler(
			ctx,
			files.New(flagPath),
		),
	)
	http.Handle("/upload/",
		logger.NewLoggerForHTTPHandler(
			ctx,
			uploads.New(),
		),
	)
	http.Handle("/",
		logger.NewLoggerForHTTPHandler(
			ctx,
			landing.New(),
		),
	)

	if err := http.ListenAndServe(fmt.Sprintf(":%d", flagPort), nil); err != nil {
		panic(err)
	}
}
