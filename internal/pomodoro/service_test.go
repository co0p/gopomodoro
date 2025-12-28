package pomodoro

import (
	"testing"
	"time"

	"github.com/benbjohnson/clock"
)

// MockNotifier is a test implementation of Notifier
type MockNotifier struct {
	SessionStartedCalls   int
	SessionTickCalls      int
	SessionCompletedCalls int
	StateChangedCalls     int
	LastSessionType       string
	LastDuration          int
	LastRemaining         int
	LastState             State
	Events                []string
}

func (m *MockNotifier) SessionStarted(sessionType string, duration int) {
	m.SessionStartedCalls++
	m.LastSessionType = sessionType
	m.LastDuration = duration
	m.Events = append(m.Events, "SessionStarted")
}

func (m *MockNotifier) SessionTick(remainingSeconds int) {
	m.SessionTickCalls++
	m.LastRemaining = remainingSeconds
	m.Events = append(m.Events, "SessionTick")
}

func (m *MockNotifier) SessionCompleted(sessionType string) {
	m.SessionCompletedCalls++
	m.LastSessionType = sessionType
	m.Events = append(m.Events, "SessionCompleted")
}

func (m *MockNotifier) StateChanged(state State) {
	m.StateChangedCalls++
	m.LastState = state
	m.Events = append(m.Events, "StateChanged")
}

func TestNewService(t *testing.T) {
	mockClock := clock.NewMock()
	service := NewService(mockClock, nil, nil)

	if service.GetState() != StateIdle {
		t.Errorf("expected StateIdle, got %v", service.GetState())
	}
}

func TestStartSession_FirstWorkSession(t *testing.T) {
	mockClock := clock.NewMock()
	mockNotifier := &MockNotifier{}
	service := NewService(mockClock, nil, mockNotifier)

	err := service.StartSession()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if service.GetState() != StateRunning {
		t.Errorf("expected StateRunning, got %v", service.GetState())
	}

	if mockNotifier.SessionStartedCalls != 1 {
		t.Errorf("expected 1 SessionStarted call, got %d", mockNotifier.SessionStartedCalls)
	}

	if mockNotifier.LastSessionType != "work" {
		t.Errorf("expected session type 'work', got '%s'", mockNotifier.LastSessionType)
	}

	if mockNotifier.LastDuration != 1500 {
		t.Errorf("expected duration 1500, got %d", mockNotifier.LastDuration)
	}
}

func TestStartSession_WhenAlreadyRunning_ReturnsError(t *testing.T) {
	mockClock := clock.NewMock()
	service := NewService(mockClock, nil, nil)
	service.StartSession()

	err := service.StartSession()
	if err != ErrAlreadyRunning {
		t.Errorf("expected ErrAlreadyRunning, got %v", err)
	}
}

func TestPauseSession_WhenRunning_Succeeds(t *testing.T) {
	mockClock := clock.NewMock()
	mockNotifier := &MockNotifier{}
	service := NewService(mockClock, nil, mockNotifier)
	service.StartSession()

	err := service.PauseSession()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if service.GetState() != StatePaused {
		t.Errorf("expected StatePaused, got %v", service.GetState())
	}

	if mockNotifier.StateChangedCalls != 1 {
		t.Errorf("expected 1 StateChanged call, got %d", mockNotifier.StateChangedCalls)
	}
}

func TestPauseSession_WhenIdle_ReturnsError(t *testing.T) {
	mockClock := clock.NewMock()
	service := NewService(mockClock, nil, nil)

	err := service.PauseSession()
	if err != ErrNotRunning {
		t.Errorf("expected ErrNotRunning, got %v", err)
	}
}

func TestResumeSession_WhenPaused_Succeeds(t *testing.T) {
	mockClock := clock.NewMock()
	mockNotifier := &MockNotifier{}
	service := NewService(mockClock, nil, mockNotifier)
	service.StartSession()
	service.PauseSession()

	err := service.ResumeSession()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if service.GetState() != StateRunning {
		t.Errorf("expected StateRunning, got %v", service.GetState())
	}

	if mockNotifier.StateChangedCalls != 2 {
		t.Errorf("expected 2 StateChanged calls, got %d", mockNotifier.StateChangedCalls)
	}
}

func TestResumeSession_WhenNotPaused_ReturnsError(t *testing.T) {
	mockClock := clock.NewMock()
	service := NewService(mockClock, nil, nil)

	err := service.ResumeSession()
	if err != ErrNotPaused {
		t.Errorf("expected ErrNotPaused, got %v", err)
	}
}

func TestSkipSession_AdvancesCycle(t *testing.T) {
	mockClock := clock.NewMock()
	mockNotifier := &MockNotifier{}
	service := NewService(mockClock, nil, mockNotifier)
	service.StartSession()

	err := service.SkipSession()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if service.GetState() != StateIdle {
		t.Errorf("expected StateIdle, got %v", service.GetState())
	}

	if service.GetCompletedSessions() != 1 {
		t.Errorf("expected 1 completed session, got %d", service.GetCompletedSessions())
	}

	if service.GetCurrentSessionType() != "short_break" {
		t.Errorf("expected 'short_break', got '%s'", service.GetCurrentSessionType())
	}
}

func TestResetCycle_ResetsToInitialState(t *testing.T) {
	mockClock := clock.NewMock()
	service := NewService(mockClock, nil, nil)
	service.StartSession()
	service.SkipSession() // Advance to break

	err := service.ResetCycle()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if service.GetState() != StateIdle {
		t.Errorf("expected StateIdle, got %v", service.GetState())
	}

	if service.GetCompletedSessions() != 0 {
		t.Errorf("expected 0 completed sessions, got %d", service.GetCompletedSessions())
	}

	if service.GetCurrentSessionType() != "work" {
		t.Errorf("expected 'work', got '%s'", service.GetCurrentSessionType())
	}
}

func TestQueryMethods_ReturnCorrectValues(t *testing.T) {
	mockClock := clock.NewMock()
	service := NewService(mockClock, nil, nil)
	service.StartSession()

	if service.GetRemainingSeconds() != 1500 {
		t.Errorf("expected 1500 remaining seconds, got %d", service.GetRemainingSeconds())
	}

	if service.GetCurrentSessionType() != "work" {
		t.Errorf("expected 'work', got '%s'", service.GetCurrentSessionType())
	}

	if service.GetCompletedSessions() != 0 {
		t.Errorf("expected 0 completed sessions, got %d", service.GetCompletedSessions())
	}

	progress := service.GetCycleProgress()
	if progress != "Session 1/4  🍅○○○" {
		t.Errorf("expected 'Session 1/4  🍅○○○', got '%s'", progress)
	}
}

func TestTimerCompletion_AdvancesToNextSession(t *testing.T) {
	mockClock := clock.NewMock()
	mockNotifier := &MockNotifier{}
	service := NewService(mockClock, nil, mockNotifier)

	service.StartSession() // Start work session

	// Advance clock to complete the session
	mockClock.Add(25 * time.Minute)

	// Wait a moment for the timer goroutine to process
	time.Sleep(50 * time.Millisecond)

	if service.GetCompletedSessions() != 1 {
		t.Errorf("expected 1 completed session, got %d", service.GetCompletedSessions())
	}

	if service.GetCurrentSessionType() != "short_break" {
		t.Errorf("expected 'short_break', got '%s'", service.GetCurrentSessionType())
	}

	// Check that SessionCompleted was called
	hasCompleted := false
	for _, event := range mockNotifier.Events {
		if event == "SessionCompleted" {
			hasCompleted = true
			break
		}
	}
	if !hasCompleted {
		t.Error("expected SessionCompleted event")
	}
}
