package gen

import (
	"fmt"
	"os"
	"path/filepath"
)

// Writer can write source to specified absolute file paths.
type Writer interface {
	Write(path, source string) error
}

// StdoutWriter prints file contents to the screen.
type StdoutWriter struct{}

// FileWriter creates normal files on disk.
type FileWriter struct{}

// Write prints source to the screen.
func (StdoutWriter) Write(path, source string) error {
	fmt.Println("////////////////////////////////////////////////////////////")
	fmt.Println("//", filepath.Base(path))
	fmt.Println("////////////////////////////////////////////////////////////")
	fmt.Println(source)
	return nil
}

// Write creates/truncates a file at path and writes source to it.
func (FileWriter) Write(path, source string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.WriteString(source)
	return err
}
