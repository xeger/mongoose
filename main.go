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
	rend := gen.StubRenderer{}
	writer := gen.StdoutWriter{}

	placed := map[*parse.Interface]string{}
	for _, intf := range pkg.Interfaces {
		placed[&intf] = placer.Place(pkg.Dir, &intf)
	}

	byFile := map[string][]*parse.Interface{}
	for intf, path := range placed {
		list := byFile[path]
		byFile[path] = append(list, intf)
	}

	for path, intfs := range byFile {
		source, err := rend.Render(pkg, intfs)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Code generation failure:", err)
			os.Exit(1)
		}
		writer.Write(path, source)
	}
}
