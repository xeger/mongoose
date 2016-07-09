package gen

import (
	"fmt"
	"path/filepath"

	"github.com/acsellers/inflections"
	"github.com/xeger/mongoose/parse"
)

// Placer decides where to write the source code for an interface's mock.
type Placer interface {
	Place(pkg parse.Package, intf parse.Interface) string
}

// Places each interface's mock in its own file located in the same
// package as the interface.
// Example: the mock of foo.Xyz is placed in foo/mock_xyz.go
type PackagePlacer struct{}

func (PackagePlacer) Place(pkg parse.Package, intf parse.Interface) string {
	path := pkg.Dir
	return filepath.Join(path, fmt.Sprintf("mock_%s.go", inflections.Underscore(intf.Name)))
}
