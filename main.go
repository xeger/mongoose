package main

import (
	"fmt"
	"os"

	"github.com/xeger/mongoose/gen"
	"github.com/xeger/mongoose/parse"
)

func main() {
	pkg, err := parse.NewPackage(os.Args[1])
	if err != nil {
		fmt.Println("// Parse failure:")
		fmt.Println("//   ", err)
		os.Exit(1)
	}

	// Proof of concept: render stubs to stdout
	rend := gen.StubRenderer{}
	source, err := rend.Render(pkg, pkg.Interfaces)
	if err != nil {
		fmt.Println("// Code generation failure:")
		fmt.Println("//   ", err)
		os.Exit(1)
	}
	os.Stdout.WriteString(source)
}
