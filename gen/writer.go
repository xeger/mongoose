package gen

import (
	"fmt"
	"os"
	"path/filepath"
)

type Writer interface {
	Write(path, source string) error
}

type StdoutWriter struct{}

type FileWriter struct{}

func (StdoutWriter) Write(path, source string) error {
	fmt.Println("////////////////////////////////////////////////////////////")
	fmt.Println("//", filepath.Base(path))
	fmt.Println("////////////////////////////////////////////////////////////")
	fmt.Println(source)
	return nil
}

func (FileWriter) Write(path, source string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.WriteString(source)
	return err
}
