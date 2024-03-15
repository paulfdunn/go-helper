package osh

import (
	"fmt"
	"os"
	"os/user"
)

func ExampleByteSliceToIntSlice() {
	fmt.Println(ByteSliceToIntSlice([]byte{}))
	fmt.Println(ByteSliceToIntSlice([]byte{0}))
	fmt.Println(ByteSliceToIntSlice([]byte{0, 1, 2, 3, 4, 5}))

	// Output:
	// []
	// [0]
	// [0 1 2 3 4 5]
}

func ExampleByteSliceToString() {
	fmt.Print(ByteSliceToString([]byte{}, 3))
	fmt.Print(ByteSliceToString([]byte{0}, 3))
	fmt.Print(ByteSliceToString([]byte{0, 1, 2}, 3))
	fmt.Print(ByteSliceToString([]byte{0, 1, 2, 3}, 3))
	fmt.Print(ByteSliceToString([]byte{0, 1, 2, 3, 4, 5}, 3))

	// Output:
	// 00
	// 00 01 02
	// 00 01 02
	// 03
	// 00 01 02
	// 03 04 05
}

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
