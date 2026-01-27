package gopomodoro_test

import (
	"testing"

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

	if c.RemainingMinutes() != int(gopomodoro.Pomodoro) {
		t.Fatalf("expected %d minutes remaining, got %d", gopomodoro.Pomodoro, c.RemainingMinutes())
	}
}

func TestTickDecrementsRemainingMinutes(t *testing.T) {
	c := &gopomodoro.Cycle{}
	c.Start()

	c.Tick()

	expected := int(gopomodoro.Pomodoro) - 1
	if c.RemainingMinutes() != expected {
		t.Fatalf("expected %d minutes remaining, got %d", expected, c.RemainingMinutes())
	}
}

func TestCycleTransitionsToIdleWhenRemainingReachesZero(t *testing.T) {
	c := &gopomodoro.Cycle{}
	c.Start()

	// Tick through all minutes
	for i := 0; i < int(gopomodoro.Pomodoro); i++ {
		c.Tick()
	}

	if !c.Is(gopomodoro.Idle) {
		t.Fatal("expected cycle to be idle after all ticks")
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
		State:     gopomodoro.Pomodoro,
		Remaining: 25,
		Observer:  observer,
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
