package runtimeh

import (
	"fmt"
	"regexp"
)

// ExampleSourceInfo is an example of SourceInfo. Note that the line number is
// replaced to prevent the example failing as the line number changes.
func ExampleSourceInfo() {
	ffl := SourceInfo()
	nonum := regexp.MustCompile(`[0-9]+`).ReplaceAllString(ffl, "#")
	fmt.Println(nonum)
	// Output:
	// | file: runtimeh_test.go| func: ExampleSourceInfo| line: #|
}

// ExampleSourceInfoError is an example of SourceInfoError. Note that the line number is
// replaced to prevent the example failing as the line number changes.
func ExampleSourceInfoError() {
	ffl := SourceInfoError("there was an error", fmt.Errorf("example error")).Error()
	nonum := regexp.MustCompile(`[0-9]+`).ReplaceAllString(ffl, "#")
	fmt.Println(nonum)

	ffle := SourceInfoError("", nil)
	if ffle != nil {
		fmt.Println("this should never print")
	}
	// Output:
	// | file: runtimeh_test.go| func: ExampleSourceInfoError| line: #| desc: there was an error| error: example error|
}
