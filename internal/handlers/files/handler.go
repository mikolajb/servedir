package files

import (
	"context"
	"io"
)

type targetHandler interface {
	Handle(context.Context, io.Writer) error
}
