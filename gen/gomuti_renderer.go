package gen

import (
	"text/template"
)

const gomutiItem = `
{{$typename := .Interface.Name | printf "Mock%s" }}// {{$typename}} is a test double for {{.Interface.Name}}.
// Generated with github.com/xeger/mongoose; do not edit by hand.
type {{$typename}} struct {
	Stub bool
	Mock gomuti.Mock
	Spy  gomuti.Spy
}

// init lazy-initializes a spy controller whenever a mock method is called.
func (m *{{$typename}}) init() {
	if m.Spy == nil {
		m.Spy = gomuti.Spy{}
	}
}

{{$locl := .Package.Name}}{{$res := .Resolver}}
{{range .Interface.Methods}}
{{$mname := .Name}}
{{$rtuple := .Results.Tuple $locl $res}}
{{$pnames := .Params.NameList}}
{{$ptypes := (.Params.TypeList $locl $res)}}
// {{$mname}} is a mock interface method.
func (m *{{$typename}}) {{$mname}}{{.Params.Tuple $locl $res}}{{if gt .Results.Len 0}} {{$rtuple}}{{end}} {
	m.init()
	m.Spy.Observe("{{$mname}}", {{.Params.NameList}})
	ret := m.Mock.Call("{{$mname}}", {{.Params.NameList}})
	if ret == nil {
		if m.Stub {
			return{{if gt .Results.Len 0}} {{.Results.ZeroList $locl $res}}{{end}}
		}
		panic("{{$typename}}: unexpected call to {{$mname}}")
	}

	{{$rlen := .Results | len}}
	if len(ret) != {{$rlen}} {
		panic(fmt.Sprintf("{{$typename}}.{{$mname}}: return value mismatch; expected %d, got %d: %v", {{$rlen}}, len(ret), ret))
	}
	{{range $idx, $typ := .Results}}
		var r{{$idx}} {{$typ.ShortName $locl $res}}
		if ret[{{$idx}}] == nil {
			r{{$idx}} = {{$typ.ZeroValue $locl $res}}
		} else {
			r{{$idx}}t, ok := ret[{{$idx}}].({{$typ.ShortName $locl $res}})
			if !ok {
				panic(fmt.Sprintf("{{$typename}}.{{$mname}}: return type mismatch; expected %T, got %T", r{{$idx}}t, ret[{{$idx}}]))
			}
			r{{$idx}} = r{{$idx}}t
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
	tr := newTemplateRenderer()
	tr.Resolver.Import("fmt", "fmt")
	tr.Resolver.Import("gomuti", "github.com/xeger/gomuti/types")
	tr.Item = template.New("gomutiItem")
	tr.Item.Parse(gomutiItem)
	return tr
}
