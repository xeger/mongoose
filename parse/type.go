package parse

import (
	"fmt"
	"go/types"
	"regexp"
	"strings"
)

// Type represents a Go data type.
type Type struct {
	typ types.Type
}

// String returns the unmodified name of the type. It's only useful for basic
// types and compositions thereof; for named types, use BareName or ShortName
// instead.
func (lt Type) String() string {
	return lt.typ.String()
}

// BareName returns the name of ultimate, underlying basic or named type, free
// from any pointer/slice/map decorators or package names. For maps, it returns
// the underlying key type instead of the value type.
func (lt Type) BareName() string {
	t := lt.typ
	for done := false; !done; {
		switch t.(type) {
		case *types.Basic, *types.Named:
			done = true
		case *types.Map:
			m := t.(*types.Map)
			t = m.Key()
		default:
			u := t.Underlying()
			if u != t && u != nil {
				t = u
			} else {
				done = true
			}
		}
	}
	name := t.String()
	if strings.Index(name, ".") > 0 {
		split := strings.Split(name, ".")
		return split[len(split)-1]
	}
	return name
}

// ShortName returns the type's name as usable from within a Go source file.
// You must pass a Resolver to handle package import names, as well as the
// package name declared in the source file in which this type name will appear.
func (lt Type) ShortName(local string, r Resolver) string {
	_, b := lt.typ.(*types.Basic)
	it, i := lt.typ.(*types.Interface)
	_, n := lt.typ.(*types.Named)
	slt, sl := lt.typ.(*types.Slice)
	mpt, mp := lt.typ.(*types.Map)
	ptt, pt := lt.typ.(*types.Pointer)

	if b {
		return lt.typ.String()
	} else if i {
		return it.String()
	} else if n {
		return r.Resolve(local, lt.typ.String())
	} else if pt {
		elem := Type{typ: ptt.Elem()}.ShortName(local, r)
		return fmt.Sprintf("*%s", elem)
	} else if sl {
		elem := Type{typ: slt.Elem()}.ShortName(local, r)
		return fmt.Sprintf("[]%s", elem)
	} else if mp {
		key := Type{typ: mpt.Key()}.ShortName(local, r)
		elem := Type{typ: mpt.Elem()}.ShortName(local, r)
		return fmt.Sprintf("map[%s]%s", key, elem)
	}
	panic(fmt.Sprintf("unhandled ShortName for type %s (%T)", lt.typ.String(), lt.typ))
}

var zeroConverts = regexp.MustCompile("(byte|u?int|float|rune)[0-9]*")

func (lt Type) ZeroValue(local string, r Resolver) string {
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

func (lt Type) zeroBasic(basic *types.Basic) string {
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
