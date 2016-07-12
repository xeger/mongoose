package gen

import (
	"fmt"
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
	fmt.Println(source)
}
