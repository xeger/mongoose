package parse

import (
	"bytes"
	"fmt"
)

// Results are the type(s) returned by an interface method.
type Results []Type

// Tuple returns a comma-separated list of types enclosed in parentheses.
// It is suitable for printing as the return values of a method declaration.
func (r Results) Tuple(local string, resolver Resolver) string {
	buf := bytes.Buffer{}
	if r.Len() > 1 {
		buf.WriteString("(")
	}
	for _, typ := range r {
		if buf.Len() > 1 {
			buf.WriteString(", ")
		}
		buf.WriteString(typ.ShortName(local, resolver))
	}
	if r.Len() > 1 {
		buf.WriteString(")")
	}

	return buf.String()
}

// ZeroList returns a comma-separated list of zero values that are suitable
// as a return value or for other purposes.
func (r Results) ZeroList(local string, resolver Resolver) string {
	buf := bytes.Buffer{}
	for _, typ := range r {
		if buf.Len() > 1 {
			buf.WriteString(", ")
		}
		buf.WriteString(typ.ZeroValue(local, resolver))
	}

	return buf.String()
}

// NameList returns a comma-separated list of parameter names.
func (r Results) NameList() string {
	buf := bytes.Buffer{}
	for i := range r {
		if buf.Len() > 0 {
			buf.WriteString(", ")
		}
		buf.WriteString(fmt.Sprintf("r%d", i))
	}
	return buf.String()
}

// Len returns the number of values returned by the method.
func (r Results) Len() int {
	if r == nil {
		return 0
	}
	return len(r)
}

// At returns the type of the return value at index i.
func (r Results) At(i int) *Type {
	if r == nil {
		return nil
	} else if i < 0 || i >= len(r) {
		return nil
	}
	return &r[i]
}
