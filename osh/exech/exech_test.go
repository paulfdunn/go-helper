package exech

import (
	"bytes"
	"context"
	"strings"
	"testing"
	"time"
)

type testExecArgs struct {
	ExecArgs
	testStdout string
	testStderr string
}

func TestExecCommandContextAndShell(t *testing.T) {
	commonTests(t, ExecCommandContext)
	commonTests(t, ExecShellContext)
}

func TestExecCommandContextWithPipe(t *testing.T) {
	// Verify redirect DOES NOT WORK and the string goes to STDOUT.
	// (ExecShellContext is required for redirects or other shell supported functions.)
	var stdout, stderr bytes.Buffer
	tea := testExecArgs{
		ExecArgs: ExecArgs{
			Args:    []string{"this is stdout 1>&2"},
			Context: context.TODO(),
			Command: "echo",
			Stderr:  &stderr,
			Stdout:  &stdout,
		},
		testStdout: "this is stdout 1>&2",
		testStderr: "",
	}
	rc, err := ExecCommandContext(&tea.ExecArgs)
	if err != nil {
		t.Errorf("Error calling ExecCommandContext: %+v", err)
	}
	// STDOUT comes back with a trailing \n
	se := strings.TrimSpace(stderr.String())
	so := strings.TrimSpace(stdout.String())
	if se != tea.testStderr || strings.TrimSpace(so) != tea.testStdout || rc != 0 || err != nil {
		t.Errorf("echo failed, so: %s, se: %s, rc: %d, err: %v", so, se, rc, err)
	}
}

func TestExecShellContextWithPipe(t *testing.T) {
	// Verify redirect DOES WORK and the string goes to STDERR.
	var stdout, stderr bytes.Buffer
	tea := testExecArgs{
		ExecArgs: ExecArgs{
			Args:    []string{"this is stdout 1>&2"},
			Context: context.TODO(),
			Command: "echo",
			Stderr:  &stderr,
			Stdout:  &stdout,
		},
		testStdout: "",
		testStderr: "this is stdout",
	}
	rc, err := ExecShellContext(&tea.ExecArgs)
	if err != nil {
		t.Errorf("Error calling ExecCommandContext: %+v", err)
	}
	// STDOUT comes back with a trailing \n
	se := strings.TrimSpace(stderr.String())
	so := strings.TrimSpace(stdout.String())
	if se != tea.testStderr || strings.TrimSpace(so) != tea.testStdout || rc != 0 || err != nil {
		t.Errorf("echo failed, so: %s, se: %s, rc: %d, err: %v", so, se, rc, err)
	}
}

func commonTests(t *testing.T, testFunc func(*ExecArgs) (int, error)) {
	var stdout, stderr bytes.Buffer
	tea := testExecArgs{
		ExecArgs: ExecArgs{
			Args:    []string{"this is stdout"},
			Context: context.TODO(),
			Command: "echo",
			Stderr:  &stderr,
			Stdout:  &stdout,
		},
		testStdout: "this is stdout",
		testStderr: "",
	}
	rc, err := testFunc(&tea.ExecArgs)
	if err != nil {
		t.Errorf("Error calling testFunc: %+v", err)
	}
	// STDOUT comes back with a trailing \n
	se := strings.TrimSpace(stderr.String())
	so := strings.TrimSpace(stdout.String())
	if se != tea.testStderr || strings.TrimSpace(so) != tea.testStdout || rc != 0 || err != nil {
		t.Errorf("echo failed, so: %s, se: %s, rc: %d, err: %v", so, se, rc, err)
	}

	// Show that a timeout works to stop a command.
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(1)*time.Second)
	defer cancel()
	stdout.Truncate(0)
	stderr.Truncate(0)
	tea = testExecArgs{
		ExecArgs: ExecArgs{
			Args:    []string{"60"},
			Context: ctx,
			Command: "sleep",
			Stderr:  &stderr,
			Stdout:  &stdout,
		},
		testStdout: "this is stdout 1>&2",
		testStderr: "",
	}
	rc, err = testFunc(&tea.ExecArgs)
	// STDOUT comes back with a trailing \n
	se = strings.TrimSpace(stderr.String())
	so = strings.TrimSpace(stdout.String())
	if rc == 0 || err == nil {
		t.Errorf("echo failed, so: %s, se: %s, rc: %d, err: %v", so, se, rc, err)
	}

	// Negative test with a bad command
	stdout.Truncate(0)
	stderr.Truncate(0)
	tea = testExecArgs{
		ExecArgs: ExecArgs{
			Args:    []string{"this is stdout"},
			Context: context.TODO(),
			Command: "echo-blah",
			Stderr:  &stderr,
			Stdout:  &stdout,
		},
		testStdout: "",
		testStderr: "",
	}
	rc, err = testFunc(&tea.ExecArgs)
	// STDOUT comes back with a trailing \n
	se = strings.TrimSpace(stderr.String())
	so = strings.TrimSpace(stdout.String())
	if rc == 0 || err == nil {
		t.Errorf("echo failed, so: %s, se: %s, rc: %d, err: %v", so, se, rc, err)
	}
}
