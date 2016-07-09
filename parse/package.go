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

type Package struct {
	Dir        string
	Names      []string
	Interfaces []Interface
}

func (lp *Package) Len() int {
	return len(lp.Interfaces)
}

func (lp *Package) String() string {
	return fmt.Sprintf("Package(size=%d)", lp.Len())
}

func (lp *Package) finalize(interfaces map[string]*types.Interface) {
	lp.Interfaces = make([]Interface, 0, len(interfaces))
	for name, gointf := range interfaces {
		intf := Interface{}
		intf.finalize(name, gointf)
		lp.Interfaces = append(lp.Interfaces, intf)
	}
}

// NewPackage loads *.go from a given directory and returns type information about the package defined in it.
func NewPackage(path string) (*Package, error) {
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

	lp := &Package{}
	lp.Dir = path
	lp.Names = pkgList
	lp.finalize(interfaces)
	return lp, nil
}
