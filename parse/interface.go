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
	name string
	intf *types.Interface
}

func (li loaderInterface) Size() int {
	return li.intf.NumMethods()
}

func (li loaderInterface) Name() string {
		return li.name
}

func (li loaderInterface) EachMethod(cb func(Method)) {
	for i := 0; i < li.intf.NumMethods(); i++ {
		meth := loaderMethod{li.intf.Method(i)}
		cb(meth)
	}
}

func (li loaderInterface) String() string {
	return fmt.Sprintf("Interface(size=%d)", li.Size())
}
