package files

import (
	"context"
	"io"
	"os"

	"github.com/rs/zerolog"
)

type singleFile struct {
	path string
}

func (sf *singleFile) Handle(ctx context.Context, w io.Writer) error {
	log := zerolog.Ctx(ctx)
	log.Info().Msg("serving a single file")

	file, err := os.Open(sf.path)
	if err != nil {
		log.Err(err).Msg("error while opening a file")
		return err
	}

	if _, err := io.Copy(w, file); err != nil {
		log.Err(err).Msg("error while reading a file")
		return err
	}

	return nil
}

func newSingleFile(path string) targetHandler {
	return &singleFile{
		path: path,
	}
}
