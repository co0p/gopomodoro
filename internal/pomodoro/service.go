package pomodoro

import (
	"errors"
	"sync"

	"github.com/co0p/gopomodoro/internal/session"
	"github.com/co0p/gopomodoro/internal/timer"
)

var (
	ErrAlreadyRunning = errors.New("session already running")
	ErrNotRunning     = errors.New("no session running")
	ErrNotPaused      = errors.New("session is not paused")
	ErrIdle           = errors.New("session is idle")
)

// Service is the core pomodoro business logic service
type Service struct {
	session  *session.Session
	timer    *timer.Timer
	clock    Clock
	notifier Notifier
	storage  Storage
	state    State
	mu       sync.RWMutex
}

// NewService creates a new pomodoro service with the given dependencies
func NewService(clk Clock, storage Storage, notifier Notifier) *Service {
	sess := session.New()
	tmr := timer.NewWithClock(clk)

	s := &Service{
		session:  sess,
		timer:    tmr,
		clock:    clk,
		storage:  storage,
		notifier: notifier,
		state:    StateIdle,
	}

	// Register timer callbacks
	tmr.OnStarted(s.handleTimerStarted)
	tmr.OnTick(s.handleTimerTick)
	tmr.OnCompleted(s.handleTimerCompleted)

	return s
}

// SetNotifier sets the notifier (used for late binding in main)
func (s *Service) SetNotifier(notifier Notifier) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.notifier = notifier
}

// GetState returns the current state of the service
func (s *Service) GetState() State {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.state
}

// StartSession starts a new pomodoro session
func (s *Service) StartSession() error {
	s.mu.Lock()
	if s.state != StateIdle {
		s.mu.Unlock()
		return ErrAlreadyRunning
	}
	s.state = StateRunning
	sessionType := s.session.CurrentType
	duration := s.session.GetDuration()
	s.mu.Unlock()

	s.timer.Start(sessionType, duration)

	if s.notifier != nil {
		s.notifier.SessionStarted(sessionType, duration)
	}

	return nil
}

// PauseSession pauses the currently running session
func (s *Service) PauseSession() error {
	s.mu.Lock()
	if s.state != StateRunning {
		s.mu.Unlock()
		return ErrNotRunning
	}
	s.state = StatePaused
	s.mu.Unlock()

	s.timer.Stop()

	if s.notifier != nil {
		s.notifier.StateChanged(StatePaused)
	}

	return nil
}

// ResumeSession resumes a paused session
func (s *Service) ResumeSession() error {
	s.mu.Lock()
	if s.state != StatePaused {
		s.mu.Unlock()
		return ErrNotPaused
	}
	s.state = StateRunning
	sessionType := s.session.CurrentType
	remaining := s.timer.GetRemaining()
	s.mu.Unlock()

	s.timer.Start(sessionType, remaining)

	if s.notifier != nil {
		s.notifier.StateChanged(StateRunning)
	}

	return nil
}

// SkipSession skips the current session and advances to the next (but does not start it)
func (s *Service) SkipSession() error {
	s.mu.Lock()
	if s.state == StateIdle {
		s.mu.Unlock()
		return ErrIdle
	}
	currentType := s.session.CurrentType
	s.mu.Unlock()

	s.timer.Stop()
	s.advanceCycle()

	if s.notifier != nil {
		s.notifier.SessionCompleted(currentType)
	}

	return nil
}

// ResetCycle resets the pomodoro cycle to the beginning
func (s *Service) ResetCycle() error {
	s.timer.Stop()

	s.mu.Lock()
	s.session.Reset()
	s.state = StateIdle
	s.mu.Unlock()

	return nil
}

// GetRemainingSeconds returns the remaining seconds in the current session
func (s *Service) GetRemainingSeconds() int {
	return s.timer.GetRemaining()
}

// GetCurrentSessionType returns the current session type
func (s *Service) GetCurrentSessionType() string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.session.CurrentType
}

// GetCompletedSessions returns the number of completed work sessions
func (s *Service) GetCompletedSessions() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.session.CompletedWorkSessions
}

// GetCycleProgress returns the formatted cycle progress indicator
func (s *Service) GetCycleProgress() string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.session.FormatCycleIndicator()
}

// GetCurrentDuration returns the duration of the current session in seconds
func (s *Service) GetCurrentDuration() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.session.GetDuration()
}

// Internal timer event handlers

func (s *Service) handleTimerStarted(sessionType string, duration int) {
	// Already published in StartSession, no-op here
}

func (s *Service) handleTimerTick(remaining int) {
	if s.notifier != nil {
		s.notifier.SessionTick(remaining)
	}
}

func (s *Service) handleTimerCompleted() {
	s.mu.Lock()
	currentType := s.session.CurrentType
	s.mu.Unlock()

	s.advanceCycle()

	if s.notifier != nil {
		s.notifier.SessionCompleted(currentType)
	}
}

// advanceCycle increments the cycle counter and determines the next session
func (s *Service) advanceCycle() {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.session.IncrementCycle()
	nextType, _ := s.session.DetermineNext()
	s.session.CurrentType = nextType
	s.state = StateIdle
}
