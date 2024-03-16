// Package exech provides wrapper functions for os/exec.
package exech

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
	"syscall"

	"github.com/paulfdunn/go-helper/osh/runtimeh"
)

const (
	ErrorWithNoReturnCode = -1
)

var (
	Shell = []string{"sh", "-c"}
)

// ExecCommand wraps os.exec to provide a function that returns: stdout, stderr, rc, err.
func ExecCommand(name string, args []string) (string, string, int, error) {
	cmd := exec.Command(name, args...)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		rerr := runtimeh.SourceInfoError(fmt.Sprintf("ExecCommand Run error: %s", stderr.String()), err)
		rc := ErrorWithNoReturnCode
		if exitError, ok := err.(*exec.ExitError); ok {
			rc = exitError.Sys().(syscall.WaitStatus).ExitStatus()
		}
		return "", "", rc, rerr
	}
	so := stdout.String()
	se := stderr.String()
	return so, se, 0, nil
}

// ExecCommand wraps os.exec to provide a function runs in a shell and
// that returns: stdout, stderr, rc, err.
func ExecShell(cmd string, args []string) (string, string, int, error) {
	var nargs []string
	cmdString := cmd + " " + strings.Join(args, " ")
	if len(Shell) > 1 {
		nargs = append(Shell[1:], cmdString)
	} else {
		nargs = []string{cmdString}
	}
	return ExecCommand(Shell[0], nargs)
}
