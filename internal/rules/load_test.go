package rules

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadFromFileYAML(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	path := filepath.Join(dir, "rules.yaml")

	content := []byte(`rules:
  - name: block-ads
    enabled: true
    priority: 1
    match:
      host_pattern: "*ads*"
    action:
      type: block
`)

	if err := os.WriteFile(path, content, 0o600); err != nil {
		t.Fatalf("failed writing test file: %v", err)
	}

	loaded, err := LoadFromFile(path)
	if err != nil {
		t.Fatalf("LoadFromFile failed: %v", err)
	}

	if len(loaded) != 1 {
		t.Fatalf("expected one loaded rule, got %d", len(loaded))
	}

	if loaded[0].Action.Type != ActionBlock {
		t.Fatalf("expected block action, got %q", loaded[0].Action.Type)
	}
}

func TestValidateRejectsUnknownAction(t *testing.T) {
	t.Parallel()

	err := Validate([]Rule{{
		Name:    "broken-rule",
		Enabled: true,
		Action: Action{
			Type: ActionType("explode"),
		},
	}})

	if err == nil {
		t.Fatalf("expected validation error, got nil")
	}
}
