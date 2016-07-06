package main

import (
	"bytes"
	"fmt"
	"os"

	"github.com/xeger/mongoose/parse"
)

func main() {
	p, err := parse.NewPackage(os.Args[1])
	if err != nil {
		fmt.Println("//", err)
		os.Exit(1)
	}
	fmt.Println("// success", "pkg=", p)
	p.EachInterface(func(intf parse.Interface) {
		fmt.Println("type", intf.Name(), "interface", "{")
		intf.EachMethod(func(meth parse.Method) {
			params := bytes.NewBufferString("(")
			meth.EachParam(func(name string, typ parse.Type) {
				if params.Len() > 1 {
					params.WriteString(",")
				}
				params.WriteString(name)
				params.WriteString(" ")
				params.WriteString(typ.Name())
			})
			params.WriteString(")")
			fmt.Printf("  %s%s\n", meth.Name(), params.String())
		})
		fmt.Println("}")
		fmt.Println()
	})
}
