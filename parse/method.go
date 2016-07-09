package parse

import (
	"go/types"
)

// Method represents a method definition of an interface or other type.
type Method interface {
	Arity() int
	Len() int
	Name() string
	EachParam(func(string, Type))
	EachResult(func(Type))
}

type loaderMethod struct {
	name    string
	params  map[string]Type
	results []Type
}

func (lm *loaderMethod) Arity() int {
	return len(lm.params)
}

func (lm *loaderMethod) Len() int {
	return len(lm.results)
}

func (lm *loaderMethod) Name() string {
	return lm.name
}

func (lm *loaderMethod) EachParam(cb func(string, Type)) {
	namer := make(namer)
	pos := 0
	for name, typ := range lm.params {
		if name == "" {
			name = namer.Name(pos, typ)
		}
		cb(name, typ)
		pos++
	}
}

func (lm *loaderMethod) EachResult(cb func(Type)) {
	for _, typ := range lm.results {
		cb(typ)
	}
}

func (lm *loaderMethod) finalize(meth *types.Func) {
	sig, ok := meth.Type().(*types.Signature)
	if !ok {
		panic("what the heck?")
	}
	if sig.Variadic() {
		panic("can't handle variadic interface methods!!!")
	}

	lm.name = meth.Name()
	params := sig.Params()
	lm.params = map[string]Type{}
	for i := 0; i < params.Len(); i++ {
		p := params.At(i)
		lm.params[p.Name()] = loaderType{p.Type()}
	}
	results := sig.Results()
	lm.results = make([]Type, results.Len())
	for i := 0; i < results.Len(); i++ {
		p := results.At(i)
		lm.results[i] = loaderType{p.Type()}
	}
}
