package parse

import (
	"go/types"
)

type Method struct {
	Name    string
	Params  Params
	Results Results
}

func (meth *Method) Arity() int {
	return meth.Params.Len()
}

func (meth *Method) Len() int {
	return meth.Results.Len()
}

func (meth *Method) String() string {
	return meth.Name
}

func (meth *Method) finalize(actual *types.Func) {
	sig, ok := actual.Type().(*types.Signature)
	if !ok {
		panic("what the heck?")
	}

	meth.Name = actual.Name()
	params := sig.Params()
	namer := make(namer)

	meth.Params = Params{data: make([]Param, 0, params.Len()), variadic: sig.Variadic()}
	pos := 0
	for i := 0; i < params.Len(); i++ {
		p := params.At(i)
		name := p.Name()
		typ := Type{p.Type()}
		if name == "" {
			name = namer.Name(pos, typ)
		}
		meth.Params.data = append(meth.Params.data, Param{Name: name, Type: typ})
		pos++
	}
	results := sig.Results()
	meth.Results = make([]Type, results.Len())
	for i := 0; i < results.Len(); i++ {
		p := results.At(i)
		meth.Results[i] = Type{p.Type()}
	}
}
