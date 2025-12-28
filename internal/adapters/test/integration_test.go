package test

import (
	"testing"
	"time"

	"github.com/benbjohnson/clock"
	"github.com/co0p/gopomodoro/internal/pomodoro"
)

// TestFullPomodoroCycle_FourSessions verifies a complete pomodoro cycle:
// 4 work sessions with 3 short breaks and 1 long break
func TestFullPomodoroCycle_FourSessions(t *testing.T) {
	mockClock := clock.NewMock()
	driver := NewTestDriver()
	service := pomodoro.NewService(mockClock, nil, driver)

	// Work session 1
	t.Log("Starting work session 1")
	if err := service.StartSession(); err != nil {
		t.Fatalf("Failed to start session 1: %v", err)
	}

	if err := driver.WaitForEvent(EventSessionStarted, 1*time.Second); err != nil {
		t.Fatalf("Session 1 start event not received: %v", err)
	}

	// Complete work session 1 (25 minutes)
	mockClock.Add(25 * time.Minute)
	if err := driver.WaitForEvent(EventSessionCompleted, 1*time.Second); err != nil {
		t.Fatalf("Session 1 completion event not received: %v", err)
	}

	if service.GetCompletedSessions() != 1 {
		t.Errorf("Expected 1 completed session, got %d", service.GetCompletedSessions())
	}

	if service.GetCurrentSessionType() != "short_break" {
		t.Errorf("Expected short_break, got %s", service.GetCurrentSessionType())
	}

	driver.ClearEvents()

	// Short break 1 (300 seconds = 30 ticks of 10 seconds each)
	t.Log("Starting short break 1")

	if err := service.StartSession(); err != nil {
		t.Fatalf("Failed to start short break 1: %v", err)
	}

	if err := driver.WaitForEvent(EventSessionStarted, 1*time.Second); err != nil {
		t.Fatalf("Short break 1 start event not received: %v", err)
	}

	// Advance clock and give time for goroutines to process each tick
	mockClock.Add(5*time.Minute + 1*time.Second) // Add extra second to ensure completion
	time.Sleep(600 * time.Millisecond)           // Give enough time for all ticks to process

	if err := driver.WaitForEvent(EventSessionCompleted, 2*time.Second); err != nil {
		t.Logf("Remaining seconds: %d", service.GetRemainingSeconds())
		t.Logf("Current state: %v", service.GetState())
		t.Logf("Events: %d", len(driver.GetEvents()))
		for i, e := range driver.GetEvents() {
			t.Logf("  Event %d: %v", i, e.Type)
		}
		t.Fatalf("Short break 1 completion event not received: %v", err)
	}

	if service.GetCurrentSessionType() != "work" {
		t.Errorf("Expected work session, got %s", service.GetCurrentSessionType())
	}

	driver.ClearEvents()

	// Work session 2
	t.Log("Starting work session 2")
	if err := service.StartSession(); err != nil {
		t.Fatalf("Failed to start session 2: %v", err)
	}

	if err := driver.WaitForEvent(EventSessionStarted, 1*time.Second); err != nil {
		t.Fatalf("Session 2 start event not received: %v", err)
	}

	mockClock.Add(25 * time.Minute)
	if err := driver.WaitForEvent(EventSessionCompleted, 1*time.Second); err != nil {
		t.Fatalf("Session 2 completion event not received: %v", err)
	}

	if service.GetCompletedSessions() != 2 {
		t.Errorf("Expected 2 completed sessions, got %d", service.GetCompletedSessions())
	}

	driver.ClearEvents()

	// Short break 2
	t.Log("Starting short break 2")
	if err := service.StartSession(); err != nil {
		t.Fatalf("Failed to start short break 2: %v", err)
	}

	if err := driver.WaitForEvent(EventSessionStarted, 1*time.Second); err != nil {
		t.Fatalf("Short break 2 start event not received: %v", err)
	}

	mockClock.Add(5 * time.Minute)
	if err := driver.WaitForEvent(EventSessionCompleted, 1*time.Second); err != nil {
		t.Fatalf("Short break 2 completion event not received: %v", err)
	}

	driver.ClearEvents()

	// Work session 3
	t.Log("Starting work session 3")
	if err := service.StartSession(); err != nil {
		t.Fatalf("Failed to start session 3: %v", err)
	}

	if err := driver.WaitForEvent(EventSessionStarted, 1*time.Second); err != nil {
		t.Fatalf("Session 3 start event not received: %v", err)
	}

	mockClock.Add(25 * time.Minute)
	if err := driver.WaitForEvent(EventSessionCompleted, 1*time.Second); err != nil {
		t.Fatalf("Session 3 completion event not received: %v", err)
	}

	if service.GetCompletedSessions() != 3 {
		t.Errorf("Expected 3 completed sessions, got %d", service.GetCompletedSessions())
	}

	driver.ClearEvents()

	// Short break 3
	t.Log("Starting short break 3")
	if err := service.StartSession(); err != nil {
		t.Fatalf("Failed to start short break 3: %v", err)
	}

	if err := driver.WaitForEvent(EventSessionStarted, 1*time.Second); err != nil {
		t.Fatalf("Short break 3 start event not received: %v", err)
	}

	mockClock.Add(5 * time.Minute)
	if err := driver.WaitForEvent(EventSessionCompleted, 1*time.Second); err != nil {
		t.Fatalf("Short break 3 completion event not received: %v", err)
	}

	driver.ClearEvents()

	// Work session 4
	t.Log("Starting work session 4")
	if err := service.StartSession(); err != nil {
		t.Fatalf("Failed to start session 4: %v", err)
	}

	if err := driver.WaitForEvent(EventSessionStarted, 1*time.Second); err != nil {
		t.Fatalf("Session 4 start event not received: %v", err)
	}

	mockClock.Add(25 * time.Minute)
	if err := driver.WaitForEvent(EventSessionCompleted, 1*time.Second); err != nil {
		t.Fatalf("Session 4 completion event not received: %v", err)
	}

	if service.GetCompletedSessions() != 4 {
		t.Errorf("Expected 4 completed sessions, got %d", service.GetCompletedSessions())
	}

	// Verify final state - should be long break
	if service.GetCurrentSessionType() != "long_break" {
		t.Errorf("Expected long_break after 4 sessions, got %s", service.GetCurrentSessionType())
	}

	t.Log("Full pomodoro cycle completed successfully")
}

// TestSkipSession_AdvancesCorrectly verifies skip functionality
func TestSkipSession_AdvancesCorrectly(t *testing.T) {
	mockClock := clock.NewMock()
	driver := NewTestDriver()
	service := pomodoro.NewService(mockClock, nil, driver)

	// Start work session 1
	service.StartSession()
	driver.WaitForEvent(EventSessionStarted, 1*time.Second)

	// Skip immediately
	if err := service.SkipSession(); err != nil {
		t.Fatalf("Failed to skip session: %v", err)
	}

	// Should advance to short break
	if service.GetCurrentSessionType() != "short_break" {
		t.Errorf("Expected short_break after skip, got %s", service.GetCurrentSessionType())
	}

	if service.GetCompletedSessions() != 1 {
		t.Errorf("Expected 1 completed session after skip, got %d", service.GetCompletedSessions())
	}

	// Verify skip event was recorded
	if !driver.HasEvent(EventSessionCompleted) {
		t.Error("Expected SessionCompleted event after skip")
	}
}

// TestPauseResume_MaintainsState verifies pause and resume functionality
func TestPauseResume_MaintainsState(t *testing.T) {
	mockClock := clock.NewMock()
	driver := NewTestDriver()
	service := pomodoro.NewService(mockClock, nil, driver)

	// Start session
	service.StartSession()
	driver.WaitForEvent(EventSessionStarted, 1*time.Second)

	if service.GetState() != pomodoro.StateRunning {
		t.Errorf("Expected StateRunning, got %v", service.GetState())
	}

	// Pause
	if err := service.PauseSession(); err != nil {
		t.Fatalf("Failed to pause: %v", err)
	}

	driver.WaitForEvent(EventStateChanged, 1*time.Second)

	if service.GetState() != pomodoro.StatePaused {
		t.Errorf("Expected StatePaused, got %v", service.GetState())
	}

	// Resume
	if err := service.ResumeSession(); err != nil {
		t.Fatalf("Failed to resume: %v", err)
	}

	driver.WaitForEvent(EventStateChanged, 1*time.Second)

	if service.GetState() != pomodoro.StateRunning {
		t.Errorf("Expected StateRunning after resume, got %v", service.GetState())
	}

	// Count state change events
	if driver.CountEvents(EventStateChanged) != 2 {
		t.Errorf("Expected 2 state change events, got %d", driver.CountEvents(EventStateChanged))
	}
}

// TestResetCycle_ResetsToInitialState verifies reset functionality
func TestResetCycle_ResetsToInitialState(t *testing.T) {
	mockClock := clock.NewMock()
	driver := NewTestDriver()
	service := pomodoro.NewService(mockClock, nil, driver)

	// Complete 2 work sessions to advance the cycle
	service.StartSession()
	mockClock.Add(25 * time.Minute)
	driver.WaitForEvent(EventSessionCompleted, 1*time.Second)

	service.StartSession() // Short break
	mockClock.Add(5 * time.Minute)
	driver.WaitForEvent(EventSessionCompleted, 1*time.Second)

	service.StartSession() // Work session 2
	mockClock.Add(25 * time.Minute)
	driver.WaitForEvent(EventSessionCompleted, 1*time.Second)

	if service.GetCompletedSessions() != 2 {
		t.Fatalf("Expected 2 completed sessions, got %d", service.GetCompletedSessions())
	}

	// Reset cycle
	if err := service.ResetCycle(); err != nil {
		t.Fatalf("Failed to reset cycle: %v", err)
	}

	// Verify reset state
	if service.GetState() != pomodoro.StateIdle {
		t.Errorf("Expected StateIdle after reset, got %v", service.GetState())
	}

	if service.GetCompletedSessions() != 0 {
		t.Errorf("Expected 0 completed sessions after reset, got %d", service.GetCompletedSessions())
	}

	if service.GetCurrentSessionType() != "work" {
		t.Errorf("Expected work session after reset, got %s", service.GetCurrentSessionType())
	}
}

// TestErrorHandling verifies state machine error handling
func TestErrorHandling(t *testing.T) {
	mockClock := clock.NewMock()
	service := pomodoro.NewService(mockClock, nil, nil)

	// Try to start when already running
	service.StartSession()
	err := service.StartSession()
	if err != pomodoro.ErrAlreadyRunning {
		t.Errorf("Expected ErrAlreadyRunning, got %v", err)
	}

	// Try to pause when idle
	service.ResetCycle()
	err = service.PauseSession()
	if err != pomodoro.ErrNotRunning {
		t.Errorf("Expected ErrNotRunning, got %v", err)
	}

	// Try to resume when not paused
	err = service.ResumeSession()
	if err != pomodoro.ErrNotPaused {
		t.Errorf("Expected ErrNotPaused, got %v", err)
	}

	// Try to skip when idle
	err = service.SkipSession()
	if err != pomodoro.ErrIdle {
		t.Errorf("Expected ErrIdle, got %v", err)
	}
}
