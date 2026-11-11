package plugin

import (
	"os"
	"strings"
	"testing"
)

func renderRunScript(t *testing.T, ocmd *OpenShiftTestsRunCommand) string {
	t.Helper()
	// OpenShiftTestsRunFile lives in /tmp/shared, created by the plugin
	// runtime; ensure it exists for the test.
	if err := os.MkdirAll(SharedDir, 0o755); err != nil {
		t.Fatalf("error creating shared dir: %v", err)
	}
	if err := ocmd.Create(); err != nil {
		t.Fatalf("error creating run file: %v", err)
	}
	content, err := os.ReadFile(OpenShiftTestsRunFile)
	if err != nil {
		t.Fatalf("error reading run file: %v", err)
	}
	return string(content)
}

func TestRunCommandProviderExternal(t *testing.T) {
	ocmd := NewOpenShiftRunCommand("run", "openshift/conformance")
	ocmd.Provider = `{"type":"external"}`

	script := renderRunScript(t, ocmd)

	if !strings.Contains(script, `--provider='{"type":"external"}'`) {
		t.Errorf("expected provider flag in run script, got:\n%s", script)
	}
}

func TestRunCommandProviderEmpty(t *testing.T) {
	ocmd := NewOpenShiftRunCommand("run", "openshift/conformance")

	script := renderRunScript(t, ocmd)

	if strings.Contains(script, "--provider") {
		t.Errorf("expected no provider flag in run script, got:\n%s", script)
	}
	if !strings.Contains(script, `--junit-dir="`) {
		t.Errorf("expected junit-dir flag in run script, got:\n%s", script)
	}
}
