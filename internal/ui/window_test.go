package ui

// Note: These tests use package ui (not ui_test) because they test private
// helper functions (formatTime, shouldStartBeEnabled, etc.) that are implementation
// details not exposed in the public API. This is an accepted exception per the ADR.

import (
	"testing"

	"github.com/co0p/gopomodoro/internal/timer"
)

func TestFormatTime(t *testing.T) {
	tests := []struct {
		seconds  int
		expected string
	}{
		{1500, "25:00"},
		{1499, "24:59"},
		{61, "01:01"},
		{3, "00:03"},
		{0, "00:00"},
		{3599, "59:59"},
	}

	for _, tt := range tests {
		result := formatTime(tt.seconds)
		if result != tt.expected {
			t.Errorf("formatTime(%d) = %s; want %s", tt.seconds, result, tt.expected)
		}
	}
}

func TestButtonStateLogic(t *testing.T) {
	tests := []struct {
		state        timer.State
		startEnabled bool
		pauseEnabled bool
		resetEnabled bool
	}{
		{timer.StateIdle, true, false, false},
		{timer.StateRunning, false, true, true},
		{timer.StatePaused, true, false, true},
	}

	for _, tt := range tests {
		if got := shouldStartBeEnabled(tt.state); got != tt.startEnabled {
			t.Errorf("shouldStartBeEnabled(%v) = %v; want %v", tt.state, got, tt.startEnabled)
		}
		if got := shouldPauseBeEnabled(tt.state); got != tt.pauseEnabled {
			t.Errorf("shouldPauseBeEnabled(%v) = %v; want %v", tt.state, got, tt.pauseEnabled)
		}
		if got := shouldResetBeEnabled(tt.state); got != tt.resetEnabled {
			t.Errorf("shouldResetBeEnabled(%v) = %v; want %v", tt.state, got, tt.resetEnabled)
		}
	}
}

func TestStateTransitions(t *testing.T) {
	// Test: Idle -> Running -> Paused -> Running -> Paused
	// This ensures button states are correct through pause/resume cycles

	// Idle state
	if !shouldStartBeEnabled(timer.StateIdle) {
		t.Error("Start should be enabled in Idle state")
	}
	if shouldPauseBeEnabled(timer.StateIdle) {
		t.Error("Pause should be disabled in Idle state")
	}

	// After Start -> Running
	if shouldStartBeEnabled(timer.StateRunning) {
		t.Error("Start should be disabled in Running state")
	}
	if !shouldPauseBeEnabled(timer.StateRunning) {
		t.Error("Pause should be enabled in Running state")
	}
	if !shouldResetBeEnabled(timer.StateRunning) {
		t.Error("Reset should be enabled in Running state")
	}

	// After Pause -> Paused
	if !shouldStartBeEnabled(timer.StatePaused) {
		t.Error("Start should be enabled in Paused state (for resume)")
	}
	if shouldPauseBeEnabled(timer.StatePaused) {
		t.Error("Pause should be disabled in Paused state")
	}
	if !shouldResetBeEnabled(timer.StatePaused) {
		t.Error("Reset should be enabled in Paused state")
	}

	// After Resume -> Running again
	if shouldStartBeEnabled(timer.StateRunning) {
		t.Error("Start should be disabled after resume (Running state)")
	}
	if !shouldPauseBeEnabled(timer.StateRunning) {
		t.Error("Pause should be enabled after resume (Running state)")
	}
	if !shouldResetBeEnabled(timer.StateRunning) {
		t.Error("Reset should be enabled after resume (Running state)")
	}
}
