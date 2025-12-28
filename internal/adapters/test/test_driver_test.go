package test

import (
	"testing"
	"time"

	"github.com/co0p/gopomodoro/internal/pomodoro"
)

func TestNewTestDriver(t *testing.T) {
	driver := NewTestDriver()
	if driver == nil {
		t.Fatal("expected driver to be created")
	}

	events := driver.GetEvents()
	if len(events) != 0 {
		t.Errorf("expected 0 events, got %d", len(events))
	}
}

func TestDriverRecordsEvents(t *testing.T) {
	driver := NewTestDriver()

	// Record a session started event
	driver.SessionStarted("work", 1500)

	events := driver.GetEvents()
	if len(events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(events))
	}

	if events[0].Type != EventSessionStarted {
		t.Errorf("expected EventSessionStarted, got %v", events[0].Type)
	}

	data := events[0].Data.(map[string]interface{})
	if data["type"] != "work" {
		t.Errorf("expected type 'work', got %v", data["type"])
	}
	if data["duration"] != 1500 {
		t.Errorf("expected duration 1500, got %v", data["duration"])
	}
}

func TestDriverRecordsMultipleEvents(t *testing.T) {
	driver := NewTestDriver()

	driver.SessionStarted("work", 1500)
	driver.SessionTick(1490)
	driver.SessionCompleted("work")
	driver.StateChanged(pomodoro.StateIdle)

	events := driver.GetEvents()
	if len(events) != 4 {
		t.Errorf("expected 4 events, got %d", len(events))
	}

	// Verify event order
	if events[0].Type != EventSessionStarted {
		t.Errorf("expected first event to be SessionStarted")
	}
	if events[1].Type != EventSessionTick {
		t.Errorf("expected second event to be SessionTick")
	}
	if events[2].Type != EventSessionCompleted {
		t.Errorf("expected third event to be SessionCompleted")
	}
	if events[3].Type != EventStateChanged {
		t.Errorf("expected fourth event to be StateChanged")
	}
}

func TestClearEvents(t *testing.T) {
	driver := NewTestDriver()

	driver.SessionStarted("work", 1500)
	driver.SessionTick(1490)

	if len(driver.GetEvents()) != 2 {
		t.Fatalf("expected 2 events before clear")
	}

	driver.ClearEvents()

	if len(driver.GetEvents()) != 0 {
		t.Errorf("expected 0 events after clear, got %d", len(driver.GetEvents()))
	}
}

func TestWaitForEvent_Success(t *testing.T) {
	driver := NewTestDriver()

	// Record event in background
	go func() {
		time.Sleep(50 * time.Millisecond)
		driver.SessionCompleted("work")
	}()

	err := driver.WaitForEvent(EventSessionCompleted, 200*time.Millisecond)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}

func TestWaitForEvent_Timeout(t *testing.T) {
	driver := NewTestDriver()

	err := driver.WaitForEvent(EventSessionCompleted, 100*time.Millisecond)
	if err == nil {
		t.Error("expected timeout error, got nil")
	}
}

func TestHasEvent(t *testing.T) {
	driver := NewTestDriver()

	if driver.HasEvent(EventSessionStarted) {
		t.Error("expected no SessionStarted event")
	}

	driver.SessionStarted("work", 1500)

	if !driver.HasEvent(EventSessionStarted) {
		t.Error("expected SessionStarted event to be present")
	}

	if driver.HasEvent(EventSessionCompleted) {
		t.Error("expected no SessionCompleted event")
	}
}

func TestCountEvents(t *testing.T) {
	driver := NewTestDriver()

	driver.SessionTick(1490)
	driver.SessionTick(1480)
	driver.SessionTick(1470)
	driver.SessionStarted("work", 1500)

	if driver.CountEvents(EventSessionTick) != 3 {
		t.Errorf("expected 3 tick events, got %d", driver.CountEvents(EventSessionTick))
	}

	if driver.CountEvents(EventSessionStarted) != 1 {
		t.Errorf("expected 1 started event, got %d", driver.CountEvents(EventSessionStarted))
	}

	if driver.CountEvents(EventSessionCompleted) != 0 {
		t.Errorf("expected 0 completed events, got %d", driver.CountEvents(EventSessionCompleted))
	}
}
