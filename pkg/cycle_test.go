package gopomodoro_test

import (
	"testing"
	"time"

	gopomodoro "github.com/co0p/gopomodoro/pkg"
	mocks "github.com/co0p/gopomodoro/pkg/testing"
)

func TestNewCycleIsIdle(t *testing.T) {
	c := &gopomodoro.Cycle{}

	if !c.Is(gopomodoro.Idle) {
		t.Fatal("expected new cycle to be idle")
	}
}

func TestStartingIdleCycleTransitionsToPomodoro(t *testing.T) {
	c := &gopomodoro.Cycle{}

	c.Start()

	if !c.Is(gopomodoro.Pomodoro) {
		t.Fatal("expected cycle to be in pomodoro state after start")
	}
}

func TestStoppingRunningCycleTransitionsToIdle(t *testing.T) {
	c := &gopomodoro.Cycle{}
	c.Start()

	c.Stop()

	if !c.Is(gopomodoro.Idle) {
		t.Fatal("expected cycle to be idle after stop")
	}
}

func TestStartedPomodoroHas25MinutesRemaining(t *testing.T) {
	c := &gopomodoro.Cycle{}

	c.Start()

	expected := time.Duration(gopomodoro.Pomodoro) * time.Minute
	if c.Remaining() != expected {
		t.Fatalf("expected %v remaining, got %v", expected, c.Remaining())
	}
}

func TestTickDecrementsRemainingMinutes(t *testing.T) {
	c := &gopomodoro.Cycle{}
	c.Start()

	c.Tick()

	expected := time.Duration(gopomodoro.Pomodoro)*time.Minute - time.Minute
	if c.Remaining() != expected {
		t.Fatalf("expected %v remaining, got %v", expected, c.Remaining())
	}
}

func TestCycle_GivenPomodoroRunning_WhenTimerReachesZero_ThenShortBreakStarts(t *testing.T) {
	c := &gopomodoro.Cycle{}
	c.Start()

	// Tick through all pomodoro minutes
	for i := 0; i < int(gopomodoro.Pomodoro); i++ {
		c.Tick()
	}

	if !c.Is(gopomodoro.ShortBreak) {
		t.Fatalf("expected cycle to be in ShortBreak state, got %v", c.State)
	}

	expected := 5 * time.Minute
	if c.Remaining() != expected {
		t.Fatalf("expected %v remaining, got %v", expected, c.Remaining())
	}
}

func TestCycle_GivenShortBreakRunning_WhenTimerReachesZero_ThenReturnsToIdle(t *testing.T) {
	c := &gopomodoro.Cycle{
		State:    gopomodoro.ShortBreak,
		TimeLeft: 5 * time.Minute,
	}

	// Tick through all short break minutes
	for i := 0; i < int(gopomodoro.ShortBreak); i++ {
		c.Tick()
	}

	if !c.Is(gopomodoro.Idle) {
		t.Fatalf("expected cycle to be Idle, got %v", c.State)
	}

	if c.Remaining() != 0 {
		t.Fatalf("expected 0 remaining, got %v", c.Remaining())
	}
}

func TestCycle_GivenShortBreakRunning_WhenStopClicked_ThenReturnsToIdle(t *testing.T) {
	c := &gopomodoro.Cycle{
		State:    gopomodoro.ShortBreak,
		TimeLeft: 3 * time.Minute,
	}

	c.Stop()

	if !c.Is(gopomodoro.Idle) {
		t.Fatalf("expected cycle to be Idle, got %v", c.State)
	}

	if c.Remaining() != 0 {
		t.Fatalf("expected 0 remaining, got %v", c.Remaining())
	}
}

func TestStartNotifiesObserverOfStateChange(t *testing.T) {
	observer := &mocks.MockObserver{}
	c := gopomodoro.Cycle{Observer: observer}

	c.Start()

	if len(observer.StateChanges) != 1 {
		t.Fatalf("expected 1 state change, got %d", len(observer.StateChanges))
	}
	if observer.StateChanges[0] != gopomodoro.Pomodoro {
		t.Errorf("expected state change to Pomodoro, got %v", observer.StateChanges[0])
	}
}

func TestStopNotifiesObserverOfStateChange(t *testing.T) {
	observer := &mocks.MockObserver{}
	c := gopomodoro.Cycle{State: gopomodoro.Pomodoro, Observer: observer}
	c.Start()

	c.Stop()

	if len(observer.StateChanges) != 1 {
		t.Fatalf("expected 1 state change, got %d", len(observer.StateChanges))
	}
	if observer.StateChanges[0] != gopomodoro.Idle {
		t.Errorf("expected state change to Idle, got %v", observer.StateChanges[0])
	}
}

func TestTickNotifiesObserverOfStateChange(t *testing.T) {
	observer := &mocks.MockObserver{}
	c := gopomodoro.Cycle{
		State:    gopomodoro.Pomodoro,
		TimeLeft: 25 * time.Minute,
		Observer: observer,
	}

	c.Tick()

	// Tick should notify of state change
	if len(observer.StateChanges) != 1 {
		t.Fatalf("expected 1 state change, got %d", len(observer.StateChanges))
	}
	if observer.StateChanges[0] != gopomodoro.Pomodoro {
		t.Errorf("expected state change to Pomodoro, got %v", observer.StateChanges[0])
	}
}

func TestTickerFireTriggersTickAndNotifiesObserver(t *testing.T) {
	ticker := mocks.NewMockTicker()
	observer := &mocks.MockObserver{}
	c := gopomodoro.Cycle{
		Ticker:   ticker,
		Observer: observer,
	}

	c.Start()
	ticker.Fire()

	// Should have 2 state changes: one from Start, one from Tick
	if len(observer.StateChanges) != 2 {
		t.Fatalf("expected 2 state changes, got %d", len(observer.StateChanges))
	}
	if observer.StateChanges[0] != gopomodoro.Pomodoro {
		t.Errorf("expected first state change to Pomodoro, got %v", observer.StateChanges[0])
	}
	if observer.StateChanges[1] != gopomodoro.Pomodoro {
		t.Errorf("expected second state change to Pomodoro, got %v", observer.StateChanges[1])
	}
}
