// Package file provides functions to deal with FS
package file

import (
	"io"
	"os"
	"path/filepath"
)

// FileCreation makes sure a file at a specific path exist with the provided content.
func FileCreation(filePath string, content io.Reader) error {
	if err := os.MkdirAll(filepath.Dir(filePath), 0700); err != nil {
		return err
	}
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = io.Copy(file, content)
	return err
}
