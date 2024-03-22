// Package exech provides wrapper functions for os/exec.
package exech

import (
	"context"
	"fmt"
	"io"
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

type ExecArgs struct {
	Args    []string
	Command string
	Context context.Context
	Stderr  io.Writer
	Stdout  io.Writer
}

// ExecCommandContext wraps os.exec.CommandContext to provide a function that returns: stdout, stderr, rc, err.
// This is a blocking call that only returns when the command completes.
// Callers that don't want to provide a context can pass in `context.TODO()`
func ExecCommandContext(execArgs *ExecArgs) (int, error) {
	cmd := exec.CommandContext(execArgs.Context, execArgs.Command, execArgs.Args...)
	cmd.Stderr = execArgs.Stderr
	cmd.Stdout = execArgs.Stdout
	err := cmd.Run()
	if err != nil {
		rerr := runtimeh.SourceInfoError(fmt.Sprintf("ExecCommandContext Run error"), err)
		rc := ErrorWithNoReturnCode
		if exitError, ok := err.(*exec.ExitError); ok {
			rc = exitError.Sys().(syscall.WaitStatus).ExitStatus()
		}
		return rc, rerr
	}
	return 0, nil
}

// ExecShellContext wraps os.exec.CommandContext to provide a function that runs in a shell and
// that returns: stdout, stderr, rc, err.
// Callers that don't want to provide a context can pass in `context.TODO()`
func ExecShellContext(execArgs *ExecArgs) (int, error) {
	cmdString := execArgs.Command + " " + strings.Join(execArgs.Args, " ")
	if len(Shell) > 1 {
		execArgs.Args = append(Shell[len(Shell)-1:], cmdString)
	} else {
		execArgs.Args = []string{cmdString}
	}
	execArgs.Command = Shell[0]
	return ExecCommandContext(execArgs)
}
