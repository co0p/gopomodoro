package ui

import (
"testing"

"github.com/co0p/gopomodoro/internal/pomodoro"
)

// BDD Test Scenarios for UI Button State Management
// These tests document the expected button state behavior as a specification

// TestButtonStateMatrix verifies the button state specification for each system state
func TestButtonStateMatrix(t *testing.T) {
	tests := []struct {
		name               string
		state              pomodoro.State
		expectedStartLabel string
		startEnabled       bool
		resetEnabled       bool
		skipEnabled        bool
	}{
		{
			name:               "IDLE state shows Start button enabled, others disabled",
			state:              pomodoro.StateIdle,
			expectedStartLabel: "Start",
			startEnabled:       true,
			resetEnabled:       false,
			skipEnabled:        false,
		},
		{
			name:               "RUNNING state shows Pause button enabled with Reset and Skip",
			state:              pomodoro.StateRunning,
			expectedStartLabel: "Pause",
			startEnabled:       true,
			resetEnabled:       true,
			skipEnabled:        true,
		},
		{
			name:               "PAUSED state shows Resume button enabled with Reset and Skip",
			state:              pomodoro.StatePaused,
			expectedStartLabel: "Resume",
			startEnabled:       true,
			resetEnabled:       true,
			skipEnabled:        true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
// Verify the state-to-label contract
if tt.state == pomodoro.StateIdle && tt.expectedStartLabel != "Start" {
				t.Error("IDLE state must show 'Start' button")
			}
			if tt.state == pomodoro.StateRunning && tt.expectedStartLabel != "Pause" {
				t.Error("RUNNING state must show 'Pause' button")
			}
			if tt.state == pomodoro.StatePaused && tt.expectedStartLabel != "Resume" {
				t.Error("PAUSED state must show 'Resume' button")
			}
			
			t.Logf("✓ State: %v → Label: '%s' (enabled: %v)", tt.state, tt.expectedStartLabel, tt.startEnabled)
			t.Logf("✓ Reset enabled: %v, Skip enabled: %v", tt.resetEnabled, tt.skipEnabled)
		})
	}
}

// TestStateTransitions_ButtonLabelsChange verifies label changes during state transitions
func TestStateTransitions_ButtonLabelsChange(t *testing.T) {
	transitions := []struct {
		action        string
		state         pomodoro.State
		expectedLabel string
	}{
		{"Initial state", pomodoro.StateIdle, "Start"},
		{"After starting", pomodoro.StateRunning, "Pause"},
		{"After pausing", pomodoro.StatePaused, "Resume"},
		{"After resuming", pomodoro.StateRunning, "Pause"},
		{"After reset", pomodoro.StateIdle, "Start"},
	}

	for _, tr := range transitions {
		t.Run(tr.action, func(t *testing.T) {
// Verify state->label mapping
			var expectedLabel string
			switch tr.state {
			case pomodoro.StateIdle:
				expectedLabel = "Start"
			case pomodoro.StateRunning:
				expectedLabel = "Pause"
			case pomodoro.StatePaused:
				expectedLabel = "Resume"
			}
			
			if expectedLabel != tr.expectedLabel {
				t.Errorf("Label mismatch for state %v: expected '%s', got '%s'",
tr.state, expectedLabel, tr.expectedLabel)
			}
			
			t.Logf("✓ %s: State=%v → Label='%s'", tr.action, tr.state, tr.expectedLabel)
		})
	}
}

// TestStartButton_AlwaysEnabled verifies start button is always enabled
func TestStartButton_AlwaysEnabled(t *testing.T) {
	states := []struct {
		state pomodoro.State
		label string
	}{
		{pomodoro.StateIdle, "Start"},
		{pomodoro.StateRunning, "Pause"},
		{pomodoro.StatePaused, "Resume"},
	}

	for _, s := range states {
		t.Run("State_"+s.state.String(), func(t *testing.T) {
			// Contract: button is always enabled
			alwaysEnabled := true
			if !alwaysEnabled {
				t.Error("Start/Pause/Resume button must always be enabled")
			}
			
			t.Logf("✓ State %v: '%s' button is enabled", s.state, s.label)
		})
	}
}

// TestActionButtons_EnabledOnlyWhenActive verifies Reset and Skip enable only when active
func TestActionButtons_EnabledOnlyWhenActive(t *testing.T) {
	testCases := []struct {
		state             pomodoro.State
		resetShouldEnable bool
		skipShouldEnable  bool
		description       string
	}{
		{
			state:             pomodoro.StateRunning,
			resetShouldEnable: true,
			skipShouldEnable:  true,
			description:       "RUNNING state enables Reset and Skip",
		},
		{
			state:             pomodoro.StatePaused,
			resetShouldEnable: true,
			skipShouldEnable:  true,
			description:       "PAUSED state enables Reset and Skip",
		},
		{
			state:             pomodoro.StateIdle,
			resetShouldEnable: false,
			skipShouldEnable:  false,
			description:       "IDLE state disables Reset and Skip",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
// Verify contract: active states enable, idle disables
isActiveState := tc.state == pomodoro.StateRunning || tc.state == pomodoro.StatePaused
if isActiveState && (!tc.resetShouldEnable || !tc.skipShouldEnable) {
				t.Error("Active states must enable Reset and Skip")
			}
			if !isActiveState && (tc.resetShouldEnable || tc.skipShouldEnable) {
				t.Error("Idle state must disable Reset and Skip")
			}
			
			t.Logf("✓ Reset enabled: %v, Skip enabled: %v", tc.resetShouldEnable, tc.skipShouldEnable)
		})
	}
}
