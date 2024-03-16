package osh

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// ByteSliceToIntSlice converts an byte slice to integer slice
func ByteSliceToIntSlice(bytes []byte) []int {
	out := make([]int, len(bytes))
	for i := range bytes {
		out[i] = int(bytes[i])
	}
	return out
}

// ByteSliceToString converts a byte slice to a hex string with bytesPerLine; no "0x" prefix.
func ByteSliceToString(in []byte, bytesPerLine int) (out string) {
	index := 0
	// Convert bytes to ints, needed for formatting later.
	inInts := make([]int, len(in))
	for i, v := range in {
		inInts[i] = int(v)
	}

	for {
		if index >= len(inInts) {
			break
		}

		// Print bytesPerLine, or a partial line if there are not enough bytes left.
		end := index + bytesPerLine
		if end > len(inInts) {
			end = len(inInts)
		}

		// Ints format nicely with this; space separated.
		s := fmt.Sprintf("%02x", inInts[index:end])
		out += fmt.Sprintf("%s", s[1:len(s)-1])
		out += "\n"
		index += bytesPerLine
	}
	return out
}

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
