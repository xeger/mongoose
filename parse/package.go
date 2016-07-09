package parse

import (
	"fmt"
	"go/ast"
	"go/build"
	"go/importer"
	"go/types"
	"log"
	"path/filepath"

	"golang.org/x/tools/go/loader"
)

// Package represents a package, which contains types.
type Package interface {
	// Absolute path to directory containing package source files.
	Dir() string
	// Number of interfaces.
	Len() int
	// All names declared in the package statement of any source file in this
	// package. For valid packages, should only contain or two entries: foo
	// and/or foo_test.
	EachName(func(string))
	// All interfaces declared within this package's source files.
	EachInterface(func(Interface))
}

type loaderPackage struct {
	dir        string
	names      []string
	interfaces []Interface
}

func (lp *loaderPackage) Dir() string {
	return lp.dir
}

func (lp *loaderPackage) Len() int {
	return len(lp.interfaces)
}

func (lp *loaderPackage) EachName(cb func(name string)) {
	for _, n := range lp.names {
		cb(n)
	}
}

func (lp *loaderPackage) EachInterface(cb func(Interface)) {
	for _, intf := range lp.interfaces {
		cb(intf)
	}
}

func (lp *loaderPackage) String() string {
	return fmt.Sprintf("Package(size=%d)", lp.Len())
}

func (lp *loaderPackage) finalize(interfaces map[string]*types.Interface) {
	lp.interfaces = make([]Interface, 0, len(interfaces))
	for name, intf := range interfaces {
		li := &loaderInterface{}
		li.finalize(name, intf)
		lp.interfaces = append(lp.interfaces, li)
	}
}

// NewPackage loads *.go from a given directory and returns type information about the package defined in it.
func NewPackage(path string) (Package, error) {
	path, err := filepath.Abs(path)
	if err != nil {
		return nil, err
	}

	buildPkg, err := build.Default.Import(".", path, build.AllowBinary|build.ImportComment)
	if err != nil {
		return nil, err
	}

	var astFiles []*ast.File
	pkgNames := map[string]bool{}
	var conf loader.Config

	conf.TypeCheckFuncBodies = func(_ string) bool { return false }
	conf.TypeChecker.DisableUnusedImportCheck = true
	conf.TypeChecker.Error = func(err error) {
		// TODO: something else about this error?
		log.Println("parse error", err.Error())
	}
	conf.TypeChecker.Importer = importer.Default()

	for _, fn := range buildPkg.GoFiles {
		f, errp := conf.ParseFile(filepath.Join(path, fn), nil)

		if errp != nil {
			return nil, errp
		}

		astFiles = append(astFiles, f)
		pkgNames[f.Name.Name] = true
	}

	conf.CreateFromFiles(path, astFiles...)

	prog, err := conf.Load()
	if err != nil {
		return nil, err
	}

	interfaces := map[string]*types.Interface{}
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
			interfaces[name] = iface
		}
	}

	pkgList := make([]string, 0, len(pkgNames))
	for k := range pkgNames {
		pkgList = append(pkgList, k)
	}

	lp := &loaderPackage{}
	lp.dir = path
	lp.names = pkgList
	lp.finalize(interfaces)
	return lp, nil
}
