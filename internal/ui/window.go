// Package ui manages the dropdown window/panel
package ui

import (
	"fmt"
	"log"
	"time"

	"github.com/co0p/gopomodoro/internal/session"
	"github.com/co0p/gopomodoro/internal/storage"
	"github.com/co0p/gopomodoro/internal/timer"
	"github.com/co0p/gopomodoro/internal/tray"
	"github.com/getlantern/systray"
)

// Window represents the dropdown panel
type Window struct {
	visible          bool
	header           *systray.MenuItem
	timerDisplay     *systray.MenuItem
	progressBar      *systray.MenuItem
	cycleIndicator   *systray.MenuItem
	btnStart         *systray.MenuItem
	btnPause         *systray.MenuItem
	btnReset         *systray.MenuItem
	btnSkip          *systray.MenuItem
	btnQuit          *systray.MenuItem
	timer            *timer.Timer
	session          *session.Session
	tray             *tray.Tray
	sessionStartTime time.Time
	sessionDuration  int
}

// CreateWindow initializes the dropdown window with placeholder UI
func CreateWindow() (*Window, error) {
	return &Window{
		visible: false,
		session: session.New(),
	}, nil
}

// Show displays the window at the specified screen coordinates
func (w *Window) Show(x, y int) error {
	w.visible = true
	log.Printf("[INFO] Window shown at position (x: %d, y: %d)", x, y)
	return nil
}

// Hide conceals the window
func (w *Window) Hide() error {
	w.visible = false
	log.Println("[INFO] Window hidden")
	return nil
}

// IsVisible returns whether the window is currently displayed
func (w *Window) IsVisible() bool {
	return w.visible
}

// InitializeMenu creates the actual menu items (called from systray.Run)
func (w *Window) InitializeMenu() {
	log.Println("[INFO] Initializing menu items...")

	w.header = systray.AddMenuItem("Ready", "Current state")
	w.header.Disable()

	w.timerDisplay = systray.AddMenuItem("25min", "Timer display")
	w.timerDisplay.Disable()

	w.progressBar = systray.AddMenuItem("○○○○○○○○○○", "Session progress")
	w.progressBar.Disable()

	w.cycleIndicator = systray.AddMenuItem("Session 1/4  🍅○○○", "Cycle progress")
	w.cycleIndicator.Disable()

	systray.AddSeparator()

	w.btnStart = systray.AddMenuItem("Start", "Start timer")
	w.btnPause = systray.AddMenuItem("Pause", "Pause timer")
	w.btnReset = systray.AddMenuItem("Reset", "Reset timer")

	systray.AddSeparator()

	w.btnSkip = systray.AddMenuItem("Skip", "Skip to next session")

	systray.AddSeparator()

	w.btnQuit = systray.AddMenuItem("Quit", "Quit the application")

	log.Println("[INFO] Menu items initialized")
}

// SetTimer sets the timer and registers event handlers
func (w *Window) SetTimer(t *timer.Timer) {
	w.timer = t

	// Register event handlers
	w.timer.OnStarted(w.handleTimerStarted)
	w.timer.OnTick(w.handleTimerTick)
	w.timer.OnCompleted(w.handleTimerCompleted)

	// Start button click handlers
	w.startClickHandlers()
}

// SetTray sets the tray reference for icon updates
func (w *Window) SetTray(t *tray.Tray) {
	w.tray = t
}

// startClickHandlers sets up goroutines to listen for button clicks
func (w *Window) startClickHandlers() {
	go w.handleStartClick()
	go w.handlePauseClick()
	go w.handleResetClick()
	go w.handleSkipClick()
	go w.handleQuitClick()
}

// handleStartClick handles Start/Resume button clicks
func (w *Window) handleStartClick() {
	for range w.btnStart.ClickedCh {
		if w.timer == nil {
			continue
		}

		state := w.timer.GetState()
		if state == timer.StateIdle {
			w.timer.Start(w.session.CurrentType, w.session.GetDuration())
		} else if state == timer.StatePaused {
			w.timer.Resume()
			// Update button states after resume
			w.UpdateButtonStates(timer.StateRunning)
			// Restore session-specific emoji based on current session type
			switch w.session.CurrentType {
			case session.TypeWork:
				w.header.SetTitle("🍅 Work Session")
			case session.TypeShortBreak:
				w.header.SetTitle("☕ Short Break")
			case session.TypeLongBreak:
				w.header.SetTitle("🌟 Long Break")
			default:
				w.header.SetTitle("Running")
			}
		}
	}
}

// handlePauseClick handles Pause button clicks
func (w *Window) handlePauseClick() {
	for range w.btnPause.ClickedCh {
		if w.timer == nil {
			continue
		}

		w.timer.Pause()
		// Update button states after pause
		w.UpdateButtonStates(timer.StatePaused)
		w.header.SetTitle("⏸️ Paused")
	}
}

// handleResetClick handles Reset button clicks
func (w *Window) handleResetClick() {
	for range w.btnReset.ClickedCh {
		if w.timer == nil {
			continue
		}

		w.timer.Reset()
		// Reset cycle state
		w.session.Reset()
		// Reset session tracking
		w.sessionDuration = 0
		// Reset progress bar
		w.progressBar.SetTitle("○○○○○○○○○○")
		// Update button states and display after reset
		w.UpdateButtonStates(timer.StateIdle)
		w.timerDisplay.SetTitle("25min")
		w.header.SetTitle("Ready")
		// Update cycle indicator
		w.cycleIndicator.SetTitle(w.session.FormatCycleIndicator())
	}
}

// handleSkipClick handles Skip button clicks
func (w *Window) handleSkipClick() {
	for range w.btnSkip.ClickedCh {
		if w.timer == nil {
			continue
		}

		// Calculate elapsed time
		currentSessionDuration := w.session.GetDuration()
		elapsed := currentSessionDuration - w.timer.GetRemaining()
		elapsedMinutes := elapsed / 60

		// Save current session type before stopping timer
		currentSessionType := w.session.CurrentType

		// Stop timer
		w.timer.Reset()

		// Reset progress bar
		w.progressBar.SetTitle("○○○○○○○○○○")

		// Log as skipped
		err := storage.LogSession(time.Now(), currentSessionType, "skipped", elapsedMinutes)
		if err != nil {
			log.Printf("[ERROR] Failed to log skipped session: %v", err)
		}

		// Increment work session counter if skipping a work session
		// This makes Skip truly advance the cycle (unlike Reset which restarts it)
		w.session.IncrementCycle()

		// Determine next session and advance the cycle
		nextSessionType, nextDuration := w.session.DetermineNext()
		w.session.CurrentType = nextSessionType
		w.timerDisplay.SetTitle(formatTime(nextDuration))

		// Set header based on next session type
		switch nextSessionType {
		case session.TypeWork:
			w.header.SetTitle("🍅 Ready for Work")
		case session.TypeShortBreak:
			w.header.SetTitle("☕ Ready for Break")
		case session.TypeLongBreak:
			w.header.SetTitle("🌟 Ready for Long Break")
		default:
			w.header.SetTitle("Ready")
		}

		// Update button states
		w.UpdateButtonStates(timer.StateIdle)

		// Update cycle indicator
		w.cycleIndicator.SetTitle(w.session.FormatCycleIndicator())
	}
}

// handleQuitClick handles Quit button clicks
func (w *Window) handleQuitClick() {
	for range w.btnQuit.ClickedCh {
		log.Println("[INFO] Quit button clicked")
		systray.Quit()
	}
}

// handleTimerStarted is called when the timer starts
func (w *Window) handleTimerStarted(sessionType string, durationSeconds int) {
	log.Println("[INFO] Timer started")

	// Store session metadata for progress tracking
	w.sessionStartTime = time.Now()
	w.sessionDuration = durationSeconds

	// Update current session type
	w.session.CurrentType = sessionType

	// Set header based on session type
	switch sessionType {
	case session.TypeWork:
		w.header.SetTitle("🍅 Work Session")
	case session.TypeShortBreak:
		w.header.SetTitle("☕ Short Break")
	case session.TypeLongBreak:
		w.header.SetTitle("🌟 Long Break")
	default:
		w.header.SetTitle("Running")
	}

	// Update button states
	w.UpdateButtonStates(timer.StateRunning)

	// Update cycle indicator
	w.cycleIndicator.SetTitle(w.session.FormatCycleIndicator())

	// Log session started
	err := storage.LogSession(time.Now(), sessionType, "started", 0)
	if err != nil {
		log.Fatalf("[ERROR] Failed to log session: %v", err)
	}
}

// handleTimerTick is called on each timer tick
func (w *Window) handleTimerTick(remaining int) {
	timeStr := formatTime(remaining)
	w.timerDisplay.SetTitle(timeStr)

	// Update progress bar
	elapsed := w.sessionDuration - remaining
	progressStr := formatProgressBar(elapsed, w.sessionDuration)
	w.progressBar.SetTitle(progressStr)
}

// handleTimerCompleted is called when the timer reaches zero
func (w *Window) handleTimerCompleted() {
	log.Println("[INFO] Timer completed")

	// Reset session tracking
	w.sessionDuration = 0

	// Increment work session counter if completing a work session
	w.session.IncrementCycle()

	// Determine next session type
	nextSessionType, nextDuration := w.session.DetermineNext()
	w.session.CurrentType = nextSessionType

	// Update display for next session
	w.timerDisplay.SetTitle(formatTime(nextDuration))

	// Set header based on next session type
	switch nextSessionType {
	case session.TypeWork:
		w.header.SetTitle("🍅 Ready for Work")
	case session.TypeShortBreak:
		w.header.SetTitle("☕ Ready for Break")
	case session.TypeLongBreak:
		w.header.SetTitle("🌟 Ready for Long Break")
	default:
		w.header.SetTitle("Ready")
	}

	// Update button states
	w.UpdateButtonStates(timer.StateIdle)

	// Update cycle indicator
	w.cycleIndicator.SetTitle(w.session.FormatCycleIndicator())

	// Log session completed (using current session type for now, will be improved in Step 11)
	err := storage.LogSession(time.Now(), w.session.CurrentType, "completed", 25)
	if err != nil {
		log.Fatalf("[ERROR] Failed to log session: %v", err)
	}
}

// UpdateButtonStates updates button enabled/disabled state based on timer state
func (w *Window) UpdateButtonStates(state timer.State) {
	if shouldStartBeEnabled(state) {
		w.btnStart.Enable()
	} else {
		w.btnStart.Disable()
	}

	if shouldPauseBeEnabled(state) {
		w.btnPause.Enable()
	} else {
		w.btnPause.Disable()
	}

	if shouldResetBeEnabled(state) {
		w.btnReset.Enable()
	} else {
		w.btnReset.Disable()
	}

	// Update tray icon to match timer state
	w.updateTrayIcon(state)
}

// updateTrayIcon updates the tray icon based on current session type and timer state
func (w *Window) updateTrayIcon(state timer.State) {
	if w.tray == nil {
		return
	}

	sessionType := ""
	if state == timer.StateRunning || state == timer.StatePaused {
		sessionType = w.session.CurrentType
	}

	w.tray.UpdateIcon(sessionType, state)
}

// formatTime converts seconds to minutes-only format
func formatTime(seconds int) string {
	minutes := seconds / 60
	return fmt.Sprintf("%dmin", minutes)
}

// formatProgressBar generates a 10-segment progress bar using Unicode circles
func formatProgressBar(elapsed, duration int) string {
	if duration == 0 {
		return "○○○○○○○○○○"
	}
	fillPercentage := float64(elapsed) / float64(duration)
	filledSegments := int(fillPercentage * 10)
	if filledSegments > 10 {
		filledSegments = 10
	}

	result := ""
	for i := 0; i < filledSegments; i++ {
		result += "●"
	}
	for i := filledSegments; i < 10; i++ {
		result += "○"
	}
	return result
}

// shouldStartBeEnabled returns true if Start button should be enabled
func shouldStartBeEnabled(state timer.State) bool {
	return state == timer.StateIdle || state == timer.StatePaused
}

// shouldPauseBeEnabled returns true if Pause button should be enabled
func shouldPauseBeEnabled(state timer.State) bool {
	return state == timer.StateRunning
}

// shouldResetBeEnabled returns true if Reset button should be enabled
func shouldResetBeEnabled(state timer.State) bool {
	return state == timer.StateRunning || state == timer.StatePaused
}
