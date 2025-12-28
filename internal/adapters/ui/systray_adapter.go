package ui

import (
	"log"

	"github.com/co0p/gopomodoro/internal/pomodoro"
	"github.com/co0p/gopomodoro/internal/timer"
	"github.com/co0p/gopomodoro/internal/tray"
	"github.com/getlantern/systray"
)

// SystrayAdapter adapts the pomodoro service to the systray UI
// It implements the Notifier interface to receive service events
type SystrayAdapter struct {
	service        *pomodoro.Service
	tray           *tray.Tray
	progressBar    *systray.MenuItem
	cycleIndicator *systray.MenuItem
	btnStart       *systray.MenuItem
	btnReset       *systray.MenuItem
	btnSkip        *systray.MenuItem
	btnQuit        *systray.MenuItem
}

// NewSystrayAdapter creates a new systray adapter
func NewSystrayAdapter(service *pomodoro.Service, trayIcon *tray.Tray) *SystrayAdapter {
	return &SystrayAdapter{
		service: service,
		tray:    trayIcon,
	}
}

// InitializeMenu creates the systray menu items and starts click handlers
func (a *SystrayAdapter) InitializeMenu() {
	log.Println("[INFO] Initializing systray menu items...")

	a.progressBar = systray.AddMenuItem("○○○○○○○○○○", "Session progress")
	a.progressBar.Disable()

	a.cycleIndicator = systray.AddMenuItem("Session 1/4  🍅○○○", "Cycle progress")
	a.cycleIndicator.Disable()

	systray.AddSeparator()

	a.btnStart = systray.AddMenuItem("Start", "Start/Pause/Resume timer")
	a.btnReset = systray.AddMenuItem("Reset", "Reset timer")

	systray.AddSeparator()

	a.btnSkip = systray.AddMenuItem("Skip", "Skip to next session")

	systray.AddSeparator()

	a.btnQuit = systray.AddMenuItem("Quit", "Quit the application")

	log.Println("[INFO] Menu items initialized")

	// Start click handlers
	a.startClickHandlers()
}

// SessionStarted handles session start events from the service
func (a *SystrayAdapter) SessionStarted(sessionType string, duration int) {
	log.Printf("[INFO] Session started: %s (%d seconds)", sessionType, duration)

	// Update cycle indicator
	if a.cycleIndicator != nil {
		a.cycleIndicator.SetTitle(a.service.GetCycleProgress())
	}

	// Reset progress bar to empty
	if a.progressBar != nil {
		a.progressBar.SetTitle("○○○○○○○○○○")
	}

	// Update tray icon
	if a.tray != nil {
		// Convert pomodoro.State to timer.State for tray
		var tState timer.State
		switch a.service.GetState() {
		case pomodoro.StateIdle:
			tState = timer.StateIdle
		case pomodoro.StateRunning:
			tState = timer.StateRunning
		case pomodoro.StatePaused:
			tState = timer.StateIdle // Show as idle when paused
		}
		a.tray.UpdateDisplay(sessionType, tState, duration)
	}
}

// SessionTick handles tick events from the service
func (a *SystrayAdapter) SessionTick(remainingSeconds int) {
	// Update progress bar
	if a.progressBar != nil {
		sessionType := a.service.GetCurrentSessionType()
		a.progressBar.SetTitle(formatProgressBar(remainingSeconds, sessionType))
	}

	// Update tray icon
	if a.tray != nil {
		sessionType := a.service.GetCurrentSessionType()
		a.tray.UpdateDisplay(sessionType, timer.StateRunning, remainingSeconds)
	}
}

// SessionCompleted handles session completion events from the service
func (a *SystrayAdapter) SessionCompleted(sessionType string) {
	log.Printf("[INFO] Session completed: %s", sessionType)

	// Update tray icon to idle
	if a.tray != nil {
		a.tray.UpdateDisplay(a.service.GetCurrentSessionType(), timer.StateIdle, 0)
	}

	// Update cycle indicator
	if a.cycleIndicator != nil {
		a.cycleIndicator.SetTitle(a.service.GetCycleProgress())
	}

	// Reset progress bar
	if a.progressBar != nil {
		a.progressBar.SetTitle("○○○○○○○○○○")
	}
}

// StateChanged handles state change events from the service
func (a *SystrayAdapter) StateChanged(state pomodoro.State) {
	log.Printf("[INFO] State changed: %v", state)

	// Update tray icon
	if a.tray != nil {
		sessionType := a.service.GetCurrentSessionType()
		remaining := a.service.GetRemainingSeconds()
		// Convert pomodoro.State to timer.State
		var tState timer.State
		switch state {
		case pomodoro.StateIdle:
			tState = timer.StateIdle
		case pomodoro.StateRunning:
			tState = timer.StateRunning
		}
		a.tray.UpdateDisplay(sessionType, tState, remaining)
	}

	// Update button states
	a.updateButtonStates(state)
}

// startClickHandlers starts goroutines to handle menu item clicks
func (a *SystrayAdapter) startClickHandlers() {
	go a.handleStartClick()
	go a.handleResetClick()
	go a.handleSkipClick()
	go a.handleQuitClick()
}

// handleStartClick handles Start/Pause/Resume button clicks
func (a *SystrayAdapter) handleStartClick() {
	for range a.btnStart.ClickedCh {
		state := a.service.GetState()
		switch state {
		case pomodoro.StateIdle:
			if err := a.service.StartSession(); err != nil {
				log.Printf("[WARN] Failed to start session: %v", err)
			}
		case pomodoro.StateRunning:
			if err := a.service.PauseSession(); err != nil {
				log.Printf("[WARN] Failed to pause session: %v", err)
			}
		case pomodoro.StatePaused:
			if err := a.service.ResumeSession(); err != nil {
				log.Printf("[WARN] Failed to resume session: %v", err)
			}
		}
	}
}

// handleResetClick handles Reset button clicks
func (a *SystrayAdapter) handleResetClick() {
	for range a.btnReset.ClickedCh {
		log.Println("[INFO] Reset button clicked")

		if err := a.service.ResetCycle(); err != nil {
			log.Printf("[WARN] Failed to reset cycle: %v", err)
		}

		// Reset UI elements
		if a.progressBar != nil {
			a.progressBar.SetTitle("○○○○○○○○○○")
		}
		if a.cycleIndicator != nil {
			a.cycleIndicator.SetTitle(a.service.GetCycleProgress())
		}
		if a.tray != nil {
			a.tray.UpdateDisplay("work", timer.StateIdle, 0)
		}

		a.updateButtonStates(pomodoro.StateIdle)
	}
}

// handleSkipClick handles Skip button clicks
func (a *SystrayAdapter) handleSkipClick() {
	for range a.btnSkip.ClickedCh {
		log.Println("[INFO] Skip button clicked")

		if err := a.service.SkipSession(); err != nil {
			log.Printf("[WARN] Failed to skip session: %v", err)
		}

		// Update tray display with next session type and idle state
		if a.tray != nil {
			nextSessionType := a.service.GetCurrentSessionType()
			nextDuration := a.service.GetCurrentDuration()
			a.tray.UpdateDisplay(nextSessionType, timer.StateIdle, nextDuration)
		}
	}
}

// handleQuitClick handles Quit button clicks
func (a *SystrayAdapter) handleQuitClick() {
	for range a.btnQuit.ClickedCh {
		log.Println("[INFO] Quit requested")
		systray.Quit()
	}
}

// updateButtonStates updates button enabled/disabled states based on timer state
func (a *SystrayAdapter) updateButtonStates(state pomodoro.State) {
	switch state {
	case pomodoro.StateIdle:
		a.btnStart.SetTitle("Start")
		a.btnStart.Enable()
		a.btnReset.Disable()
		a.btnSkip.Disable()

	case pomodoro.StateRunning:
		a.btnStart.SetTitle("Pause")
		a.btnStart.Enable()
		a.btnReset.Enable()
		a.btnSkip.Enable()

	case pomodoro.StatePaused:
		a.btnStart.SetTitle("Resume")
		a.btnStart.Enable()
		a.btnReset.Enable()
		a.btnSkip.Enable()
	}
}

// formatProgressBar creates a visual progress bar based on remaining time
func formatProgressBar(remainingSeconds int, sessionType string) string {
	totalSeconds := 1500 // Default to work session
	if sessionType == "short_break" {
		totalSeconds = 300
	} else if sessionType == "long_break" {
		totalSeconds = 900
	}

	// Calculate progress (0-10 filled circles)
	progress := 10 - (remainingSeconds * 10 / totalSeconds)
	if progress < 0 {
		progress = 0
	}
	if progress > 10 {
		progress = 10
	}

	// Build progress bar
	bar := ""
	for i := 0; i < progress; i++ {
		bar += "●"
	}
	for i := progress; i < 10; i++ {
		bar += "○"
	}

	return bar
}
