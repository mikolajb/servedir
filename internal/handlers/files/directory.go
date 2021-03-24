package files

import (
	"context"
	"io"
	"os"
	"path"
	"strings"
	"time"

	"github.com/mikolajb/servedir/internal/listing"
	"github.com/rs/zerolog"
)

const RootURLPath = "/files/"

type dirHandler struct {
	osPath  string
	urlPath string
}

func newDirHandler(osPath, urlPath string) targetHandler {
	return &dirHandler{
		osPath:  osPath,
		urlPath: urlPath,
	}
}

func (dh *dirHandler) Handle(ctx context.Context, w io.Writer) error {
	log := zerolog.Ctx(ctx)
	log.Info().
		Str("url1", dh.urlPath).
		Str("url2", pathOneUp(dh.urlPath)).
		Msg("listing a directory")
	listingPage := listing.New(dh.urlPath, pathOneUp(dh.urlPath))

	dirEntries, err := os.ReadDir(dh.osPath)
	if err != nil {
		log.Err(err).Msg("cannot read a directory")
		return err
	}

	for _, i := range dirEntries {
		fileInfo, err := i.Info()
		if err != nil {
			log.Err(err).Msg("error while reading file info while listing a directory")
			return err
		}

		newFileInfo := &listing.FileInfo{
			Name:      i.Name(),
			Path:      path.Join(dh.urlPath, i.Name()),
			Size:      listing.FileSize(0),
			Mod:       fileInfo.ModTime().Format(time.Stamp),
			Directory: true,
		}

		if !i.IsDir() {
			newFileInfo.Size = listing.FileSize(fileInfo.Size())
			newFileInfo.Directory = false
		}

		listingPage.AddFile(newFileInfo)
	}

	if err := listingPage.WritePage(w); err != nil {
		log.Err(err).Msg("error while rendering a listing page")
		return err
	}
	return nil
}

func pathOneUp(p string) string {
	splittedURLPath := strings.Split(p, "/")
	result := "/"
	skipped := false
	for i := len(splittedURLPath) - 1; i >= 0; i-- {
		if !skipped && len(splittedURLPath[i]) > 0 {
			skipped = true
			continue
		}

		if len(splittedURLPath[i]) == 0 {
			continue
		}

		result = "/" + splittedURLPath[i] + result
	}

	return result
}
