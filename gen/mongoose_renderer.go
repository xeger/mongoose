package gen

import (
	"text/template"

	"github.com/xeger/mongoose/parse"
)

const mongooseItem = `
{{$typename := .Interface.Name | printf "Mock%s" }}type {{$typename}} struct {
	mock.Mock
}

{{$locl := .Package.Name}}{{$res := .Resolver}}{{range .Interface.Methods}}
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
