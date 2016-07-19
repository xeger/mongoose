package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/xeger/mongoose/gen"
	"github.com/xeger/mongoose/parse"
)

var mockPackage = flag.String("mock", "gomuti", "framework: testify,...")
var mockOutput = flag.String("out", ".", "dir/subdir for mock files (- for stdout)")

func gomuti() bool {
	return strings.Index(*mockPackage, "gomuti") >= 0
}

func testify() bool {
	return strings.Index(*mockPackage, "testify") >= 0
}

func placer() gen.Placer {
	return &gen.PackagePlacer{FilePerInterface: testify()}
}

func writer() gen.Writer {
	if *mockOutput == "-" {
		return &gen.StdoutWriter{}
	}
	return &gen.FileWriter{}
}

func renderer() gen.Renderer {
	if gomuti() {
		return gen.NewGomutiRenderer()
	} else if testify() {
		return gen.NewTestifyRenderer()
	}
	panic("not implemented")
}

func main() {
	flag.Parse()

	placer := placer()
	writer := writer()

	pkg, err := parse.NewPackage(flag.Arg(0))
	if err != nil {
		fmt.Println("Parse failure:", err)
		os.Exit(1)
	}

	placed := map[string][]parse.Interface{}
	for _, intf := range pkg.Interfaces {
		place := placer.Place(pkg.Dir, &intf)
		placed[place] = append(placed[place], intf)
	}

	for path, intfs := range placed {
		rend := renderer()
		source, err := rend.Render(pkg, intfs)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Code generation failure:", err)
			os.Exit(1)
		}
		writer.Write(path, source)
	}
}
