package exech

import (
	"strings"
	"testing"
)

func TestExecCommand(t *testing.T) {
	var cmd, testStdout, testStderr string
	cmd = "echo"
	testStdout = "this is stdout"
	testStderr = ""
	so, se, rc, err := ExecCommand(cmd, []string{testStdout})
	// STDOUT comes back with a trailing \n
	if strings.TrimSpace(so) != testStdout || strings.TrimSpace(se) != testStderr || rc != 0 || err != nil {
		t.Errorf("echo failed, so: %s, se: %s, rc: %d, err: %v", so, se, rc, err)
	}

	// Verify pipe does not work and the string goes to STDOUT
	cmd = "echo"
	testStdout = "this is stdout; echo this is stderr 1>&2"
	testStderr = ""
	so, se, rc, err = ExecCommand(cmd, []string{testStdout})
	// STDOUT comes back with a trailing \n
	if strings.TrimSpace(so) != testStdout || strings.TrimSpace(se) != testStderr || rc != 0 || err != nil {
		t.Errorf("echo failed, so: %s, se: %s, rc: %d, err: %v", so, se, rc, err)
	}

	cmd = "echo-blah"
	testStdout = "this is stdout"
	testStderr = ""
	so, se, rc, err = ExecCommand(cmd, []string{testStdout})
	if rc == 0 || err == nil {
		t.Errorf("echo failed, so: %s, se: %s, rc: %d, err: %v", so, se, rc, err)
	}
}

func TestExecShell(t *testing.T) {
	var cmd, testStdout, testStderr string
	cmd = "echo this is stdout"
	testStdout = "this is stdout"
	testStderr = ""
	so, se, rc, err := ExecShell(cmd, []string{})
	// STDOUT comes back with a trailing \n
	if strings.TrimSpace(so) != testStdout || strings.TrimSpace(se) != testStderr || rc != 0 || err != nil {
		t.Errorf("echo failed, so: %s, se: %s, rc: %d, err: %v", so, se, rc, err)
	}

	cmd = "echo this is stderr 1>&2"
	testStdout = ""
	testStderr = "this is stderr"
	so, se, rc, err = ExecShell(cmd, []string{})
	// STDOUT comes back with a trailing \n
	if strings.TrimSpace(so) != testStdout || strings.TrimSpace(se) != testStderr || rc != 0 || err != nil {
		t.Errorf("echo failed, so: %s, se: %s, rc: %d, err: %v", so, se, rc, err)
	}

	cmd = "echo-blah this is stderr 1>&2"
	testStdout = ""
	testStderr = "this is stderr"
	so, se, rc, err = ExecShell(cmd, []string{})
	if rc == 0 || err == nil {
		t.Errorf("echo failed, so: %s, se: %s, rc: %d, err: %v", so, se, rc, err)
	}
}
