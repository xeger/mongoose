package parse

import (
	"bytes"
	"go/types"
)

type Param struct {
	Name string
	Type Type
}

type Params struct {
	data     []Param
	variadic bool
}

func (p *Params) Tuple(local string, resolver Resolver) string {
	buf := bytes.NewBufferString("(")
	for i := 0; i < p.Len(); i++ {
		pd := p.data[i]
		if buf.Len() > 1 {
			buf.WriteString(",")
		}
		buf.WriteString(pd.Name)
		buf.WriteRune(' ')
		buf.WriteString(pd.Type.ShortName(local, resolver))
	}
	if p.variadic {
		vp := p.Variadic()
		ut := Type{typ: vp.Type.typ.(*types.Slice).Elem()}
		buf.WriteString(vp.Name)
		buf.WriteString(" ...")
		buf.WriteString(ut.ShortName(local, resolver))
	}
	buf.WriteString(")")
	return buf.String()
}

func (p *Params) NameList() string {
	buf := bytes.Buffer{}

	for i := 0; i < p.Len(); i++ {
		pd := p.data[i]
		if buf.Len() > 0 {
			buf.WriteString(",")
		}
		buf.WriteString(pd.Name)
	}
	if p.variadic {
		if buf.Len() > 0 {
			buf.WriteString(",")
		}
		v := p.Variadic()
		buf.WriteString(v.Name)
		buf.WriteString("...")
	}
	return buf.String()
}

func (p *Params) TypeList(local string, resolver Resolver) string {
	buf := bytes.Buffer{}

	for i := 0; i < p.Len(); i++ {
		pd := p.data[i]
		if buf.Len() > 0 {
			buf.WriteString(",")
		}
		buf.WriteString(pd.Type.ShortName(local, resolver))
	}
	if p.variadic {
		if buf.Len() > 0 {
			buf.WriteString(",")
		}
		vp := p.Variadic()
		ut := Type{typ: vp.Type.typ.(*types.Slice).Elem()}
		buf.WriteString("...")
		buf.WriteString(ut.ShortName(local, resolver))
	}
	return buf.String()
}

func (p *Params) Len() int {
	if p.data == nil {
		return 0
	} else if p.variadic && len(p.data) > 0 {
		return len(p.data) - 1
	}
	return len(p.data)
}

func (p *Params) At(i int) *Param {
	if p.data == nil {
		return nil
	} else if i < 0 || i >= len(p.data) {
		return nil
	} else {
		return &p.data[i]
	}
}

func (p *Params) Variadic() *Param {
	if p.data == nil || len(p.data) < 1 {
		return nil
	} else if p.variadic {
		return &p.data[len(p.data)-1]
	}
	return nil
}
