package session

import (
	"testing"
)

func TestNew(t *testing.T) {
	s := New()
	if s.CurrentType != TypeWork {
		t.Errorf("New session should start with TypeWork, got %s", s.CurrentType)
	}
	if s.CompletedWorkSessions != 0 {
		t.Errorf("New session should start with 0 completed work sessions, got %d", s.CompletedWorkSessions)
	}
}

func TestDetermineNext(t *testing.T) {
	tests := []struct {
		name                  string
		currentType           string
		completedWorkSessions int
		expectedType          string
		expectedDuration      int
		expectedCompleted     int // after DetermineNext
	}{
		{
			name:                  "After 1st work session → short break",
			currentType:           TypeWork,
			completedWorkSessions: 1,
			expectedType:          TypeShortBreak,
			expectedDuration:      DurationShortBreak,
			expectedCompleted:     1,
		},
		{
			name:                  "After 4th work session → long break",
			currentType:           TypeWork,
			completedWorkSessions: 4,
			expectedType:          TypeLongBreak,
			expectedDuration:      DurationLongBreak,
			expectedCompleted:     4,
		},
		{
			name:                  "After short break → work",
			currentType:           TypeShortBreak,
			completedWorkSessions: 2,
			expectedType:          TypeWork,
			expectedDuration:      DurationWork,
			expectedCompleted:     2,
		},
		{
			name:                  "After long break → work (cycle reset)",
			currentType:           TypeLongBreak,
			completedWorkSessions: 4,
			expectedType:          TypeWork,
			expectedDuration:      DurationWork,
			expectedCompleted:     0, // Reset after long break
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Session{
				CurrentType:           tt.currentType,
				CompletedWorkSessions: tt.completedWorkSessions,
			}

			sessionType, duration := s.DetermineNext()

			if sessionType != tt.expectedType {
				t.Errorf("Expected session type %s, got %s", tt.expectedType, sessionType)
			}
			if duration != tt.expectedDuration {
				t.Errorf("Expected duration %d, got %d", tt.expectedDuration, duration)
			}
			if s.CompletedWorkSessions != tt.expectedCompleted {
				t.Errorf("Expected completed work sessions %d, got %d", tt.expectedCompleted, s.CompletedWorkSessions)
			}
		})
	}
}

func TestIncrementCycle(t *testing.T) {
	tests := []struct {
		name                  string
		currentType           string
		completedWorkSessions int
		expectedCompleted     int
	}{
		{
			name:                  "Increment during work session",
			currentType:           TypeWork,
			completedWorkSessions: 1,
			expectedCompleted:     2,
		},
		{
			name:                  "No increment during short break",
			currentType:           TypeShortBreak,
			completedWorkSessions: 2,
			expectedCompleted:     2,
		},
		{
			name:                  "No increment during long break",
			currentType:           TypeLongBreak,
			completedWorkSessions: 4,
			expectedCompleted:     4,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Session{
				CurrentType:           tt.currentType,
				CompletedWorkSessions: tt.completedWorkSessions,
			}

			s.IncrementCycle()

			if s.CompletedWorkSessions != tt.expectedCompleted {
				t.Errorf("Expected completed work sessions %d, got %d", tt.expectedCompleted, s.CompletedWorkSessions)
			}
		})
	}
}

func TestReset(t *testing.T) {
	s := &Session{
		CurrentType:           TypeShortBreak,
		CompletedWorkSessions: 3,
	}

	s.Reset()

	if s.CurrentType != TypeWork {
		t.Errorf("Reset should set CurrentType to TypeWork, got %s", s.CurrentType)
	}
	if s.CompletedWorkSessions != 0 {
		t.Errorf("Reset should set CompletedWorkSessions to 0, got %d", s.CompletedWorkSessions)
	}
}

func TestGetDuration(t *testing.T) {
	tests := []struct {
		name             string
		currentType      string
		expectedDuration int
	}{
		{
			name:             "Work session duration",
			currentType:      TypeWork,
			expectedDuration: DurationWork,
		},
		{
			name:             "Short break duration",
			currentType:      TypeShortBreak,
			expectedDuration: DurationShortBreak,
		},
		{
			name:             "Long break duration",
			currentType:      TypeLongBreak,
			expectedDuration: DurationLongBreak,
		},
		{
			name:             "Unknown type defaults to work duration",
			currentType:      "unknown",
			expectedDuration: DurationWork,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Session{
				CurrentType: tt.currentType,
			}

			duration := s.GetDuration()

			if duration != tt.expectedDuration {
				t.Errorf("Expected duration %d, got %d", tt.expectedDuration, duration)
			}
		})
	}
}

func TestFormatCycleIndicator(t *testing.T) {
	tests := []struct {
		name                  string
		currentType           string
		completedWorkSessions int
		expected              string
	}{
		{
			name:                  "First work session in progress",
			currentType:           TypeWork,
			completedWorkSessions: 0,
			expected:              "Session 1/4  🍅○○○",
		},
		{
			name:                  "Second work session in progress",
			currentType:           TypeWork,
			completedWorkSessions: 1,
			expected:              "Session 2/4  🍅🍅○○",
		},
		{
			name:                  "On short break after 1st work session",
			currentType:           TypeShortBreak,
			completedWorkSessions: 1,
			expected:              "Session 1/4  🍅○○○",
		},
		{
			name:                  "Fourth work session in progress",
			currentType:           TypeWork,
			completedWorkSessions: 3,
			expected:              "Session 4/4  🍅🍅🍅🍅",
		},
		{
			name:                  "On long break after 4th work session",
			currentType:           TypeLongBreak,
			completedWorkSessions: 4,
			expected:              "Session 4/4  🍅🍅🍅🍅",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Session{
				CurrentType:           tt.currentType,
				CompletedWorkSessions: tt.completedWorkSessions,
			}

			result := s.FormatCycleIndicator()

			if result != tt.expected {
				t.Errorf("Expected cycle indicator '%s', got '%s'", tt.expected, result)
			}
		})
	}
}
