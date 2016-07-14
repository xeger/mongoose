package gen

import (
	"text/template"

	"github.com/xeger/mongoose/parse"
)

const mongooseItem = `
{{$typename := .Interface.Name | printf "Mock%s" }}type {{$typename}} struct {
	Mock mock.Mock
	Stub bool
}

{{$locl := .Package.Name}}{{$res := .Resolver}}{{range .Interface.Methods}}
func (m {{$typename}}) {{.Name}}{{.Params.Tuple $locl $res}}{{$rtuple := .Results.Tuple $locl $res}}{{if gt .Results.Len 0}} {{$rtuple}}{{end}} {
	{{$pnames := .Params.NameList}}{{$ptypes := (.Params.TypeList $locl $res)}}ret := mock.Ø(m.Mock,"{{.Name}}",{{.Params.NameList}})
	if ret == nil {
		if m.Stub {
			return{{if gt .Results.Len 0}} {{.Results.ZeroTuple $locl $res}}{{end}}
		}
		panic("mock: unexpected call to {{.Name}}")
	}
	{{range $idx, $typ := .Results}}
	var r{{$idx}} {{$typ.ShortName $locl $res}}
	{{if eq $typ.String "error"}}r{{$idx}} = ret.Error({{$idx}}){{else}}r{{$idx}} = ret[{{$idx}}].({{$typ.ShortName $locl $res}}){{end}}
	{{end}}
	return {{.Results.NameList}}
}
{{end}}
`

// NewRenderer creates a code generator using github.com/xeger/mongoose/mock.
// The mock type is derived from mock.Mock and can be programmed using the
// mock.Allow() method on instances.
func NewRenderer() Renderer {
	r := parse.NewResolver()
	r.Import("mock", "github.com/xeger/mongoose/mock")
	tr := &templateRenderer{}
	tr.Resolver = r
	tr.Header = template.New("mongooseHeader")
	tr.Header.Parse(templateHeader)
	tr.Item = template.New("mongooseItem")
	tr.Item.Parse(mongooseItem)
	return tr
}
