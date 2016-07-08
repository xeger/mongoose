package parse

import (
	"go/types"
)

// Method represents a method definition of an interface or other type.
type Method interface {
	Size() int
	Name() string
	EachParam(func(string, Type))
	EachResult(func(Type))
}

type loaderMethod struct {
	funk *types.Func
}

func (lm loaderMethod) Size() int {
	sig := lm.sig()
	return sig.Params().Len() + sig.Results().Len()
}

func (lm loaderMethod) Name() string {
	return lm.funk.Name()
}

func (lm loaderMethod) sig() *types.Signature {
	sig, ok := lm.funk.Type().(*types.Signature)
	if !ok {
		panic("what the heck?")
	}
	if sig.Variadic() {
		panic("can't handle variadic interface methods!!!")
	}
	return sig
}

func (lm loaderMethod) EachParam(cb func(string, Type)) {
	params := lm.sig().Params()
	for i := 0; i < params.Len(); i++ {
		p := params.At(i)
		cb(p.Name(), loaderType{p.Type()})
	}
}

func (lm loaderMethod) EachResult(cb func(Type)) {
	results := lm.sig().Results()
	for i := 0; i < results.Len(); i++ {
		p := results.At(i)
		cb(loaderType{p.Type()})
	}
}
