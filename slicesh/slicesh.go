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
		out += fmt.Sprintf("%s", s[1:len(s)-1])
		out += "\n"
		index += bytesPerLine
	}
	return out
}
