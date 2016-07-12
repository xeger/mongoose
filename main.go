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
	placer := gen.PackagePlacer{}
	rend := gen.NewTestifyRenderer()
	writer := gen.StdoutWriter{}

	placed := map[string][]parse.Interface{}
	for _, intf := range pkg.Interfaces {
		place := placer.Place(pkg.Dir, &intf)
		placed[place] = append(placed[place], intf)
	}

	for path, intfs := range placed {
		source, err := rend.Render(pkg, intfs)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Code generation failure:", err)
			os.Exit(1)
		}
		writer.Write(path, source)
	}
}
