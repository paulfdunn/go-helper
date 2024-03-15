package jsonh

import "fmt"

func ExamplePrettyJSON() {
	testJSON := []byte(
		`"field": [
1.1,
2,
3    ],`)
	pj := PrettyJSON(testJSON)
	fmt.Println(string(pj))

	testJSON = []byte(
		`"field": [ 1.1,
2,
3    ],`)
	pj = PrettyJSON(testJSON)
	fmt.Println(string(pj))

	testJSON = []byte(
		`"field": [ 1.1,

		2,
3    ],`)
	pj = PrettyJSON(testJSON)
	fmt.Println(string(pj))

	// Output:
	// "field": [1.1,2,3],
	// "field": [1.1,2,3],
	// "field": [1.1,2,3],
}
