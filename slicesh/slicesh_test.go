package slicesh

import "fmt"

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
