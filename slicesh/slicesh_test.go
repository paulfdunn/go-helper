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

func ExampleCombinations() {
	// Example usage with integers
	intArray := []interface{}{1, 2, 3, 4, 5}
	nInt := 3
	intCombinations := Combinations(intArray, nInt)
	fmt.Printf("Combinations of %d elements from %v:\n", nInt, intArray)
	for _, combo := range intCombinations {
		fmt.Println(combo)
	}

	fmt.Println("\n---")

	// Example usage with strings
	stringArray := []interface{}{"apple", "banana", "cherry", "date"}
	nString := 2
	stringCombinations := Combinations(stringArray, nString)
	fmt.Printf("Combinations of %d elements from %v:\n", nString, stringArray)
	for _, combo := range stringCombinations {
		fmt.Println(combo)
	}

	// Output:
	// Combinations of 3 elements from [1 2 3 4 5]:
	// [1 2 3]
	// [1 2 4]
	// [1 2 5]
	// [1 3 4]
	// [1 3 5]
	// [1 4 5]
	// [2 3 4]
	// [2 3 5]
	// [2 4 5]
	// [3 4 5]
	//
	// ---
	// Combinations of 2 elements from [apple banana cherry date]:
	// [apple banana]
	// [apple cherry]
	// [apple date]
	// [banana cherry]
	// [banana date]
	// [cherry date]
}
