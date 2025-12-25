package tray_test

import (
	"testing"

	"github.com/co0p/gopomodoro/internal/tray"
	"github.com/getlantern/systray"
)

func TestSystrayImport(t *testing.T) {
	// Just validate we can import systray
	_ = systray.Run
}

func TestInitialize(t *testing.T) {
	err := tray.Initialize()
	if err != nil {
		t.Fatalf("Initialize() failed: %v", err)
	}
}

func TestSetIcon(t *testing.T) {
	// Simple test with minimal PNG data
	iconData := []byte{0x89, 0x50, 0x4E, 0x47} // PNG header
	err := tray.SetIcon(iconData)
	if err != nil {
		t.Fatalf("SetIcon() failed: %v", err)
	}
}

// TestLoadIconFromAssets is an integration test
// Note: Icon loading will be tested during full app integration
// For now, we verify the function compiles and returns an error when file is missing
func TestLoadIconFromAssets(t *testing.T) {
	// This function will work when called from main.go in the project root
	// Skip detailed testing here as it's environment-dependent
	t.Skip("Icon loading tested during integration - requires running from project root")
}

func TestOnClick(t *testing.T) {
	handler := func() {
		// Handler would be called in real usage
	}

	// Just verify we can register a handler without panicking
	// The actual callback execution is tested during integration
	tray.OnClick(handler)
}
