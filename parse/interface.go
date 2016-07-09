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

func (li *Interface) Len() int {
	return len(li.Methods)
}

func (li *Interface) String() string {
	return fmt.Sprintf("Interface(size=%d)", li.Len())
}

func (li *Interface) finalize(name string, intf *types.Interface) {
	li.Name = name
	li.Methods = make([]Method, intf.NumMethods())
	for i := 0; i < len(li.Methods); i++ {
		meth := Method{}
		meth.finalize(intf.Method(i))
		li.Methods[i] = meth
	}
}
