package osh

import (
	"io"
	"os"
	"path/filepath"
)

// DirIsEmpty returns true if the directory exists and is empty.
func DirIsEmpty(path string) (bool, error) {
	f, err := os.Open(path)
	defer f.Close()
	if err != nil {
		// Return false if the directory does not exist.
		return false, err
	}

	// Readdirnames does NOT return "." and ".."; so a single file indicates the dir
	// is not empty.
	_, err = f.Readdirnames(1)
	if err == io.EOF {
		return true, nil
	}

	return false, err
}

// RemoveAllFiles will remove all files meeting the Glob pattern.
func RemoveAllFiles(pattern string) error {
	files, err := filepath.Glob(pattern)
	if err != nil {
		return err
	}
	for _, f := range files {
		if err := os.Remove(f); err != nil {
			return err
		}
	}
	return nil
}
