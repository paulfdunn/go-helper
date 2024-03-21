// Package testingh provides helper functions for the testing package.
package testingh

import (
	"bufio"
	"io"
	"os"
	"path/filepath"
	"testing"
)

// StringReader implements io.Reader and can be used to generate data for test files.
// Use it with a io.LimitedReader to generate a specific amount of data.
type StringReader struct {
	Data []byte
}

type TestFile struct {
	BufferSize int64
	FileName   string
	FilePath   string // filePath is an output of CreateTestFile
	FileSize   int64
	Reader     *io.Reader
}

func (sr *StringReader) Read(p []byte) (int, error) {
	copy(p, sr.Data)
	return min(len(p), len(sr.Data)), nil
}

// CreateTestFile is a function that will create a file, filled with data from the provided
// io.Reader. Call with outputPath created from a call to t.TempDir() to create a directory
// that is cleaned up automatically on test exit by the Testing package.
// To prevent filling up memory the data is read, from the testFile.Reader, into a local
// buffer, in BufferSize increments. The fully qualified filepath of the created file is returned
// in the testFile.FilePath.
// See the tests for example of how to create files with random binary data, or repeating strings.
func CreateTestFile(t *testing.T, outputPath string, testFile *TestFile) error {
	testFile.FilePath = filepath.Join(outputPath, testFile.FileName)
	limitReader := io.LimitReader(*testFile.Reader, testFile.FileSize)

	f, err := os.Create(testFile.FilePath)
	if err != nil {
		t.Errorf("file creation error: %v", err)
		return err
	}
	defer f.Close()

	w := bufio.NewWriter(f)
	buf := make([]byte, testFile.BufferSize)
	for {
		nBuf, err := limitReader.Read(buf)
		if err == io.EOF {
			break
		}
		// fmt.Printf("nBuf: %d, err: %+v, data: %s\n", nBuf, err, base64.StdEncoding.EncodeToString(buf))

		if _, err = w.Write(buf[:nBuf]); err != nil {
			t.Errorf("writing temp file: %+v", err)
		}
	}

	if err = w.Flush(); err != nil {
		t.Errorf("flushing file: %+v", err)
	}

	return nil
}
