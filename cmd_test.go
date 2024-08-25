package main

import (
	"io"
	"os"
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

	oldStdout := os.Stdout
	oldStderr := os.Stderr

	stdoutReader, stdoutWriter, err := os.Pipe()
	if err != nil {
		t.Fatalf("os.Pipe() stdout error = %v", err)
	}
	stderrReader, stderrWriter, err := os.Pipe()
	if err != nil {
		t.Fatalf("os.Pipe() stderr error = %v", err)
	}

	os.Stdout = stdoutWriter
	os.Stderr = stderrWriter
	defer func() {
		os.Stdout = oldStdout
		os.Stderr = oldStderr
	}()

	stdoutCh := make(chan []byte, 1)
	stderrCh := make(chan []byte, 1)
	readErrCh := make(chan error, 2)
	go func() {
		stdoutBytes, readErr := io.ReadAll(stdoutReader)
		if readErr != nil {
			readErrCh <- readErr
			return
		}
		stdoutCh <- stdoutBytes
	}()
	go func() {
		stderrBytes, readErr := io.ReadAll(stderrReader)
		if readErr != nil {
			readErrCh <- readErr
			return
		}
		stderrCh <- stderrBytes
	}()

	err = execute(comment)

	if closeErr := stdoutWriter.Close(); closeErr != nil {
		t.Fatalf("stdoutWriter.Close() error = %v", closeErr)
	}
	if closeErr := stderrWriter.Close(); closeErr != nil {
		t.Fatalf("stderrWriter.Close() error = %v", closeErr)
	}

	var stdoutBytes, stderrBytes []byte
	for i := 0; i < 2; i++ {
		select {
		case stdoutBytes = <-stdoutCh:
		case stderrBytes = <-stderrCh:
		case readErr := <-readErrCh:
			t.Fatalf("io.ReadAll(...) error = %v", readErr)
		}
	}

	return string(stdoutBytes), string(stderrBytes), err
}
