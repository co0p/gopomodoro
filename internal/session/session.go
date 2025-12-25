// Package session manages pomodoro session state and cycle logic
package session

import "fmt"

const (
	TypeWork       = "work"
	TypeShortBreak = "short_break"
	TypeLongBreak  = "long_break"

	DurationWork       = 1500 // 25 minutes in seconds
	DurationShortBreak = 300  // 5 minutes in seconds
	DurationLongBreak  = 900  // 15 minutes in seconds

	SessionsPerCycle = 4
)

// Session manages the current pomodoro session state and cycle tracking
type Session struct {
	CurrentType           string
	CompletedWorkSessions int
}

// New creates a new Session initialized to work session
func New() *Session {
	return &Session{
		CurrentType:           TypeWork,
		CompletedWorkSessions: 0,
	}
}

// DetermineNext determines the next session type based on current state
// and returns the session type and its duration in seconds
func (s *Session) DetermineNext() (sessionType string, duration int) {
	switch s.CurrentType {
	case TypeWork:
		// After work session, check if it's time for long break
		if s.CompletedWorkSessions >= SessionsPerCycle {
			return TypeLongBreak, DurationLongBreak
		}
		return TypeShortBreak, DurationShortBreak
	case TypeShortBreak:
		return TypeWork, DurationWork
	case TypeLongBreak:
		s.CompletedWorkSessions = 0 // Reset cycle
		return TypeWork, DurationWork
	default:
		return TypeWork, DurationWork
	}
}

// IncrementCycle increments the completed work sessions counter
// Should be called when a work session completes or is skipped
func (s *Session) IncrementCycle() {
	if s.CurrentType == TypeWork {
		s.CompletedWorkSessions++
	}
}

// Reset resets the session to initial state
func (s *Session) Reset() {
	s.CurrentType = TypeWork
	s.CompletedWorkSessions = 0
}

// GetDuration returns the duration in seconds for the current session type
func (s *Session) GetDuration() int {
	switch s.CurrentType {
	case TypeWork:
		return DurationWork
	case TypeShortBreak:
		return DurationShortBreak
	case TypeLongBreak:
		return DurationLongBreak
	default:
		return DurationWork
	}
}

// FormatCycleIndicator formats the cycle progress indicator string
func (s *Session) FormatCycleIndicator() string {
	// Determine display session number and tomato count based on current state
	var displaySession int
	var tomatoCount int

	if s.CurrentType == TypeWork {
		// During work session, show current session number (which we're working on)
		displaySession = s.CompletedWorkSessions + 1
		// Show tomatoes including the current work session in progress
		tomatoCount = s.CompletedWorkSessions + 1
	} else {
		// During break, show the work session just completed
		displaySession = s.CompletedWorkSessions
		// Show tomatoes for completed sessions only
		tomatoCount = s.CompletedWorkSessions
	}

	// Build tomato progress string
	tomatoes := ""
	for i := 0; i < tomatoCount; i++ {
		tomatoes += "🍅"
	}
	for i := tomatoCount; i < SessionsPerCycle; i++ {
		tomatoes += "○"
	}

	return fmt.Sprintf("Session %d/%d  %s", displaySession, SessionsPerCycle, tomatoes)
}
