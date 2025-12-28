package pomodoro

import (
	"time"

	"github.com/benbjohnson/clock"
)

// Notifier is the outbound port for publishing events to the UI or other listeners
type Notifier interface {
	SessionStarted(sessionType string, duration int)
	SessionTick(remainingSeconds int)
	SessionCompleted(sessionType string)
	StateChanged(state State)
}

// Storage is the outbound port for persisting session data
type Storage interface {
	LogSession(timestamp time.Time, sessionType, status string, duration int) error
}

// Clock is an alias for the clock.Clock interface for time source abstraction
type Clock = clock.Clock

// State represents the current state of the pomodoro service
type State int

const (
	StateIdle State = iota
	StateRunning
	StatePaused
)

func (s State) String() string {
	switch s {
	case StateIdle:
		return "idle"
	case StateRunning:
		return "running"
	case StatePaused:
		return "paused"
	default:
		return "unknown"
	}
}
