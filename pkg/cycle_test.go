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

	c.AdvanceMinute()

	expected := time.Duration(gopomodoro.Pomodoro)*time.Minute - time.Minute
	if c.Remaining() != expected {
		t.Fatalf("expected %v remaining, got %v", expected, c.Remaining())
	}
}

func TestCycle_GivenPomodoroRunning_WhenTimerReachesZero_ThenShortBreakStarts(t *testing.T) {
	c := &gopomodoro.Cycle{}
	c.Start()

	mocks.CompleteCycle(c)

	if !c.Is(gopomodoro.ShortBreak) {
		t.Fatalf("expected cycle to be in ShortBreak state, got %v", c.State)
	}

	expected := 5 * time.Minute
	if c.Remaining() != expected {
		t.Fatalf("expected %v remaining, got %v", expected, c.Remaining())
	}
}

func TestCycle_GivenShortBreakRunning_WhenTimerReachesZero_ThenNextPomodoroStarts(t *testing.T) {
	c := &gopomodoro.Cycle{
		State:    gopomodoro.ShortBreak,
		TimeLeft: 5 * time.Minute,
	}

	mocks.CompleteCycle(c)

	if !c.Is(gopomodoro.Pomodoro) {
		t.Fatalf("expected cycle to automatically start next Pomodoro, got %v", c.State)
	}

	expected := 25 * time.Minute
	if c.Remaining() != expected {
		t.Fatalf("expected %v remaining, got %v", expected, c.Remaining())
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

	c.AdvanceMinute()

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

func TestCycle_Given3CompletedPomodoros_When4thCompletes_ThenLongBreakStarts(t *testing.T) {
	c := &gopomodoro.Cycle{}

	// Complete 4 pomodoros (each followed by a break, except the 4th)
	for i := 0; i < 4; i++ {
		c.Start()
		mocks.CompleteCycle(c)

		// After 1st, 2nd, 3rd pomodoro: should be in short break
		if i < 3 {
			if !c.Is(gopomodoro.ShortBreak) {
				t.Fatalf("pomodoro %d: expected ShortBreak, got %v", i+1, c.State)
			}
			// Complete the short break - should auto-start next pomodoro
			mocks.CompleteCycle(c)
			// Should automatically start next Pomodoro
			if !c.Is(gopomodoro.Pomodoro) {
				t.Fatalf("after break %d: expected Pomodoro to auto-start, got %v", i+1, c.State)
			}
		}
	}

	// After 4th pomodoro completes, should be in LongBreak
	if !c.Is(gopomodoro.LongBreak) {
		t.Fatalf("expected LongBreak after 4th pomodoro, got %v", c.State)
	}

	expected := 15 * time.Minute
	if c.Remaining() != expected {
		t.Fatalf("expected %v remaining in long break, got %v", expected, c.Remaining())
	}
}

func TestCycle_GivenLongBreakRunning_WhenTimerReachesZero_ThenReturnsToIdle(t *testing.T) {
	c := &gopomodoro.Cycle{
		State:    gopomodoro.LongBreak,
		TimeLeft: 15 * time.Minute,
	}

	mocks.CompleteCycle(c)

	if !c.Is(gopomodoro.Idle) {
		t.Fatalf("expected cycle to be Idle, got %v", c.State)
	}

	if c.Remaining() != 0 {
		t.Fatalf("expected 0 remaining, got %v", c.Remaining())
	}
}

func TestCycle_GivenLongBreakRunning_WhenStopClicked_ThenReturnsToIdle(t *testing.T) {
	c := &gopomodoro.Cycle{
		State:    gopomodoro.LongBreak,
		TimeLeft: 10 * time.Minute,
	}

	c.Stop()

	if !c.Is(gopomodoro.Idle) {
		t.Fatalf("expected cycle to be Idle, got %v", c.State)
	}

	if c.Remaining() != 0 {
		t.Fatalf("expected 0 remaining, got %v", c.Remaining())
	}
}

func TestCycle_GivenPomodoroRunning_WhenStopClicked_ThenCounterResets(t *testing.T) {
	c := &gopomodoro.Cycle{}

	// Complete 2 pomodoros to increment counter
	for i := 0; i < 2; i++ {
		c.Start()
		mocks.CompleteCycle(c)
		// Complete short break
		mocks.CompleteCycle(c)
	}

	// Start 3rd pomodoro and stop it mid-way
	c.Start()
	for j := 0; j < 10; j++ {
		c.AdvanceMinute()
	}

	c.Stop()

	// Start a new pomodoro and complete it
	c.Start()
	mocks.CompleteCycle(c)

	// Should be in ShortBreak (counter was reset, so this is the 1st pomodoro)
	if !c.Is(gopomodoro.ShortBreak) {
		t.Fatalf("expected ShortBreak after first pomodoro (counter reset), got %v", c.State)
	}
}
