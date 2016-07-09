package parse

import (
	"fmt"
	"go/types"
)

// Interface represents an interface
type Interface interface {
	Size() int
	Name() string
	EachMethod(func(Method))
}

type loaderInterface struct {
	name    string
	methods []Method
}

func (li *loaderInterface) Size() int {
	return len(li.methods)
}

func (li *loaderInterface) Name() string {
	return li.name
}

func (li *loaderInterface) EachMethod(cb func(Method)) {
	for _, meth := range li.methods {
		cb(meth)
	}
}

func (li *loaderInterface) String() string {
	return fmt.Sprintf("Interface(size=%d)", li.Size())
}

func (li *loaderInterface) finalize(name string, intf *types.Interface) {
	li.name = name
	li.methods = make([]Method, intf.NumMethods())
	for i := 0; i < len(li.methods); i++ {
		meth := &loaderMethod{}
		meth.finalize(intf.Method(i))
		li.methods[i] = meth
	}
}
