// Package jsonh provides helper functions for the encoding/json package.
package jsonh

import "regexp"

// PrettyJSON transforms JSON for more friendly screen output.
// Transforms this:
// "SomeJSONField": [1,
//
//	2,
//	3,
//	4      ],
//
// into:
// "SomeJSONField": [1,2,3,4],
func PrettyJSON(json []byte) []byte {
	// re1: remove all CRLF from lines that only have a number followed by
	// comma. This gets rid of all CRLF, but leaves the initial CRLF
	// after the opening "["
	re1 := regexp.MustCompile(`(?m:^\s*?([0-9.]+,?)\s*?\r?\n?)`)
	json = re1.ReplaceAll(json, []byte("$1"))
	// re:2 Now get rid of the CRLF immediately after "[" if it is followed by a
	//  number and comma.
	re2 := regexp.MustCompile(`(?m:\[\s*?\r?\n?([0-9.]+,)\r?\n?)`)
	json = re2.ReplaceAll(json, []byte("[$1"))
	// re3: remove the trailing spaces after the final number and prior to the final "]"
	re3 := regexp.MustCompile(`([0-9.])\s*?]`)
	json = re3.ReplaceAll(json, []byte("$1]"))

	// JSON converts actual \n to "\n"; undo that
	re4 := regexp.MustCompile(`\n`)
	json = re4.ReplaceAll(json, []byte("\n"))

	// Remove trailing whitespace from any line so that output is
	// compatible with Golang Examples.
	re5 := regexp.MustCompile(`(?m)\s*?$`)
	json = re5.ReplaceAll(json, []byte(""))

	return json
}
