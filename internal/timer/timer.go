package timer

import (
	"sync"
	"time"
)

// State represents the current state of the timer
type State int

const (
	StateIdle State = iota
	StateRunning
	StatePaused

	defaultDuration = 1500 // 25 minutes in seconds
	tickInterval    = 10 * time.Second
)

// Timer manages a countdown timer
type Timer struct {
	mu          sync.Mutex
	state       State
	sessionType string
	remaining   int
	ticker      *time.Ticker
	stopChan    chan bool
	onStarted   func(string, int)
	onTick      func(int)
	onCompleted func()
}

// New creates a new Timer initialized to idle state with 25 minutes
func New() *Timer {
	return &Timer{
		state:     StateIdle,
		remaining: defaultDuration,
		stopChan:  make(chan bool, 1),
	}
}

// GetState returns the current state of the timer
func (t *Timer) GetState() State {
	t.mu.Lock()
	defer t.mu.Unlock()
	return t.state
}

// GetRemaining returns the remaining time in seconds
func (t *Timer) GetRemaining() int {
	t.mu.Lock()
	defer t.mu.Unlock()
	return t.remaining
}

// GetSessionType returns the current session type
func (t *Timer) GetSessionType() string {
	t.mu.Lock()
	defer t.mu.Unlock()
	return t.sessionType
}

// OnStarted registers a callback to be called when the timer starts
func (t *Timer) OnStarted(handler func(string, int)) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.onStarted = handler
}

// OnTick registers a callback to be called on each tick
func (t *Timer) OnTick(handler func(int)) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.onTick = handler
}

// OnCompleted registers a callback to be called when the timer completes
func (t *Timer) OnCompleted(handler func()) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.onCompleted = handler
}

// Start begins the timer countdown with the specified session type and duration
func (t *Timer) Start(sessionType string, durationSeconds int) {
	t.mu.Lock()
	defer t.mu.Unlock()

	// Only start if currently idle (no-op otherwise)
	if t.state != StateIdle {
		return
	}

	t.sessionType = sessionType
	t.remaining = durationSeconds
	t.state = StateRunning
	t.startTicker()

	if t.onStarted != nil {
		t.onStarted(sessionType, durationSeconds)
	}
}

// Pause pauses the timer
func (t *Timer) Pause() {
	t.mu.Lock()
	defer t.mu.Unlock()

	// Only pause if currently running
	if t.state != StateRunning {
		return
	}

	t.stopTicker()
	t.state = StatePaused
}

// Resume resumes the timer from paused state
func (t *Timer) Resume() {
	t.mu.Lock()
	defer t.mu.Unlock()

	// Only resume if currently paused
	if t.state != StatePaused {
		return
	}

	t.state = StateRunning
	t.startTicker()
}

// Reset resets the timer to idle state
func (t *Timer) Reset() {
	t.mu.Lock()
	defer t.mu.Unlock()

	t.stopTicker()
	t.state = StateIdle
	t.sessionType = ""
	t.remaining = 0
}

// startTicker starts the ticker goroutine (must be called with lock held)
func (t *Timer) startTicker() {
	t.ticker = time.NewTicker(tickInterval)
	go t.tickLoop()
}

// stopTicker stops the ticker goroutine (must be called with lock held)
func (t *Timer) stopTicker() {
	if t.ticker != nil {
		t.ticker.Stop()
		t.stopChan <- true
		t.ticker = nil
	}
}

// tickLoop runs the ticker loop in a goroutine
func (t *Timer) tickLoop() {
	for {
		select {
		case _, ok := <-t.ticker.C:
			if !ok {
				// Ticker was closed
				return
			}

			t.mu.Lock()

			if t.state != StateRunning {
				t.mu.Unlock()
				return
			}

			t.remaining -= 10
			if t.remaining < 0 {
				t.remaining = 0
			}
			currentRemaining := t.remaining

			// Check for completion
			if t.remaining == 0 {
				t.handleCompletion()
				t.mu.Unlock()
				return
			}

			// Call onTick callback
			if t.onTick != nil {
				callback := t.onTick
				t.mu.Unlock()
				callback(currentRemaining)
			} else {
				t.mu.Unlock()
			}

		case <-t.stopChan:
			return
		}
	}
}

// handleCompletion handles timer completion (must be called with lock held)
func (t *Timer) handleCompletion() {
	t.stopTicker()
	t.state = StateIdle

	if t.onCompleted != nil {
		callback := t.onCompleted
		// Call callback without lock to avoid deadlock
		go callback()
	}
}
