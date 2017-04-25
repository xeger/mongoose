// Package main implements the Mongoose CLI tool. For usage information, see the README: https://github.com/xeger/mongoose
package main

import (
	"flag"
	"fmt"
	"go/format"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/xeger/mongoose/gen"
	"github.com/xeger/mongoose/parse"
)

var mockPackage = flag.String("mock", "gomuti", "mocking framework: gomuti|testify")
var recurse = flag.Bool("r", false, "recurse into subdirectories")
var mockOutput = flag.String("out", ".", "package-relative subdir for mock files (- for stdout)")
var name = flag.String("name", "", "name or matching regular expression of interface to generate mock for")

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

func findPackages(basedirs []string, recurse bool) ([]string, error) {
	packages := make([]string, 0, 10)

	for _, basedir := range basedirs {
		abs, err := filepath.Abs(basedir)
		if err != nil {
			return nil, err
		}

		if recurse {
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
				return nil, err
			}
		} else {
			packages = append(packages, abs)
		}
	}
	return packages, nil
}

func doPackage(dir string, name *regexp.Regexp, oc chan outcome) {
	placer := placer()
	writer := writer()

	pkg, err := parse.NewPackage(dir)
	if strings.Index(pkg.Name(), "fuckme") >= 0 {
		panic(pkg.Name())
	}
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
		if name != nil && !name.MatchString(intf.Name) {
			continue
		}
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

		formatted, err := format.Source([]byte(source))
		if err != nil {
			oc <- outcome{dir, "Code generation", total, err}
			return
		}
		source = string(formatted)

		err = writer.Write(path, source)
		if err != nil {
			oc <- outcome{dir, "Code generation", total, err}
			return
		}

		total += len(intfs)
	}

	oc <- outcome{dir, "Code generation", total, nil}
}

// Regexp that matches a valid golang identifier (i.e. type name); used to
// determine whether we should treat the -name flag as a single name, or as
// a regular expression that matches a range of names.
var identifierPat = regexp.MustCompile(`[\p{L}_][\p{L}\p{N}_]*`)

func parseName(name string) (*regexp.Regexp, error) {
	if name != "" && identifierPat.MatchString(name) {
		return regexp.Compile(fmt.Sprintf("^%s$", regexp.QuoteMeta(name)))
	} else if name != "" {
		return regexp.Compile(name)
	} else {
		return nil, nil
	}
}

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [flags] <dir> [dir,dir,...]\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Generates a mock for every golang interface defined in any named dir\n\n")
		fmt.Fprintf(os.Stderr, "Flags:\n")
		flag.PrintDefaults()
	}

	flag.Parse()

	var dirs []string

        if len(flag.Args()) == 0 {  // if gopackage := os.Getenv("GOPACKAGE"); gopackage != "" {
		// go-generate mode: PWD is the one and only package
		pwd, err := os.Getwd()
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %s\n", err)
		}
		dirs = []string{pwd}
	} else {
		// normal mode: packages are specified as CLI args, possibly influenced
		// by the -r flag
		found, err := findPackages(flag.Args(), *recurse)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %s\n", err)
			os.Exit(1)
		}
		dirs = found
	}

	outcomes := make(chan outcome, 3)

	namePat, err := parseName(*name)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", err)
	}

	for _, dir := range dirs {
		go doPackage(dir, namePat, outcomes)
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
