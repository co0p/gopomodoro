package test

import (
	"fmt"
	"sync"
	"time"

	"github.com/co0p/gopomodoro/internal/pomodoro"
)

// EventType represents the type of event recorded
type EventType int

const (
	EventSessionStarted EventType = iota
	EventSessionTick
	EventSessionCompleted
	EventStateChanged
)

func (e EventType) String() string {
	switch e {
	case EventSessionStarted:
		return "SessionStarted"
	case EventSessionTick:
		return "SessionTick"
	case EventSessionCompleted:
		return "SessionCompleted"
	case EventStateChanged:
		return "StateChanged"
	default:
		return "Unknown"
	}
}

// Event represents a recorded event from the service
type Event struct {
	Type      EventType
	Timestamp time.Time
	Data      interface{}
}

// TestDriver implements the Notifier interface and records all events for testing
type TestDriver struct {
	events []Event
	mu     sync.Mutex
}

// NewTestDriver creates a new test driver
func NewTestDriver() *TestDriver {
	return &TestDriver{
		events: []Event{},
	}
}

// SessionStarted records a session started event
func (d *TestDriver) SessionStarted(sessionType string, duration int) {
	d.recordEvent(EventSessionStarted, map[string]interface{}{
		"type":     sessionType,
		"duration": duration,
	})
}

// SessionTick records a tick event
func (d *TestDriver) SessionTick(remainingSeconds int) {
	d.recordEvent(EventSessionTick, map[string]interface{}{
		"remaining": remainingSeconds,
	})
}

// SessionCompleted records a session completed event
func (d *TestDriver) SessionCompleted(sessionType string) {
	d.recordEvent(EventSessionCompleted, map[string]interface{}{
		"type": sessionType,
	})
}

// StateChanged records a state change event
func (d *TestDriver) StateChanged(state pomodoro.State) {
	d.recordEvent(EventStateChanged, map[string]interface{}{
		"state": state,
	})
}

// recordEvent is a helper to record an event
func (d *TestDriver) recordEvent(eventType EventType, data interface{}) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.events = append(d.events, Event{
		Type:      eventType,
		Timestamp: time.Now(),
		Data:      data,
	})
}

// GetEvents returns a copy of all recorded events
func (d *TestDriver) GetEvents() []Event {
	d.mu.Lock()
	defer d.mu.Unlock()
	// Return a copy to avoid race conditions
	eventsCopy := make([]Event, len(d.events))
	copy(eventsCopy, d.events)
	return eventsCopy
}

// ClearEvents clears all recorded events
func (d *TestDriver) ClearEvents() {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.events = []Event{}
}

// WaitForEvent waits for a specific event type to be recorded, with timeout
func (d *TestDriver) WaitForEvent(eventType EventType, timeout time.Duration) error {
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		d.mu.Lock()
		for _, e := range d.events {
			if e.Type == eventType {
				d.mu.Unlock()
				return nil
			}
		}
		d.mu.Unlock()
		time.Sleep(10 * time.Millisecond)
	}
	return fmt.Errorf("event %v not received within timeout", eventType)
}

// HasEvent checks if a specific event type has been recorded
func (d *TestDriver) HasEvent(eventType EventType) bool {
	d.mu.Lock()
	defer d.mu.Unlock()
	for _, e := range d.events {
		if e.Type == eventType {
			return true
		}
	}
	return false
}

// CountEvents returns the count of a specific event type
func (d *TestDriver) CountEvents(eventType EventType) int {
	d.mu.Lock()
	defer d.mu.Unlock()
	count := 0
	for _, e := range d.events {
		if e.Type == eventType {
			count++
		}
	}
	return count
}
