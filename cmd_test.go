package main

import (
	"bytes"
	"errors"
	"os/exec"
	"strings"
	"testing"
)

func TestExecuteSkipsMarkdownFencesAndComments(t *testing.T) {
	comment := "# The generated plan can include prose comments that should not execute.\n" +
		"```bash\n" +
		"printf 'ran-script-line\\n'\n" +
		"# This explanatory comment should not be sent to bash.\n" +
		"```\n"
	stdout, stderr, err := captureExecuteOutput(t, comment)
	if err != nil {
		t.Fatalf("execute(...) error = %v, want nil", err)
	}

	if stdout != "ran-script-line\n" {
		t.Errorf("execute(...).stdout = %q, want %q", stdout, "ran-script-line\n")
	}
	if stderr != "" {
		t.Errorf("execute(...).stderr = %q, want empty", stderr)
	}
}

func TestExecuteFlushesFailedCommandOutput(t *testing.T) {
	stdout, stderr, err := captureExecuteOutput(t, `
for i in $(seq 1 5000); do
  printf 'stdout-line-%04d\n' "$i"
  printf 'stderr-line-%04d\n' "$i" >&2
done
exit 7
`)

	if err == nil {
		t.Fatal("execute(...) error = nil, want failure")
	}
	var exitErr *exec.ExitError
	if !errors.As(err, &exitErr) {
		t.Fatalf("execute(...) error = %v, want *exec.ExitError", err)
	}
	if exitErr.ExitCode() != 7 {
		t.Errorf("execute(...) exit code = %d, want 7", exitErr.ExitCode())
	}

	for _, want := range []string{"stdout-line-0001", "stdout-line-2500", "stdout-line-5000"} {
		if !strings.Contains(stdout, want) {
			t.Errorf("execute(...).stdout missing %q", want)
		}
	}
	for _, want := range []string{"stderr-line-0001", "stderr-line-2500", "stderr-line-5000"} {
		if !strings.Contains(stderr, want) {
			t.Errorf("execute(...).stderr missing %q", want)
		}
	}
}

func captureExecuteOutput(t *testing.T, comment string) (stdout string, stderr string, err error) {
	t.Helper()

	var stdoutBuffer bytes.Buffer
	var stderrBuffer bytes.Buffer

	err = executeWithWriters(comment, &stdoutBuffer, &stderrBuffer)

	return stdoutBuffer.String(), stderrBuffer.String(), err
}
