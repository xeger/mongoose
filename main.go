package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/xeger/mongoose/gen"
	"github.com/xeger/mongoose/parse"
)

var mockPackage = flag.String("mock", "gomuti", "framework: testify,...")
var recurse = flag.Bool("r", false, "recurse into subdirectories")
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

type outcome struct {
	dir        string
	phase      string
	interfaces int
	err        error
}

var nonPackage = regexp.MustCompile("^(.bzr|.git|.hg|.svn|vendor)$")

func findPackages(basedir string) ([]string, error) {
	if *recurse {
		packages := make([]string, 0, 10)

		werr := filepath.Walk(basedir, func(path string, fi os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if fi.IsDir() {
				if nonPackage.MatchString(filepath.Base(path)) {
					return filepath.SkipDir
				}
				packages = append(packages, path)
			}
			return nil
		})

		if werr != nil {
			return nil, werr
		}

		return packages, nil
	}

	return []string{basedir}, nil
}

func doPackage(dir string, oc chan outcome) {
	placer := placer()
	writer := writer()

	pkg, err := parse.NewPackage(dir)
	if err != nil {
		if strings.Index(err.Error(), "no buildable Go source files") == 0 {
			// not really an error...
			oc <- outcome{dir, "Parse", 0, nil}
			return
		}
		oc <- outcome{dir, "Parse", 0, err}
		return
	}

	placed := map[string][]parse.Interface{}
	for _, intf := range pkg.Interfaces {
		place := placer.Place(pkg.Dir, &intf)
		placed[place] = append(placed[place], intf)
	}

	total := 0
	for path, intfs := range placed {
		rend := renderer()
		source, err := rend.Render(pkg, intfs)
		if err != nil {
			oc <- outcome{dir, "Code generation", total, err}
			return
		}
		err = writer.Write(path, source)
		if err != nil {
			oc <- outcome{dir, "Code generation", total, err}
			return
		}
		total += len(intfs)
	}

	oc <- outcome{dir, "Code generation", total, nil}
}

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [flags] <dir>\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Generates a mock for every golang interface defined in <dir>\n\n")
		fmt.Fprintf(os.Stderr, "Flags:\n")
		flag.PrintDefaults()
	}

	flag.Parse()

	dirs, err := findPackages(flag.Arg(0))
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", err)
		os.Exit(1)
	}

	outcomes := make(chan outcome, 3)

	for _, dir := range dirs {
		go doPackage(dir, outcomes)
	}

	failed := 0
	interfaces := 0
	for i := 0; i < len(dirs); i++ {
		oc := <-outcomes
		interfaces += oc.interfaces
		if oc.err != nil {
			fmt.Fprintf(os.Stderr, "%s: %s: %s\n", oc.dir, oc.phase, oc.err)
			failed++
		} else if oc.interfaces > 0 {
			fmt.Fprintf(os.Stderr, "%s: mocked %d interfaces\n", oc.dir, oc.interfaces)
		}
	}

	if failed > 0 {
		fmt.Fprintf(os.Stderr, "\nmongoose: encountered errors in %d packages!\n", failed)
		os.Exit(failed)
	} else if interfaces == 0 {
		fmt.Fprintf(os.Stderr, "\nmongoose: did not find any interfaces to mock!\n")
		os.Exit(-1)
	}
}
