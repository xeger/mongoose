package parse

import (
	"bytes"
	"go/types"
)

// Param is a method parameter.
type Param struct {
	Name string
	Type Type
}

// Params is the set of parameters to a method.
type Params struct {
	data     []Param
	variadic bool
}

// Tuple is the formal parameters declaration surrounded by parentheses
// e.g. "(alice string, bob int)"
func (p *Params) Tuple(local string, resolver Resolver) string {
	buf := bytes.NewBufferString("(")
	for i := 0; i < p.Len(); i++ {
		pd := p.data[i]
		if buf.Len() > 1 {
			buf.WriteString(", ")
		}
		buf.WriteString(pd.Name)
		buf.WriteRune(' ')
		buf.WriteString(pd.Type.ShortName(local, resolver))
	}
	if p.variadic {
		if buf.Len() > 1 {
			buf.WriteString(", ")
		}
		vp := p.Variadic()
		ut := Type{typ: vp.Type.typ.(*types.Slice).Elem()}
		buf.WriteString(vp.Name)
		buf.WriteString(" ...")
		buf.WriteString(ut.ShortName(local, resolver))
	}
	buf.WriteString(")")
	return buf.String()
}

// NameList is a comma-separated list of parameter names, including variadic (without trailing dots)
func (p *Params) NameList() string {
	buf := bytes.Buffer{}

	for i := 0; i < p.Len(); i++ {
		pd := p.data[i]
		if buf.Len() > 0 {
			buf.WriteString(", ")
		}
		buf.WriteString(pd.Name)
	}
	if p.variadic {
		if buf.Len() > 0 {
			buf.WriteString(", ")
		}
		v := p.Variadic()
		buf.WriteString(v.Name)
		//buf.WriteString("...")
	}
	return buf.String()
}

// TypeList is a comma-separated list of parameter types, including variadic (without trailing dots)
func (p *Params) TypeList(local string, resolver Resolver) string {
	buf := bytes.Buffer{}

	for i := 0; i < p.Len(); i++ {
		pd := p.data[i]
		if buf.Len() > 0 {
			buf.WriteString(", ")
		}
		buf.WriteString(pd.Type.ShortName(local, resolver))
	}
	if p.variadic {
		if buf.Len() > 0 {
			buf.WriteString(", ")
		}
		vp := p.Variadic()
		ut := Type{typ: vp.Type.typ.(*types.Slice).Elem()}
		buf.WriteString("...")
		buf.WriteString(ut.ShortName(local, resolver))
	}
	return buf.String()
}

// Len is the number of method parameters, excluding variadic .
func (p *Params) Len() int {
	if p.data == nil {
		return 0
	} else if p.variadic && len(p.data) > 0 {
		return len(p.data) - 1
	}
	return len(p.data)
}

// At returns the parameter in a specified position i.
func (p *Params) At(i int) *Param {
	if p.data == nil {
		return nil
	} else if i < 0 || i >= len(p.data) {
		return nil
	} else {
		return &p.data[i]
	}
}

// Variadic returns the variadic parameter, or nil if the method is
// not variadic.
func (p *Params) Variadic() *Param {
	if p.data == nil || len(p.data) < 1 {
		return nil
	} else if p.variadic {
		return &p.data[len(p.data)-1]
	}
	return nil
}
