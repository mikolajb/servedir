package files

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPathOneUp(t *testing.T) {
	cases := []struct {
		path  string
		oneUp string
	}{
		{
			path:  "/",
			oneUp: "/",
		},
		{
			path:  "/files/",
			oneUp: "/",
		},
		{
			path:  "/files",
			oneUp: "/",
		},
		{
			path:  "/foo/bar",
			oneUp: "/foo/",
		},
	}

	for _, c := range cases {
		t.Run(c.path, func(t *testing.T) {
			assert.Equal(t, c.oneUp, pathOneUp(c.path))
		})
	}
}
