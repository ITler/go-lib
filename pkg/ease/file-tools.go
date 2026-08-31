package ease

import (
	"fmt"
	"io"
	"os"
)

// NewFileOrStdinReader creates a reader for a file or stdin
func NewFileOrStdinReader(filename string) (io.Reader, func() error, error) {
	if filename == "-" {
		return os.Stdin, nil, nil
	}
	file, err := os.Open(filename)
	if err != nil {
		return nil, nil, fmt.Errorf("open file readonly: %s - %w", filename, err)
	}
	return file, file.Close, nil
}

// NewFileOrStdoutWriter creates a writer for a file or stdout
func NewFileOrStdoutWriter(filename string) (io.Writer, func() error, error) {
	if filename == "-" {
		return os.Stdout, nil, nil
	}
	file, err := os.Create(filename)
	if err != nil {
		return nil, nil, fmt.Errorf("create file: %s - %w", filename, err)
	}
	return file, file.Close, nil
}
