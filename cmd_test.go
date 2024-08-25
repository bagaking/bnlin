package main

import (
	"bytes"
	"errors"
	"os/exec"
	"strings"
	"testing"
)

func TestExecuteSkipsMarkdownFencesAndComments(t *testing.T) {
	tests := []struct {
		name    string
		comment string
		want    string
	}{
		{
			name: "plain bash fence",
			comment: "# The generated plan can include prose comments that should not execute.\n" +
				"```bash\n" +
				"printf 'ran-script-line\\n'\n" +
				"# This explanatory comment should not be sent to bash.\n" +
				"```\n",
			want: "ran-script-line\n",
		},
		{
			name: "bash fence with markdown attributes",
			comment: "```bash title=\"generated-script.sh\"\n" +
				"printf 'ran-attributed-fence\\n'\n" +
				"```\n",
			want: "ran-attributed-fence\n",
		},
		{
			name: "sh fence with markdown attributes",
			comment: "```sh linenos\n" +
				"printf 'ran-sh-fence\\n'\n" +
				"```\n",
			want: "ran-sh-fence\n",
		},
		{
			name: "plain sh fence",
			comment: "```sh\n" +
				"printf 'ran-plain-sh-fence\\n'\n" +
				"```\n",
			want: "ran-plain-sh-fence\n",
		},
		{
			name: "markdown fence inside here-doc",
			comment: "cat <<'EOF'\n" +
				"```bash title=\"x.sh\"\n" +
				"printf 'kept-here-doc-fence\\n'\n" +
				"```\n" +
				"EOF\n",
			want: "```bash title=\"x.sh\"\n" +
				"printf 'kept-here-doc-fence\\n'\n" +
				"```\n",
		},
		{
			name: "compact here-doc operator keeps markdown",
			comment: "cat<<'EOF'\n" +
				"```bash title=\"compact.sh\"\n" +
				"# kept as here-doc content\n" +
				"printf 'kept-compact-here-doc\\n'\n" +
				"```\n" +
				"EOF\n",
			want: "```bash title=\"compact.sh\"\n" +
				"# kept as here-doc content\n" +
				"printf 'kept-compact-here-doc\\n'\n" +
				"```\n",
		},
		{
			name: "tab stripped here-doc before attributed fence",
			comment: "cat <<-'EOF'\n" +
				"printf 'kept-tabbed-here-doc\\n'\n" +
				"\tEOF\n" +
				"```bash title=\"outer.sh\"\n" +
				"printf 'ran-after-tabbed-here-doc\\n'\n" +
				"```\n",
			want: "printf 'kept-tabbed-here-doc\\n'\n" +
				"ran-after-tabbed-here-doc\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stdout, stderr, err := captureExecuteOutput(t, tt.comment)
			if err != nil {
				t.Fatalf("execute(%q) error = %v, want nil", tt.comment, err)
			}

			if stdout != tt.want {
				t.Errorf("execute(%q).stdout = %q, want %q", tt.comment, stdout, tt.want)
			}
			if stderr != "" {
				t.Errorf("execute(%q).stderr = %q, want empty", tt.comment, stderr)
			}
		})
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
