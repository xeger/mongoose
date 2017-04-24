package gen

import (
	"bytes"
	"path/filepath"
	"text/template"

	"github.com/xeger/mongoose/parse"
)

// Renderer generates the complete contents of a Go source file, including
// imports, that defines mocks and/or stubs for a number of interfaces.
type Renderer interface {
	Render(*parse.Package, []parse.Interface) (string, error)
}

// Resolve all types in a list of interfaces. Useful to generate all of the
// imports we'll require before we actually generating any code.
func resolve(local string, r parse.Resolver, intfs []parse.Interface) {
	for _, intf := range intfs {
		for _, meth := range intf.Methods {
			for i := 0; i < meth.Params.Len(); i++ {
				meth.Params.At(i).Type.ShortName(local, r)
			}
			if v := meth.Params.Variadic(); v != nil {
				v.Type.ShortName(local, r)
			}
			for _, typ := range meth.Results {
				typ.ShortName(local, r)
			}
		}
	}
}

type templateRenderer struct {
	// Resolver stores import names -> packages
	Resolver parse.Resolver
	// Header is called before items for import statements, etc.
	Header *template.Template
	// Item is called once for each interface in the file
	Item *template.Template
}

type headerContext struct {
	Resolver parse.Resolver
	Package  *parse.Package
}

type itemContext struct {
	Resolver  parse.Resolver
	Package   *parse.Package
	Interface *parse.Interface
}

const templateHeader = `package {{.Package.Name}}

import ({{range $nick, $pkg := .Resolver.Imports}}
	{{$nick}} "{{$pkg}}"{{end}}
)
`

// newTemplateRenderer initializes a renderer, its Resolver, and its Header.
// Other fields must be initialized by the caller before rendering.
func newTemplateRenderer() *templateRenderer {
	tr := &templateRenderer{}
	tr.Resolver = parse.NewResolver()
	tr.Header = template.New("templateHeader")
	tr.Header.Parse(templateHeader)
	return tr
}

func (tr *templateRenderer) Render(pkg *parse.Package, intfs []parse.Interface) (string, error) {
	out := bytes.Buffer{}
	local := filepath.Base(pkg.Dir)

	resolve(local, tr.Resolver, intfs)

	if tr.Header != nil {
		err := tr.Header.Execute(&out, headerContext{Resolver: tr.Resolver, Package: pkg})
		if err != nil {
			return "", err
		}
	}
	if tr.Item != nil {
		for _, intf := range intfs {
			err := tr.Item.Execute(&out, itemContext{Resolver: tr.Resolver, Package: pkg, Interface: &intf})
			if err != nil {
				return "", err
			}
		}
	}

	return out.String(), nil
}
