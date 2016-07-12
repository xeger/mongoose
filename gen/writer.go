package gen

import (
	"fmt"
	"os"
	"path/filepath"
)

type Writer interface {
	Write(path, source string)
}

type StdoutWriter struct{}

func (StdoutWriter) Write(path, source string) {
	fmt.Println("////////////////////////////////////////////////////////////")
	fmt.Println("//", filepath.Base(path))
	fmt.Println("////////////////////////////////////////////////////////////")
	fmt.Println()
	os.Stdout.WriteString(source)
}
