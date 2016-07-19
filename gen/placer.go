package gen

import (
	"fmt"
	"path/filepath"

	"github.com/acsellers/inflections"
	"github.com/xeger/mongoose/parse"
)

// Placer decides where to write the source code for an interface's mock.
type Placer interface {
	Place(path string, intf *parse.Interface) string
}

// PackagePlacer puts mocks into the same package as their interfaces, either in
// a single file (if FilePerInterface is false) or one mock per file (if true).
//
// Without FilePerInterface: the mock of foo.Bar is placed in foo/mocks.go
// With FilePerInterface: the mock of foo.Bar is placed in foo/mock_bar.go
type PackagePlacer struct {
	FilePerInterface bool
}

// Place decides on the filename and directory that an interface's mock should
// be written to.
func (pp PackagePlacer) Place(path string, intf *parse.Interface) string {
	if pp.FilePerInterface {
		return filepath.Join(path, fmt.Sprintf("mock_%s.go", inflections.Underscore(intf.Name)))
	}
	return filepath.Join(path, "mocks.go")
}
