package gen

import (
	"text/template"

	"github.com/xeger/mongoose/parse"
)

const gomutiItem = `
{{$typename := .Interface.Name | printf "Mock%s" }}// {{$typename}} is a test double for {{.Interface.Name}}.
// Generated with github.com/xeger/mongoose; do not edit by hand.
type {{$typename}} struct {
	Stub bool
	Mock gomuti.Mock
	Spy  gomuti.Spy
}

// Lazy-initialize spy controller whenever a mock method is called.
func (m *{{$typename}}) init() {
	if m.Spy == nil {
		m.Spy = gomuti.Spy{}
	}
}

{{$locl := .Package.Name}}{{$res := .Resolver}}{{range .Interface.Methods}}
func (m *{{$typename}}) {{.Name}}{{.Params.Tuple $locl $res}}{{$rtuple := .Results.Tuple $locl $res}}{{if gt .Results.Len 0}} {{$rtuple}}{{end}} {
	m.init()
	m.Spy.Observe("{{.Name}}", {{.Params.NameList}})
	{{$pnames := .Params.NameList}}{{$ptypes := (.Params.TypeList $locl $res)}}ret := m.Mock.Call("{{.Name}}",{{.Params.NameList}})
	if ret == nil {
		if m.Stub {
			return{{if gt .Results.Len 0}} {{.Results.ZeroList $locl $res}}{{end}}
		}
		panic("{{$typename}}: unexpected call to {{.Name}}")
	}
	{{range $idx, $typ := .Results}}
	var r{{$idx}} {{$typ.ShortName $locl $res}}
	if ret[{{$idx}}] == nil {
		r{{$idx}} = {{$typ.ZeroValue $locl $res}}
	} else {
		r{{$idx}} = ret[{{$idx}}].({{$typ.ShortName $locl $res}})
	}
	{{end}}
	return {{.Results.NameList}}
}
{{end}}
`

// NewGomutiRenderer creates a code generator using github.com/xeger/gomuti.
// The mock type contains a gomuti.Mock and can be programmed using the
// gomuti.Allow() method.
func NewGomutiRenderer() Renderer {
	r := parse.NewResolver()
	r.Import("gomuti", "github.com/xeger/gomuti/types")
	tr := &templateRenderer{}
	tr.Resolver = r
	tr.Header = template.New("templateHeader")
	tr.Header.Parse(templateHeader)
	tr.Item = template.New("gomutiItem")
	tr.Item.Parse(gomutiItem)
	return tr
}
