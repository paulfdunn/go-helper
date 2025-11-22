// Package slicesh provides helper functions for the slices package.
package slicesh

import "fmt"

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
		out += s[1 : len(s)-1]
		out += "\n"
		index += bytesPerLine
	}
	return out
}

// Combinations generates all combinations of size 'n' from the input array 'elements'.
func Combinations(elements []interface{}, n int) [][]interface{} {
	var result [][]interface{}
	// Call the recursive helper function to build combinations.
	generateCombinations(elements, n, 0, []interface{}{}, &result)
	return result
}

// generateCombinations is a recursive helper function to find combinations.
func generateCombinations(elements []interface{}, n int, start int, currentCombination []interface{}, result *[][]interface{}) {
	// Base case: if the current combination has 'n' elements, it's a valid combination.
	if len(currentCombination) == n {
		// Create a copy of the current combination and add it to the results.
		temp := make([]interface{}, n)
		copy(temp, currentCombination)
		*result = append(*result, temp)
		return
	}

	// If we've reached the end of the elements array, or if there aren't enough remaining
	// elements to form a combination of size 'n', stop.
	if start >= len(elements) || len(currentCombination)+(len(elements)-start) < n {
		return
	}

	// Option 1: Include the current element.
	generateCombinations(elements, n, start+1, append(currentCombination, elements[start]), result)

	// Option 2: Exclude the current element.
	generateCombinations(elements, n, start+1, currentCombination, result)
}
