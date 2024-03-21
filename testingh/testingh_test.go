package testingh

import (
	"bytes"
	"crypto/rand"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"testing"
)

// ExampleStringReader shows how to use StringReader with a io.LimitedReader.
func ExampleStringReader() {
	var sr io.Reader
	sr = &StringReader{[]byte("1234567890")}
	limit := int64(15)
	buf := new(bytes.Buffer)
	limitReader := io.LimitReader(sr, limit)
	_, err := buf.ReadFrom(limitReader)
	if err != nil {
		fmt.Printf("io.Copy error: %+v", err)
	}
	rdbuf := make([]byte, limit)
	buf.Read(rdbuf)
	fmt.Printf("%+v", string(rdbuf))

	// Output:
	// 123456789012345
}

func TestStringReader(t *testing.T) {
	var sr io.Reader
	limit := int64(1e6)
	sr = &StringReader{[]byte("1234567890")}
	buf := new(bytes.Buffer)
	limitReader := io.LimitReader(sr, limit)
	n, err := buf.ReadFrom(limitReader)
	if err != nil {
		t.Errorf("io.Copy error: %+v", err)
	}

	if n != limit || int64(buf.Len()) != limit {
		t.Errorf("output data is not correct size, got: %d, want: %d", buf.Len(), limit)
	}
}

// TestCreateTestFileInMultipleDirectories tests creating testFiles, using
// different io.Readers in a directory per []testFile.
func TestCreateTestFileInMultipleDirectories(t *testing.T) {
	// Create slices of slices, where the index +1 == number of files.
	var sr io.Reader
	sr = &StringReader{[]byte("1234567890")}
	testFiles := [][]TestFile{
		{{int64(1e1), "f1", "", int64(9), &rand.Reader}},
		{{int64(1e1), "f2", "", int64(10), &sr},
			{int64(1e1), "f3", "", int64(11), &rand.Reader}},
		{{int64(1e1), "f4", "", int64(150), &sr},
			{int64(1e6), "f5", "", int64(1e7), &sr},
			{int64(1e6), "f6", "", int64(1e7) + int64(1), &rand.Reader}},
	}

	for _, tfs := range testFiles {
		outputPath := t.TempDir()
		for j := range tfs {
			err := CreateTestFile(t, outputPath, &tfs[j])
			if err != nil {
				t.Errorf("Error calling CreateTestFile: %+v", err)
			}

			testFileValidation(t, tfs[j])
		}

		di, err := os.ReadDir(filepath.Dir(tfs[0].FilePath))
		if err != nil {
			t.Errorf("ReadDir error: %+v", err)
		}
		if len(di) != len(tfs) {
			t.Errorf("Wrong number of files created, want: %d, got, %d", len(tfs), len(di))
		}
	}
}

// testFileValidation verifies a testFile was created, is a file, and is the correct size.
func testFileValidation(t *testing.T, tf TestFile) {
	fi, err := os.Stat(tf.FilePath)
	if err != nil {
		t.Errorf("Error getting os.Stat: %+v", err)
	}
	if fi.IsDir() {
		t.Error("File type is directory.")
	}
	if fi.Size() != tf.FileSize {
		t.Errorf("File size incorrect, want: %d, got: %d", tf.FileSize, fi.Size())
	}
}
