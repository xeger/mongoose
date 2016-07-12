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

func resolve(r parse.Resolver, pkgName string, intfs []parse.Interface) {
	for _, intf := range intfs {
		for _, meth := range intf.Methods {
			for i := 0; i < meth.Params.Len(); i++ {
				p := meth.Params.At(i)
				r.Resolve(pkgName, p.Type.Name())
			}
			for _, typ := range meth.Results {
				r.Resolve(pkgName, typ.Name())
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

func (tr *templateRenderer) Render(pkg *parse.Package, intfs []parse.Interface) (string, error) {
	out := bytes.Buffer{}
	pkgName := filepath.Base(pkg.Dir)

	resolve(tr.Resolver, pkgName, intfs)

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
