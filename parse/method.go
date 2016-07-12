package parse

import (
	"go/types"
)

type Method struct {
	Name    string
	Params  Params
	Results Results
}

func (lm *Method) Arity() int {
	return lm.Params.Len()
}

func (lm *Method) Len() int {
	return lm.Results.Len()
}

func (lm *Method) finalize(meth *types.Func) {
	sig, ok := meth.Type().(*types.Signature)
	if !ok {
		panic("what the heck?")
	}

	lm.Name = meth.Name()
	params := sig.Params()
	namer := make(namer)

	lm.Params = Params{data: make([]Param, 0, params.Len()), variadic: sig.Variadic()}
	pos := 0
	for i := 0; i < params.Len(); i++ {
		p := params.At(i)
		name := p.Name()
		typ := Type{p.Type()}
		if name == "" {
			name = namer.Name(pos, typ)
		}
		lm.Params.data = append(lm.Params.data, Param{Name: name, Type: typ})
		pos++
	}
	results := sig.Results()
	lm.Results = make([]Type, results.Len())
	for i := 0; i < results.Len(); i++ {
		p := results.At(i)
		lm.Results[i] = Type{p.Type()}
	}
}
