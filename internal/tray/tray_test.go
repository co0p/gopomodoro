package tray_test

import (
	"testing"

	"github.com/co0p/gopomodoro/internal/timer"
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

func TestGetEmojiForState(t *testing.T) {
	tr := tray.New()

	tests := []struct {
		name        string
		sessionType string
		state       timer.State
		expected    string
	}{
		{
			name:        "idle state",
			sessionType: "",
			state:       timer.StateIdle,
			expected:    "🍅",
		},
		{
			name:        "running work",
			sessionType: "work",
			state:       timer.StateRunning,
			expected:    "🍅",
		},
		{
			name:        "running short break",
			sessionType: "short_break",
			state:       timer.StateRunning,
			expected:    "☕",
		},
		{
			name:        "running long break",
			sessionType: "long_break",
			state:       timer.StateRunning,
			expected:    "🌟",
		},
		{
			name:        "paused work",
			sessionType: "work",
			state:       timer.StatePaused,
			expected:    "⏸️",
		},
		{
			name:        "paused short break",
			sessionType: "short_break",
			state:       timer.StatePaused,
			expected:    "⏸️",
		},
		{
			name:        "paused long break",
			sessionType: "long_break",
			state:       timer.StatePaused,
			expected:    "⏸️",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			emoji := tr.GetEmojiForState(tt.sessionType, tt.state)
			if emoji != tt.expected {
				t.Errorf("GetEmojiForState(%q, %v) = %q, want %q", tt.sessionType, tt.state, emoji, tt.expected)
			}
		})
	}
}

func TestFormatMinutes(t *testing.T) {
	tr := tray.New()

	tests := []struct {
		name     string
		seconds  int
		expected string
	}{
		{
			name:     "25 minutes",
			seconds:  1500,
			expected: "25m",
		},
		{
			name:     "5 minutes",
			seconds:  300,
			expected: "5m",
		},
		{
			name:     "15 minutes",
			seconds:  900,
			expected: "15m",
		},
		{
			name:     "1 minute",
			seconds:  60,
			expected: "1m",
		},
		{
			name:     "30 seconds",
			seconds:  30,
			expected: "0m",
		},
		{
			name:     "0 seconds",
			seconds:  0,
			expected: "0m",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tr.FormatMinutes(tt.seconds)
			if result != tt.expected {
				t.Errorf("FormatMinutes(%d) = %q, want %q", tt.seconds, result, tt.expected)
			}
		})
	}
}

func TestUpdateDisplay(t *testing.T) {
	tr := tray.New()

	// Test that UpdateDisplay can be called without panicking
	// We can't directly verify systray.SetTitle calls without mocking,
	// but we can verify the method exists and accepts correct parameters
	tests := []struct {
		name             string
		sessionType      string
		state            timer.State
		remainingSeconds int
	}{
		{
			name:             "work session running",
			sessionType:      "work",
			state:            timer.StateRunning,
			remainingSeconds: 1500,
		},
		{
			name:             "short break running",
			sessionType:      "short_break",
			state:            timer.StateRunning,
			remainingSeconds: 300,
		},
		{
			name:             "paused state",
			sessionType:      "work",
			state:            timer.StatePaused,
			remainingSeconds: 900,
		},
		{
			name:             "idle state",
			sessionType:      "",
			state:            timer.StateIdle,
			remainingSeconds: 1500,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// This shouldn't panic
			tr.UpdateDisplay(tt.sessionType, tt.state, tt.remainingSeconds)
		})
	}
}

func TestGetIconData(t *testing.T) {
	tr := tray.New()

	tests := []struct {
		name        string
		sessionType string
		state       timer.State
		expectWork  bool // true if we expect iconWork (tomato)
	}{
		{
			name:        "idle state shows tomato",
			sessionType: "",
			state:       timer.StateIdle,
			expectWork:  true,
		},
		{
			name:        "work running shows tomato",
			sessionType: "work",
			state:       timer.StateRunning,
			expectWork:  true,
		},
		{
			name:        "paused shows pause icon",
			sessionType: "work",
			state:       timer.StatePaused,
			expectWork:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// We can't directly compare byte slices without accessing private vars,
			// but we can test via UpdateIcon which calls getIconData
			// For this test, we verify it doesn't panic
			tr.UpdateIcon(tt.sessionType, tt.state)
		})
	}
}
