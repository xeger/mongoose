package gen

import (
	"fmt"
	"path/filepath"

	"github.com/acsellers/inflections"
	"github.com/xeger/mongoose/parse"
)

// Placer decides where to write the source code for an interface's mock.
// TODO ---- stop passing a package in; just pass a dir in which the package's source is defined ----
type Placer interface {
	Place(path string, intf *parse.Interface) string
}

// PakagePlacer puts each interface's mock in its own file located in the same
// package as the interface.
// Example: the mock of foo.Xyz is placed in foo/mock_xyz.go
type PackagePlacer struct{}

func (PackagePlacer) Place(path string, intf *parse.Interface) string {
	return filepath.Join(path, fmt.Sprintf("mock_%s.go", inflections.Underscore(intf.Name)))
}
