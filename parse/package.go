package parse

import (
	"fmt"
	"go/ast"
	"go/build"
	"go/importer"
	"go/types"
	"path/filepath"

	"golang.org/x/tools/go/loader"
)

// Package is a package that contains interfaces.
type Package struct {
	// Absolute path to this package
	Dir string
	// All package names declared by any source file (usually just 1 or 2!)
	Names []string
	// All interfaces defined in any source file
	Interfaces []Interface
}

// Name is the "natural" name of this package (i.e. base name of path it's located in).
func (lp *Package) Name() string {
	return filepath.Base(lp.Dir)
}

// Len is the number of interfaces defined in the sources.
func (lp *Package) Len() int {
	return len(lp.Interfaces)
}

// String returns a pseudocode package definition.
func (lp *Package) String() string {
	return fmt.Sprintf("package %s {%v}", lp.Name(), lp.Interfaces)
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
		// swallow the error for now; they end up getting reported as an error return
		//fmt.Fprintf(os.Stderr, "    %s\n", err.Error())
	}
	conf.TypeChecker.Importer = importer.Default()

	for _, fn := range buildPkg.GoFiles {
		fp := filepath.Join(path, fn)
		f, errp := conf.ParseFile(fp, nil)

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
