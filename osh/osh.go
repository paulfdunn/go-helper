// Package osh provides helper functions for the os package.
package osh

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"
)

// DirIsEmpty returns true if the directory exists and is empty.
func DirIsEmpty(path string) (bool, error) {
	f, err := os.Open(path)
	if err != nil {
		// Return false if the directory does not exist.
		return false, err
	}
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Printf("f.Close() error:%+v\n", err)
		}
	}()

	// Readdirnames does NOT return "." and ".."; so a single file indicates the dir
	// is not empty.
	_, err = f.Readdirnames(1)
	if err == io.EOF {
		return true, nil
	}

	return false, err
}

// FileModifiedAfterFilter will filter a slice of input files and remove files that
// were not modified within the last modifiedSeconds.
func FileModifiedAfterFilter(filepaths []string, modifiedSeconds int) ([]string, error) {
	output := make([]string, 0, len(filepaths))
	for _, fp := range filepaths {
		fi, err := os.Stat(fp)
		if err != nil {
			return nil, err
		}
		if fi.ModTime().After(time.Now().Add(-time.Duration(modifiedSeconds) * time.Second)) {
			output = append(output, fp)
		}
	}
	return output, nil
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
