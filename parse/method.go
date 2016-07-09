package parse

import (
	"go/types"
)

type Method struct {
	Name    string
	Params  map[string]Type
	Results []Type
}

func (lm *Method) Arity() int {
	return len(lm.Params)
}

func (lm *Method) Len() int {
	return len(lm.Results)
}

func (lm *Method) finalize(meth *types.Func) {
	sig, ok := meth.Type().(*types.Signature)
	if !ok {
		panic("what the heck?")
	}
	if sig.Variadic() {
		panic("can't handle variadic interface methods!!!")
	}

	lm.Name = meth.Name()
	params := sig.Params()
	namer := make(namer)

	lm.Params = map[string]Type{}
	pos := 0
	for i := 0; i < params.Len(); i++ {
		p := params.At(i)
		name := p.Name()
		typ := Type{p.Type()}
		if name == "" {
			name = namer.Name(pos, typ)
		}
		lm.Params[name] = typ
		pos++
	}
	results := sig.Results()
	lm.Results = make([]Type, results.Len())
	for i := 0; i < results.Len(); i++ {
		p := results.At(i)
		lm.Results[i] = Type{p.Type()}
	}
}
