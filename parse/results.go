package parse

import (
	"bytes"
	"fmt"
)

type Results []Type

func (r Results) Tuple(local string, resolver Resolver) string {
	buf := bytes.Buffer{}
	if r.Len() > 1 {
		buf.WriteString("(")
	}
	for _, typ := range r {
		if buf.Len() > 1 {
			buf.WriteRune(',')
		}
		buf.WriteString(typ.ShortName(local, resolver))
	}
	if r.Len() > 1 {
		buf.WriteString(")")
	}

	return buf.String()
}

func (r Results) ZeroTuple(local string, resolver Resolver) string {
	buf := bytes.Buffer{}
	if r.Len() > 1 {
		buf.WriteString("(")
	}
	for _, typ := range r {
		if buf.Len() > 1 {
			buf.WriteRune(',')
		}
		buf.WriteString(typ.ZeroValue(local, resolver))
	}
	if r.Len() > 1 {
		buf.WriteString(")")
	}

	return buf.String()
}

func (r Results) NameList() string {
	buf := bytes.Buffer{}
	for i := range r {
		if buf.Len() > 0 {
			buf.WriteRune(',')
		}
		buf.WriteString(fmt.Sprintf("r%d", i))
	}
	return buf.String()
}

func (r Results) Len() int {
	if r == nil {
		return 0
	}
	return len(r)
}

func (r Results) At(i int) *Type {
	if r == nil {
		return nil
	} else if i < 0 || i >= len(r) {
		return nil
	}
	return &r[i]
}
