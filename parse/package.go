package parse

import (
	"fmt"
	"go/ast"
	"go/importer"
	"go/types"
	"io/ioutil"
	"log"
	"path/filepath"

	"golang.org/x/tools/go/loader"
)

// Package represents a package, which contains types.
type Package interface {
	Size() int
	EachInterface(func(Interface))
}

type loaderPackage struct {
	interfaces map[string]*types.Interface
}

func (lp loaderPackage) Size() int {
	return len(lp.interfaces)
}

func (lp loaderPackage) EachInterface(cb func(Interface)) {
	for n, i := range lp.interfaces {
		cb(loaderInterface{n,i})
	}
}

func (lp loaderPackage) String() string {
	return fmt.Sprintf("Package(size=%d)", lp.Size())
}

// NewPackage loads *.go from a given directory and returns type information about the package defined in it.
func NewPackage(path string) (Package, error) {
	path, err := filepath.Abs(path)
	if err != nil {
		return nil, err
	}

	files, err := ioutil.ReadDir(path)
	if err != nil {
		return nil, err
	}

	var astFiles []*ast.File
	var conf loader.Config

	conf.TypeCheckFuncBodies = func(_ string) bool { return false }
	conf.TypeChecker.DisableUnusedImportCheck = true
	conf.TypeChecker.Error = func(err error) {
		// TODO: something else about this error?
		log.Println("parse error", err.Error())
	}
	conf.TypeChecker.Importer = importer.Default()

	for _, fi := range files {
		if filepath.Ext(fi.Name()) != ".go" {
			continue
		}

		fpath := filepath.Join(path, fi.Name())
		f, errp := conf.ParseFile(fpath, nil)
		if errp != nil {
			return nil, errp
		}

		astFiles = append(astFiles, f)
	}

	conf.CreateFromFiles(path, astFiles...)

	prog, err := conf.Load()
	if err != nil {
		return nil, err
	}

	lp := &loaderPackage{interfaces: map[string]*types.Interface{}}
	for _, pinfo := range prog.Created {
		scope := pinfo.Pkg.Scope()
		for _, name := range scope.Names() {
			obj := scope.Lookup(name)
			if obj == nil {
				continue
			}
			typ, ok := obj.Type().(*types.Named)
			if !ok {
				continue
			}
			name = typ.Obj().Name()
			iface, ok := typ.Underlying().(*types.Interface)
			if !ok {
				continue
			}
			iface = iface.Complete()
			lp.interfaces[name] = iface
		}
	}
	return lp, nil
}
