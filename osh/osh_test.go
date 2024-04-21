package osh

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"testing"
	"time"
)

func ExampleDirIsEmpty() {
	u, _ := user.Current()
	b, _ := DirIsEmpty(u.HomeDir)
	fmt.Printf("User dir is empty? %+v\n", b)

	tmpDir, _ := os.MkdirTemp("", "")
	b, _ = DirIsEmpty(tmpDir)
	fmt.Printf("Temp dir is empty? %+v\n", b)
	os.Remove(tmpDir)
	// Output:
	// User dir is empty? false
	// Temp dir is empty? true
}

func TestFileModifiedFilter(t *testing.T) {
	modifiedSeconds := 30
	tmpDir := t.TempDir()
	oldFile := filepath.Join(tmpDir, "oldfile")
	_, err := os.Create(oldFile)
	if err != nil {
		t.Errorf("creating file, error: %+v", err)
	}
	os.Chtimes(oldFile, time.Time{}, time.Now().Add(-time.Duration(modifiedSeconds+1)*time.Second))
	newFile := filepath.Join(tmpDir, "newfile")
	_, err = os.Create(newFile)
	if err != nil {
		t.Errorf("creating file, error: %+v", err)
	}
	filteredFiles, err := FileModifiedFilter([]string{oldFile, newFile}, modifiedSeconds)
	if len(filteredFiles) != 1 || filteredFiles[0] != oldFile {
		t.Error("filtering did not work.")
	}
}

func TestRemoveAllFiles(t *testing.T) {
	tempDir := t.TempDir()
	testFiles := []string{"killme01.txt", "killme02.txt", "killme03.txt"}
	for _, v := range testFiles {
		_, err := os.Create(filepath.Join(tempDir, v))
		if err != nil {
			t.Errorf("error creating file, error:%v", err)
			return
		}
	}

	if err := RemoveAllFiles(filepath.Join(tempDir, "killme01*")); err != nil {
		t.Errorf("RemoveAllFiles error: %+v", err)
	}
	files, err := filepath.Glob(filepath.Join(tempDir, "killme*"))
	if err != nil {
		t.Errorf("getting file list, error: %v", err)
		return
	}
	if len(files) != 2 {
		t.Errorf("wrong number of files, len=%d", len(files))
		return
	}

	if err := RemoveAllFiles(filepath.Join(tempDir, "killme0*")); err != nil {
		t.Errorf("RemoveAllFiles error: %+v", err)
	}
	files, err = filepath.Glob(filepath.Join(tempDir, "killme*"))
	if err != nil {
		t.Errorf("getting file list, error: %v", err)
		return
	}
	if len(files) != 0 {
		t.Errorf("wrong number of files, len=%d", len(files))
		return
	}

}
