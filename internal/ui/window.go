// Package ui manages the dropdown window/panel
package ui

import (
	"fmt"
	"log"
	"time"

	"github.com/co0p/gopomodoro/internal/storage"
	"github.com/co0p/gopomodoro/internal/timer"
	"github.com/getlantern/systray"
)

const (
	sessionType     = "work"
	sessionDuration = 25
)

// Window represents the dropdown panel
type Window struct {
	visible      bool
	header       *systray.MenuItem
	timerDisplay *systray.MenuItem
	btnStart     *systray.MenuItem
	btnPause     *systray.MenuItem
	btnReset     *systray.MenuItem
	btnQuit      *systray.MenuItem
	timer        *timer.Timer
}

// CreateWindow initializes the dropdown window with placeholder UI
func CreateWindow() (*Window, error) {
	return &Window{visible: false}, nil
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

	w.timerDisplay = systray.AddMenuItem("25:00", "Timer display")
	w.timerDisplay.Disable()

	systray.AddSeparator()

	w.btnStart = systray.AddMenuItem("Start", "Start timer")
	w.btnPause = systray.AddMenuItem("Pause", "Pause timer")
	w.btnReset = systray.AddMenuItem("Reset", "Reset timer")

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

// startClickHandlers sets up goroutines to listen for button clicks
func (w *Window) startClickHandlers() {
	go w.handleStartClick()
	go w.handlePauseClick()
	go w.handleResetClick()
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
			w.timer.Start()
		} else if state == timer.StatePaused {
			w.timer.Resume()
			// Update button states after resume
			w.UpdateButtonStates(timer.StateRunning)
			w.header.SetTitle("Running")
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
		w.header.SetTitle("Paused")
	}
}

// handleResetClick handles Reset button clicks
func (w *Window) handleResetClick() {
	for range w.btnReset.ClickedCh {
		if w.timer == nil {
			continue
		}

		w.timer.Reset()
		// Update button states and display after reset
		w.UpdateButtonStates(timer.StateIdle)
		w.timerDisplay.SetTitle("25:00")
		w.header.SetTitle("Ready")
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
func (w *Window) handleTimerStarted() {
	log.Println("[INFO] Timer started")
	w.header.SetTitle("Running")

	// Update button states
	w.UpdateButtonStates(timer.StateRunning)

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
}

// handleTimerCompleted is called when the timer reaches zero
func (w *Window) handleTimerCompleted() {
	log.Println("[INFO] Timer completed")
	w.header.SetTitle("Ready")
	w.timerDisplay.SetTitle("25:00")

	// Update button states
	w.UpdateButtonStates(timer.StateIdle)

	// Log session completed
	err := storage.LogSession(time.Now(), sessionType, "completed", sessionDuration)
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
}

// formatTime converts seconds to MM:SS format
func formatTime(seconds int) string {
	minutes := seconds / 60
	secs := seconds % 60
	return fmt.Sprintf("%02d:%02d", minutes, secs)
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
