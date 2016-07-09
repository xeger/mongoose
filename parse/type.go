package parse

import (
	"fmt"
	"go/types"
	"path/filepath"
	"regexp"
	"strings"
)

// Type represents a type, either built-in or user-defined.
type Type interface {
	// Fully-qualified name of this type. Consists of the absolute path to
	// the type's package on disk, a dot, and the type's name within that
	// package.
	// Example: /usr/local/src/awesome.Widget
	Name() string

	// BareName provides the type's name without any package prefix.
	BareName() string

	// Name that can be used to refer to this type in Go source code. ShortName
	// is subjective and depends on the context in which it is used, hence the
	// need for some method parameters to provide that context.
	//
	// The first parameter is the package statement from the "local" source,
	// used to check whether the type is declared in the same package as the
	// source.
	//
	// The second parameter is a Resolver that can be used to map this type's
	// package to an import name within the source.
	ShortName(string, Resolver) string

	// Zero value that can be used to create literals of this type in Go source
	// code.
	ZeroValue(string, Resolver) string
}

type loaderType struct {
	typ types.Type
}

func (lt loaderType) Name() string {
	return lt.typ.String()
}

func (lt loaderType) BareName() string {
	pkgWithName := filepath.Base(lt.Name())
	parts := strings.SplitN(pkgWithName, ".", 2)
	return parts[len(parts)-1]
}

func (lt loaderType) ShortName(local string, r Resolver) string {
	_, b := lt.typ.(*types.Basic)
	_, n := lt.typ.(*types.Named)

	typ := lt.typ.String()

	if b {
		// basic types: nothing to do
		return typ
	} else if n {
		return r.Resolve(local, typ)
	}
	panic(fmt.Sprintf("unhandled ShortName for type %s", typ))
}

var zeroConverts = regexp.MustCompile("(byte|u?int|float|rune)[0-9]*")

func (lt loaderType) ZeroValue(local string, r Resolver) string {
	basic, b := lt.typ.(*types.Basic)
	named, n := lt.typ.(*types.Named)

	typ := lt.typ.String()

	if b {
		return lt.zeroBasic(basic)
	} else if n {
		under := named.Underlying()
		// It's a named type; check underlying type.
		_, i := under.(*types.Interface)
		_, s := under.(*types.Struct)
		basic, b = under.(*types.Basic)

		if i {
			// nil interface value
			return "nil"
		} else if s {
			// zero-valued struct literal
			return fmt.Sprintf("%s{}", lt.ShortName(local, r))
		} else if b {
			// conversion of underlying basic zero
			return fmt.Sprintf("%s(%s)", lt.ShortName(local, r), lt.zeroBasic(basic))
		}
		panic(fmt.Sprintf("unhandled ZeroValue for type %s", typ))
	}

	panic(fmt.Sprintf("unhandled ZeroValue for type %s", typ))
}

func (lt loaderType) zeroBasic(basic *types.Basic) string {
	typ := basic.String()
	if typ == "string" {
		return `""`
	} else if typ == "bool" {
		return "false"
	} else if zeroConverts.Match([]byte(typ)) {
		return fmt.Sprintf("%s(0)", typ)
	}
	panic(fmt.Sprintf("cannot handle basic type %s", typ))
}
