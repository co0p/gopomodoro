package gopomodoro_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestCI_GivenPushToMain_WhenWorkflowConfigured_ThenTriggered(t *testing.T) {
	path := filepath.Join("..", ".github", "workflows", "ci.yml")

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("expected workflow at %s: %v", path, err)
	}

	content := string(data)
	if !strings.Contains(content, "push") || !strings.Contains(content, "main") {
		t.Fatalf("expected workflow to run on push to main")
	}
}

func TestCI_GivenWorkflow_WhenConfigured_ThenBuildsMacOSBundleAndUploadsArtifact(t *testing.T) {

	path := filepath.Join("..", ".github", "workflows", "ci.yml")

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("expected workflow at %s: %v", path, err)
	}

	content := string(data)
	if !strings.Contains(content, "make bundle") || !strings.Contains(content, "actions/upload-artifact") {
		t.Fatalf("expected workflow to build macos bundle and upload it")
	}
}

func TestCI_GivenWorkflow_WhenConfigured_ThenCreatesRelease(t *testing.T) {
	path := filepath.Join("..", ".github", "workflows", "ci.yml")

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("expected workflow at %s: %v", path, err)
	}

	content := string(data)
	if !strings.Contains(content, "softprops/action-gh-release") {
		t.Fatalf("expected workflow to create a release")
	}
}
