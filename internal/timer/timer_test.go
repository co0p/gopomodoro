package timer

// Note: These tests use package timer (not timer_test) because they need to
// access private fields (mu, remaining) for comprehensive testing of the timer
// state machine. This is an accepted exception per the ADR on testing approach.

import (
	"sync"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	tmr := New()

	if tmr == nil {
		t.Fatal("New() returned nil")
	}

	if tmr.GetState() != StateIdle {
		t.Errorf("Expected initial state to be StateIdle, got %v", tmr.GetState())
	}

	if tmr.GetRemaining() != 1500 {
		t.Errorf("Expected initial remaining to be 1500 seconds, got %d", tmr.GetRemaining())
	}
}

func TestStart(t *testing.T) {
	tmr := New()

	startedCalled := false
	var mu sync.Mutex
	tmr.OnStarted(func(sessionType string, durationSeconds int) {
		mu.Lock()
		startedCalled = true
		mu.Unlock()
	})

	tmr.Start("work", 1500)

	if tmr.GetState() != StateRunning {
		t.Errorf("Expected state to be StateRunning after Start(), got %v", tmr.GetState())
	}

	mu.Lock()
	called := startedCalled
	mu.Unlock()
	if !called {
		t.Error("OnStarted callback was not called")
	}

	// Starting again should be a no-op
	startedCalledAgain := false
	tmr.OnStarted(func(sessionType string, durationSeconds int) {
		mu.Lock()
		startedCalledAgain = true
		mu.Unlock()
	})
	tmr.Start("work", 1500)

	mu.Lock()
	called = startedCalledAgain
	mu.Unlock()
	if called {
		t.Error("OnStarted should not be called when already running")
	}
}

func TestTick(t *testing.T) {
	tmr := New()

	var mu sync.Mutex
	tickCount := 0
	lastRemaining := 0
	tmr.OnTick(func(remaining int) {
		mu.Lock()
		tickCount++
		lastRemaining = remaining
		mu.Unlock()
	})

	// Use short duration for faster testing
	tmr.Start("work", 3)

	// Wait for a few ticks
	time.Sleep(2100 * time.Millisecond)

	mu.Lock()
	count := tickCount
	last := lastRemaining
	mu.Unlock()

	if count < 2 {
		t.Errorf("Expected at least 2 ticks, got %d", count)
	}

	if last >= 3 {
		t.Errorf("Expected remaining to decrease, got %d", last)
	}
}

func TestPauseAndResume(t *testing.T) {
	tmr := New()

	tmr.Start("work", 10)

	time.Sleep(100 * time.Millisecond)

	tmr.Pause()

	if tmr.GetState() != StatePaused {
		t.Errorf("Expected state to be StatePaused, got %v", tmr.GetState())
	}

	remainingAfterPause := tmr.GetRemaining()
	time.Sleep(1100 * time.Millisecond)
	remainingAfterWait := tmr.GetRemaining()

	if remainingAfterPause != remainingAfterWait {
		t.Errorf("Timer should not tick while paused. Before: %d, After: %d",
			remainingAfterPause, remainingAfterWait)
	}

	tmr.Resume()

	if tmr.GetState() != StateRunning {
		t.Errorf("Expected state to be StateRunning after Resume(), got %v", tmr.GetState())
	}

	time.Sleep(1100 * time.Millisecond)
	remainingAfterResume := tmr.GetRemaining()

	if remainingAfterResume >= remainingAfterWait {
		t.Error("Timer should tick after resume")
	}
}

func TestReset(t *testing.T) {
	tmr := New()

	tmr.Start("work", 1500)
	time.Sleep(100 * time.Millisecond)

	tmr.Reset()

	if tmr.GetState() != StateIdle {
		t.Errorf("Expected state to be StateIdle after Reset(), got %v", tmr.GetState())
	}

	if tmr.GetRemaining() != 0 {
		t.Errorf("Expected remaining to be reset to 0, got %d", tmr.GetRemaining())
	}

	// Reset from paused state
	tmr.Start("work", 1500)
	time.Sleep(100 * time.Millisecond)
	tmr.Pause()
	tmr.Reset()

	if tmr.GetState() != StateIdle {
		t.Error("Reset from paused should return to idle")
	}
}

func TestReset_ClearsSessionType(t *testing.T) {
	tmr := New()

	tmr.Start("work", 1500)
	time.Sleep(100 * time.Millisecond)

	tmr.Reset()

	if tmr.GetSessionType() != "" {
		t.Errorf("Expected GetSessionType() to return empty string after Reset, got %q", tmr.GetSessionType())
	}

	if tmr.GetRemaining() != 0 {
		t.Errorf("Expected GetRemaining() to return 0 after Reset, got %d", tmr.GetRemaining())
	}

	if tmr.GetState() != StateIdle {
		t.Errorf("Expected state to be StateIdle after Reset, got %v", tmr.GetState())
	}
}

func TestCompletion(t *testing.T) {
	tmr := New()

	// Set remaining to 1 for quick completion
	tmr.mu.Lock()
	tmr.remaining = 1
	tmr.mu.Unlock()

	var mu sync.Mutex
	completedCalled := false
	tmr.OnCompleted(func() {
		mu.Lock()
		completedCalled = true
		mu.Unlock()
	})

	tmr.Start("work", 1)

	// Wait for completion
	time.Sleep(1200 * time.Millisecond)

	mu.Lock()
	called := completedCalled
	mu.Unlock()

	if !called {
		t.Error("OnCompleted callback was not called")
	}

	if tmr.GetState() != StateIdle {
		t.Errorf("Expected state to be StateIdle after completion, got %v", tmr.GetState())
	}
}

func TestGetSessionType_ReturnsEmptyWhenIdle(t *testing.T) {
	tmr := New()

	sessionType := tmr.GetSessionType()

	if sessionType != "" {
		t.Errorf("Expected GetSessionType() to return empty string for new timer, got %q", sessionType)
	}
}

func TestStart_SetsSessionTypeAndDuration(t *testing.T) {
	tmr := New()

	tmr.Start("work", 1500)

	if tmr.GetSessionType() != "work" {
		t.Errorf("Expected GetSessionType() to return \"work\", got %q", tmr.GetSessionType())
	}

	if tmr.GetRemaining() != 1500 {
		t.Errorf("Expected GetRemaining() to return 1500, got %d", tmr.GetRemaining())
	}

	if tmr.GetState() != StateRunning {
		t.Errorf("Expected state to be StateRunning, got %v", tmr.GetState())
	}
}

func TestOnStarted_CallbackReceivesSessionContext(t *testing.T) {
	tmr := New()

	var receivedSessionType string
	var receivedDuration int
	var mu sync.Mutex

	tmr.OnStarted(func(sessionType string, durationSeconds int) {
		mu.Lock()
		receivedSessionType = sessionType
		receivedDuration = durationSeconds
		mu.Unlock()
	})

	tmr.Start("work", 1500)

	mu.Lock()
	sessionType := receivedSessionType
	duration := receivedDuration
	mu.Unlock()

	if sessionType != "work" {
		t.Errorf("Expected callback to receive sessionType \"work\", got %q", sessionType)
	}

	if duration != 1500 {
		t.Errorf("Expected callback to receive duration 1500, got %d", duration)
	}
}
