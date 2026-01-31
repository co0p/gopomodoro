package gopomodoro_test

import (
	"testing"

	gopomodoro "github.com/co0p/gopomodoro/pkg"
	pomotest "github.com/co0p/gopomodoro/pkg/testing"
)

func TestSoundNotification_GivenPomodoroCompletes_WhenShortBreakStarts_ThenSoundPlays(t *testing.T) {
	ticker := &pomotest.MockTicker{}
	notifier := &pomotest.MockNotifier{}
	cycle := &gopomodoro.Cycle{
		Ticker:   ticker,
		Notifier: notifier,
	}

	cycle.Start()
	pomotest.CompleteCycle(cycle)

	if notifier.NotifyCallCount != 1 {
		t.Errorf("expected NotifyCallCount = 1, got %d", notifier.NotifyCallCount)
	}
}

func TestSoundNotification_GivenCycle4Completes_WhenLongBreakStarts_ThenSoundPlays(t *testing.T) {
	ticker := &pomotest.MockTicker{}
	notifier := &pomotest.MockNotifier{}
	cycle := &gopomodoro.Cycle{
		Ticker:   ticker,
		Notifier: notifier,
	}

	cycle.Start()
	for range 7 {
		pomotest.CompleteCycle(cycle)
	}

	// After 4th pomodoro completes, should be in LongBreak state
	if !cycle.Is(gopomodoro.LongBreak) {
		t.Errorf("expected state LongBreak, got %v", cycle.State)
	}

	// Should have gotten notification when transitioning to long break
	// (7 CompleteCycles = 4 pomodoro->break + 3 break->pomodoro = 7 notifications)
	if notifier.NotifyCallCount != 7 {
		t.Errorf("expected NotifyCallCount = 7, got %d", notifier.NotifyCallCount)
	}
}

func TestSoundNotification_GivenShortBreakCompletes_WhenPomodoroStarts_ThenSoundPlays(t *testing.T) {
	ticker := &pomotest.MockTicker{}
	notifier := &pomotest.MockNotifier{}
	cycle := &gopomodoro.Cycle{
		Ticker:   ticker,
		Notifier: notifier,
	}

	cycle.Start()
	pomotest.CompleteCycle(cycle) // Pomodoro → ShortBreak (notify #1)
	pomotest.CompleteCycle(cycle) // ShortBreak → Pomodoro (should notify #2)

	if notifier.NotifyCallCount != 2 {
		t.Errorf("expected NotifyCallCount = 2, got %d", notifier.NotifyCallCount)
	}
}

func TestSoundNotification_GivenLongBreakCompletes_WhenLongBreakEnds_ThenSoundPlays(t *testing.T) {
	ticker := &pomotest.MockTicker{}
	notifier := &pomotest.MockNotifier{}
	cycle := &gopomodoro.Cycle{
		Ticker:   ticker,
		Notifier: notifier,
	}

	cycle.Start()
	for range 8 {
		pomotest.CompleteCycle(cycle) // Complete full cycle including long break
	}

	// After long break completes, should be back in Idle state
	if !cycle.Is(gopomodoro.Idle) {
		t.Errorf("expected state Idle, got %v", cycle.State)
	}

	// Should have gotten 8 notifications (7 + 1 for long break ending)
	if notifier.NotifyCallCount != 8 {
		t.Errorf("expected NotifyCallCount = 8, got %d", notifier.NotifyCallCount)
	}
}
