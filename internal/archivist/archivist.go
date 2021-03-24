package archivist

import (
	"archive/zip"
	"io"
)

type Archivist interface {
	Add(path string, data io.Reader) error
	Close() error
}

type archivist struct {
	writer *zip.Writer
}

func New(output io.Writer) Archivist {
	return &archivist{
		writer: zip.NewWriter(output),
	}
}

func (a *archivist) Add(path string, data io.Reader) error {
	if len(path) == 0 {
		return nil
	}

	f, err := a.writer.Create(path)
	if err != nil {
		return err
	}

	_, err = io.Copy(f, data)
	if err != nil {
		return err
	}

	return nil
}

func (a *archivist) Close() error {
	return a.writer.Close()
}
