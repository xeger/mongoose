package gen

import (
	"text/template"

	"github.com/xeger/mongoose/parse"
)

const testifyItem = `
{{$typename := .Interface.Name | printf "Mock%s" }}type {{$typename}} struct {
	mock.Mock
}

{{$locl := .Package.Name}}{{$res := .Resolver}}{{range .Interface.Methods}}
func (_m *{{$typename}}) {{.Name}}{{.Params.Tuple $locl $res}}{{$rtuple := .Results.Tuple $locl $res}}{{if gt .Results.Len 0}} {{$rtuple}}{{end}} {
	{{$pnames := .Params.NameList}}{{$ptypes := (.Params.TypeList $locl $res)}}{{if gt .Results.Len 0}}ret := {{end}}_m.Called({{.Params.NameList}})
	{{range $idx, $typ := .Results}}
	var r{{$idx}} {{$typ.ShortName $locl $res}}

	if r{{$idx}}f, ok := ret.Get({{$idx}}).(func({{$ptypes}}) {{$typ.ShortName $locl $res}}); ok {
			r{{$idx}} = r{{$idx}}f({{$pnames}})
	} else {
			{{if eq $typ.String "error"}}r{{$idx}} = ret.Error({{$idx}}){{else}}r{{$idx}} = ret.Get({{$idx}}).({{$typ.ShortName $locl $res}}){{end}}
	}{{end}}

	return {{.Results.NameList}}
}
{{end}}
`

// NewTestifyRenderer creates a code generator for github.com/stretchr/testify.
// The mock type embeds tesify/mock.Mock and can be programmed using the
// embedded methods.
func NewTestifyRenderer() Renderer {
	r := parse.NewResolver()
	r.Import("mock", "github.com/stretchr/testify/mock")
	tr := &templateRenderer{}
	tr.Resolver = r
	tr.Header = template.New("testifyHeader")
	tr.Header.Parse(templateHeader)
	tr.Item = template.New("testifyItem")
	tr.Item.Parse(testifyItem)
	return tr
}
