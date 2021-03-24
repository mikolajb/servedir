package files

import (
	"context"
	"io"
	"os"
	"path/filepath"

	"github.com/mikolajb/servedir/internal/archivist"
	"github.com/rs/zerolog"
)

type zipHandler struct {
	path string
}

func newZipHandler(path string) targetHandler {
	return &zipHandler{
		path: path,
	}
}

func (zh *zipHandler) Handle(ctx context.Context, w io.Writer) error {
	log := zerolog.Ctx(ctx)

	log.Info().Msg("serving a zip file")
	archivist := archivist.New(w)

	err := filepath.Walk(zh.path,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				log.Err(err).Msg("filepath Walk error")
				return err
			}

			archRelpath, err := filepath.Rel(zh.path, path)
			if err != nil {
				log.Err(err).Msg("error while resolving a relative path")
				return err
			}

			if info.IsDir() {
				return nil
			}

			file, err := os.Open(path)
			if err != nil {
				log.Err(err).Msg("error while opening a file")
				return err
			}

			return archivist.Add(archRelpath, file)
		})
	if err != nil {
		log.Err(err).Msg("error while traversing a directory")
		return err
	}

	if err := archivist.Close(); err != nil {
		log.Err(err).Msg("archivist error on closing")
		return err
	}

	return err
}
