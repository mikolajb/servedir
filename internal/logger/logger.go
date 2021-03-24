package logger

import (
	"context"
	"io"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/justinas/alice"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/hlog"
)

func NewLoggerForHTTPHandler(ctx context.Context, handler http.Handler) http.Handler {
	log, _ := NewLogger(ctx)

	c := alice.New(
		hlog.NewHandler(log),
		hlog.UserAgentHandler("userAgent"),
		hlog.MethodHandler("method"),
		hlog.RemoteAddrHandler("remoteAddress"),
	)

	return c.Then(handler)
}

// NewLogger instantiates a new logger (console writer for now).
func NewLogger(ctx context.Context) (zerolog.Logger, context.Context) {
	level := zerolog.InfoLevel

	if os.Getenv("DEBUG") == "on" {
		level = zerolog.DebugLevel
	}

	logger := zerolog.New(newConsoleWriter(os.Stderr, false)).Level(level)

	return logger, logger.WithContext(ctx)
}

func newConsoleWriter(writer io.Writer, noColor bool) zerolog.ConsoleWriter {
	return zerolog.NewConsoleWriter(
		func(w *zerolog.ConsoleWriter) {
			w.Out = writer
			w.NoColor = noColor
			w.TimeFormat = time.RFC822
		},
	)
}

// NewTestLogger insatiates a new logger that uses testing.T as a backend.
func NewTestLogger(t testing.TB, ctx context.Context) (zerolog.Logger, context.Context) {
	zerolog.SetGlobalLevel(zerolog.DebugLevel)
	logger := zerolog.New(newConsoleWriter(testWriter{TB: t}, true))

	return logger, logger.WithContext(ctx)
}

type testWriter struct {
	testing.TB
}

func (tw testWriter) Write(p []byte) (int, error) {
	tw.Helper()
	tw.Logf("%s", p)
	return len(p), nil
}
