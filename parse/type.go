package parse

import (
	"fmt"
	"go/types"
	"regexp"
)

// Type represents a type, either built-in or user-defined.
type Type interface {
	Name() string
	ShortName(string, Resolver) string
	ZeroValue(string, Resolver) string
}

type loaderType struct {
	typ types.Type
}

func (lt loaderType) Name() string {
	return lt.typ.String()
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
	_, b := lt.typ.(*types.Basic)
	_, n := lt.typ.(*types.Named)

	typ := lt.typ.String()

	if b {
		if typ == "string" {
			return `""`
		} else if typ == "bool" {
			return "false"
		} else if zeroConverts.Match([]byte(typ)) {
			return fmt.Sprintf("%s(0)", typ)
		}
		panic(fmt.Sprintf("cannot handle basic type %s", typ))
	} else if n {
		return fmt.Sprintf("%s{}", lt.ShortName(local, r))
	}
	panic(fmt.Sprintf("unhandled ZeroValue for type %s", typ))
}
