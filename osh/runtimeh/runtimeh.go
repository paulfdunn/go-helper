// Package runtimeh provides features that extend os/runtime.
package runtimeh

import (
	"fmt"
	"path/filepath"
	"runtime"
)

// SourceInfo returns a string with "file: FILENAME| func: FUNCTION_NAME| line: LINE_NUMBER|"
// of the caller.
func SourceInfo() string {
	file, fnc, line := sourceInfoCommon()
	return fmt.Sprintf("| file: %s| func: %s| line: %d|", file, fnc, line)
}

// SourceInfoError wraps the provided error with:
// "file: FILENAME| func: FUNCTION_NAME| line: LINE_NUMBER| desc: DESCRIPTION| error: ERROR|"
// of the caller.
func SourceInfoError(description string, err error) error {
	if err == nil {
		return nil
	}
	file, fnc, line := sourceInfoCommon()
	return fmt.Errorf("| file: %s| func: %s| line: %d| desc: %s| error: %w|", file, fnc, line, description, err)
}

func sourceInfoCommon() (file string, function string, line int) {
	var ok bool
	var fnc uintptr
	fnc, file, line, ok = runtime.Caller(2)
	if ok {
		// Get only the function name; FuncForPC will return something like filename.functionname.
		// I.E. runtimeh.SourceInfo
		fs := filepath.Ext(runtime.FuncForPC(fnc).Name())
		if len(fs) >= 1 {
			fs = fs[1:]
		}
		return filepath.Base(file), fs, line
	}
	return "", "", 0
}
