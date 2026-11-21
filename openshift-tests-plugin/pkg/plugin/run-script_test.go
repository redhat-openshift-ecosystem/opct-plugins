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

func TestNewOpenShiftRunCommand(t *testing.T) {
	tests := []struct {
		name     string
		provider string
		// wantContains lists substrings that must be present in the rendered script.
		wantContains []string
		// wantNotContains lists substrings that must be absent from the rendered script.
		wantNotContains []string
	}{
		{
			name:         "external provider sets --provider flag",
			provider:     `{"type":"external"}`,
			wantContains: []string{`--provider='{"type":"external"}'`},
		},
		{
			name:            "empty provider omits --provider flag",
			provider:        "",
			wantContains:    []string{`--junit-dir="`},
			wantNotContains: []string{"--provider"},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ocmd := NewOpenShiftRunCommand("run", "openshift/conformance")
			ocmd.Provider = tc.provider

			script := renderRunScript(t, ocmd)

			for _, want := range tc.wantContains {
				if !strings.Contains(script, want) {
					t.Errorf("expected %q in run script, got:\n%s", want, script)
				}
			}
			for _, unwanted := range tc.wantNotContains {
				if strings.Contains(script, unwanted) {
					t.Errorf("expected %q to be absent from run script, got:\n%s", unwanted, script)
				}
			}
		})
	}
}
