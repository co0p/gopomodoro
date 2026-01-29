package tray_test

import (
	"testing"
	"time"

	gopomodoro "github.com/co0p/gopomodoro/pkg"
	"github.com/co0p/gopomodoro/pkg/tray"
)

func TestTray_GivenIdle_WhenDisplayed_ThenShowsOnlyIcon(t *testing.T) {
	formatter := tray.Formatter{}

	result := formatter.Format(gopomodoro.Idle, 0)

	expected := "🍅"
	if result != expected {
		t.Fatalf("expected %q, got %q", expected, result)
	}
}

func TestTray_GivenPomodoroRunning_WhenDisplayed_ThenShowsTimeWithMSuffix(t *testing.T) {
	formatter := tray.Formatter{}

	result := formatter.Format(gopomodoro.Pomodoro, 25*time.Minute)

	expected := "🍅 25m"
	if result != expected {
		t.Fatalf("expected %q, got %q", expected, result)
	}
}
