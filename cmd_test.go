package main

import (
	"bytes"
	"strings"
	"testing"
)

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

	for _, want := range []string{"stdout-line-0001", "stdout-line-2500", "stdout-line-5000"} {
		if !strings.Contains(stdout, want) {
			t.Fatalf("stdout missing %q", want)
		}
	}
	for _, want := range []string{"stderr-line-0001", "stderr-line-2500", "stderr-line-5000"} {
		if !strings.Contains(stderr, want) {
			t.Fatalf("stderr missing %q", want)
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
