package main

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"

	"github.com/xeger/mongoose/parse"
)

func main() {
	r := parse.NewResolver()
	local := filepath.Base(os.Args[1])
	p, err := parse.NewPackage(os.Args[1])
	if err != nil {
		fmt.Println("//", err)
		os.Exit(1)
	}

	p.EachInterface(func(intf parse.Interface) {
		fmt.Printf("type %s interface {\n", intf.Name())
		idx := 0
		intf.EachMethod(func(meth parse.Method) {
			params := bytes.NewBufferString("(")
			meth.EachParam(func(name string, typ parse.Type) {
				if params.Len() > 1 {
					params.WriteString(",")
				}
				if name == "" {
					idx++
					name = fmt.Sprintf("p%d", idx)
				}
				params.WriteString(name)
				params.WriteString(" ")
				params.WriteString(typ.ShortName(local, r))
			})
			params.WriteString(")")

			result := bytes.Buffer{}
			zeroes := bytes.Buffer{}
			meth.EachResult(func(typ parse.Type) {
				if result.Len() > 1 {
					result.WriteString(",")
					zeroes.WriteString(",")
				}
				result.WriteString(typ.ShortName(local, r))
				zeroes.WriteString(typ.ZeroValue(local, r))
			})

			fmt.Printf("  %s%s %s // return (%s)\n", meth.Name(), params.String(), result.String(), zeroes.String())
		})
		fmt.Println("}")
		fmt.Println()
	})
}
