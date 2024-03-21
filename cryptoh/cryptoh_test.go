package cryptoh

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

var (
	testDir string
)

func init() {
	t := testing.T{}
	testDir = t.TempDir()
}

// ExampleSha256FileHash is an example of computing a SHA256 hash on a file.
// Verify the value by stopping the code in the debugger generating the hash:
// shasum -a 256 <testFilePath>
func ExampleSha256FileHash() {
	testFilePath := filepath.Join(testDir, "testHash.txt")
	os.WriteFile(testFilePath, []byte("this is a test"), 0644)
	hash, err := Sha256FileHash(testFilePath)
	if err != nil {
		fmt.Println("Error getting hash....")
	}
	fmt.Printf("%x", hash)

	// Output:
	// 2e99758548972a8e8822ad47fa1017ff72f06f3ff6a016851f45c398732bc50c
}
