package ziph

import (
	"bytes"
	"crypto/rand"
	"fmt"
	"io"
	"path/filepath"
	"testing"

	"github.com/paulfdunn/go-helper/cryptoh"
	"github.com/paulfdunn/go-helper/testingh"
)

// TestZipUnzipShaCompare tests a round trip operation of creating a files, zipping,
// checking GetZipStats, unzipping, and comparing the checksum of the input and unzipped files.
func TestZipUnzipShaCompare(t *testing.T) {
	// Create test input files.
	testFileDir := t.TempDir()
	tfRand := testingh.TestFile{BufferSize: int64(1e1), FileName: "test_binary", FileSize: int64(1e6), Reader: &rand.Reader}
	err := testingh.CreateTestFile(t, testFileDir, &tfRand)
	if err != nil {
		t.Errorf("Error creating binary file: %+v", err)
	}
	var sr io.Reader
	sr = &testingh.StringReader{Data: []byte("1234567890")}
	tfStr := testingh.TestFile{BufferSize: int64(1e1), FileName: "test_string", FileSize: int64(1e6), Reader: &sr}
	err = testingh.CreateTestFile(t, testFileDir, &tfStr)
	if err != nil {
		t.Errorf("Error creating string file: %+v", err)
	}
	testFilePaths := []string{tfRand.FilePath, tfStr.FilePath}

	// Create the zip file and wait for zip completion
	zipDir := t.TempDir()
	fmt.Printf("zipDir: %s\n", zipDir)
	zipFilePath := filepath.Join(zipDir, "test_asynczip.zip")
	done, processedPaths, errs := AsyncZip(zipFilePath, testFilePaths)
	var dn bool
	var pathCount, errCount int
	dn = <-done
	fmt.Printf("done: %t\n", dn)
	noMessage := false
	for {
		select {
		case pp := <-processedPaths:
			pathCount++
			fmt.Printf("processed path: %s\n", pp)
		case err := <-errs:
			errCount++
			fmt.Printf("error: %v\n", err)
		default:
			noMessage = true
		}

		if noMessage {
			break
		}
	}

	// Test gitZipStats
	n, err := GetZipStats(zipFilePath)
	if n.FileCount != len(testFilePaths) || err != nil {
		t.Errorf("gitZipStats issue, n: %d, err: +%v", n, err)
	}

	// Test AsyncZip outputs.
	if pathCount != len(testFilePaths) || errCount != 0 {
		t.Errorf("AsyncZip pathCount: %d, errCount: %d", pathCount, errCount)
	}

	// AsyncUnzip the files into a new TempDir
	unzipDir := t.TempDir()
	done, processedPaths, errs = AsyncUnzip(zipFilePath, unzipDir, 0755)
	pathCount = 0
	errCount = 0
	dn = <-done
	fmt.Printf("done: %t\n", dn)
	noMessage = false
	for {
		select {
		case pp := <-processedPaths:
			pathCount++
			fmt.Printf("processed path: %s\n", pp)
		case err := <-errs:
			errCount++
			fmt.Printf("error: %v\n", err)
		default:
			noMessage = true
		}

		if noMessage {
			break
		}
	}

	// Test AsyncUnzip outputs.
	if pathCount != len(testFilePaths) || errCount != 0 {
		t.Errorf("AsyncUnzip pathCount: %d, errCount: %d", pathCount, errCount)
	}

	// Compare the hashes of the input files to the output files.
	for _, tp := range testFilePaths {
		testInputHash, err := cryptoh.Sha256FileHash(tp)
		if err != nil {
			t.Errorf("getting hash, error: %+v", err)
		}
		outputFileHash, err := cryptoh.Sha256FileHash(filepath.Join(unzipDir, tp))
		if err != nil {
			t.Errorf("getting hash, error: %+v", err)
		}
		if !bytes.Equal(testInputHash, outputFileHash) {
			t.Error("input and output hashes are not equal.")
		}
	}
}
