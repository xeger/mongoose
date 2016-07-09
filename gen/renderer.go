package gen

import (
	"bytes"
	"fmt"
	"path/filepath"

	"github.com/xeger/mongoose/parse"
)

// Renderer generates the complete contents of a Go source file, including
// imports, that defines mocks and/or stubs for a number of interfaces.
type Renderer interface {
	Render(*parse.Package, []parse.Interface) (string, error)
}

// BasicRenderer uses sprintf and stuff to generate a behavior-less stub
// for every mocked interface. It's not very useful and serves mostly as a
// simple proof-of-concept.
type StubRenderer struct{}

func (StubRenderer) Render(pkg *parse.Package, intfs []parse.Interface) (string, error) {
	pkgName := filepath.Base(pkg.Dir)
	r := parse.NewResolver()

	out := bytes.Buffer{}

	// Resolve all imports before generating any code.
	for _, intf := range intfs {
		for _, meth := range intf.Methods {
			for _, typ := range meth.Params {
				r.Resolve(pkgName, typ.Name())
			}
			for _, typ := range meth.Results {
				r.Resolve(pkgName, typ.Name())
			}
		}
	}

	// Generate import statement.
	fmt.Fprintln(&out, "import (")
	r.EachImport(func(pkg, name string) {
		fmt.Fprintln(&out, "  ", name, fmt.Sprintf(`"%s"`, pkg))
	})
	fmt.Fprintln(&out, ")")

	for _, intf := range intfs {
		mock := fmt.Sprintf("Mock%s", intf.Name)

		// Mock type definition.
		fmt.Fprintln(&out)
		fmt.Fprintln(&out, "type", mock, "struct {")
		fmt.Fprintln(&out, "}")
		fmt.Fprintln(&out)

		// Mock method implementations
		for _, meth := range intf.Methods {
			// Declare method parameters
			buf := bytes.Buffer{}
			for name, typ := range meth.Params {
				if buf.Len() > 0 {
					buf.WriteString(",")
				}
				buf.WriteString(fmt.Sprintf("%s %s", name, typ.ShortName(pkgName, r)))
			}
			params := buf.String()

			// Declare returns (if any)
			buf = bytes.Buffer{}
			multi := false
			for _, typ := range meth.Results {
				if buf.Len() > 0 {
					buf.WriteString(",")
					multi = true
				}
				buf.WriteString(typ.ShortName(pkgName, r))
			}
			var results string
			if multi {
				results = fmt.Sprintf("(%s)", buf.String())
			} else {
				results = buf.String()
			}

			fmt.Fprintln(&out)
			fmt.Fprintf(&out, "func (m *%s) %s(%s) %s {\n", mock, meth.Name, params, results)
			multi = false
			buf = bytes.Buffer{}
			for _, typ := range meth.Results {
				if buf.Len() > 0 {
					buf.WriteString(",")
					multi = true
				}
				buf.WriteString(typ.ZeroValue(pkgName, r))
			}
			var zeroes string
			if multi {
				zeroes = fmt.Sprintf("(%s)", buf.String())
			} else {
				zeroes = buf.String()
			}
			fmt.Fprintln(&out, "  return", zeroes)
			fmt.Fprintln(&out, "}")
		}
	}

	return out.String(), nil
}
