package parse

import (
	"fmt"
	"go/types"
	"path/filepath"
	"regexp"
	"strings"
)

// Type represents a Go data type.
type Type struct {
	typ types.Type
}

func (lt Type) Name() string {
	return lt.typ.String()
}

func (lt Type) BareName() string {
	pkgWithName := filepath.Base(lt.Name())
	parts := strings.SplitN(pkgWithName, ".", 2)
	return parts[len(parts)-1]
}

func (lt Type) ShortName(local string, r Resolver) string {
	_, b := lt.typ.(*types.Basic)
	slt, sl := lt.typ.(*types.Slice)
	mpt, mp := lt.typ.(*types.Map)
	ptt, pt := lt.typ.(*types.Pointer)
	_, n := lt.typ.(*types.Named)

	typ := lt.typ.String()

	if b {
		return typ
	} else if n {
		return r.Resolve(local, typ)
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
	panic(fmt.Sprintf("unhandled ShortName for type %s", typ))
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
