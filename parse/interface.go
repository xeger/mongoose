package parse

import (
	"fmt"
	"go/types"
)

// Interface represents an interface defined in a Package.
type Interface struct {
	Name    string
	Methods []Method
}

func (intf *Interface) Len() int {
	return len(intf.Methods)
}

func (intf *Interface) String() string {
	return fmt.Sprintf("type %s interface {%v}", intf.Name, intf.Methods)
}

func (intf *Interface) finalize(name string, actual *types.Interface) {
	intf.Name = name
	intf.Methods = make([]Method, actual.NumMethods())
	for i := 0; i < len(intf.Methods); i++ {
		meth := Method{}
		meth.finalize(actual.Method(i))
		intf.Methods[i] = meth
	}
}
