package shellUtils

import (
	"strings"
	"testing"
)

func TestExecuteCommandAllowsSuccessfulStderr(t *testing.T) {
	out, err := ExecuteCommand(`printf 'normal output'; printf 'informational stderr' >&2`)
	if err != nil {
		t.Fatalf("ExecuteCommand returned unexpected error: %v", err)
	}
	if !strings.Contains(out, "normal output") || !strings.Contains(out, "informational stderr") {
		t.Fatalf("combined output missing stdout/stderr: %q", out)
	}
}

func TestExecuteCommandIncludesStdoutAndStderrOnFailure(t *testing.T) {
	out, err := ExecuteCommand(`printf 'stdout before failure'; printf 'stderr before failure' >&2; exit 8`)
	if err == nil {
		t.Fatalf("ExecuteCommand returned nil error for failing command, output=%q", out)
	}
	if !strings.Contains(out, "stdout before failure") || !strings.Contains(out, "stderr before failure") {
		t.Fatalf("output missing stdout/stderr: %q", out)
	}
	errText := err.Error()
	if !strings.Contains(errText, "exit status 8") || !strings.Contains(errText, "stdout before failure") || !strings.Contains(errText, "stderr before failure") {
		t.Fatalf("error missing exit status or combined output: %q", errText)
	}
}
