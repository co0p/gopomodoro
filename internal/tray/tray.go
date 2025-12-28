package tray

import (
	_ "embed"
	"fmt"
	"log"

	"github.com/co0p/gopomodoro/internal/timer"
	"github.com/getlantern/systray"
)

var clickHandler func()

//go:embed assets/icon-idle.png
var iconIdle []byte

//go:embed assets/icon-work.png
var iconWork []byte

//go:embed assets/icon-short-break.png
var iconShortBreak []byte

//go:embed assets/icon-long-break.png
var iconLongBreak []byte

//go:embed assets/icon-paused.png
var iconPaused []byte

// Tray manages the system tray icon
type Tray struct{}

// New creates a new Tray instance
func New() *Tray {
	return &Tray{}
}

// UpdateIcon updates the tray icon based on session type and timer state
func (t *Tray) UpdateIcon(sessionType string, state timer.State) {
	iconData := t.getIconData(sessionType, state)
	systray.SetIcon(iconData)
}

// getIconData determines which icon bytes to use based on session type and state
func (t *Tray) getIconData(sessionType string, state timer.State) []byte {
	if state == timer.StateRunning {
		switch sessionType {
		case "work":
			return iconWork
		case "short_break":
			return iconShortBreak
		case "long_break":
			return iconLongBreak
		}
	}

	if state == timer.StatePaused {
		return iconPaused
	}

	// Default to tomato icon (work icon) for idle state
	return iconWork
}

// GetEmojiForState returns the appropriate emoji for the current session state
func (t *Tray) GetEmojiForState(sessionType string, state timer.State) string {
	// Paused state always shows pause emoji
	if state == timer.StatePaused {
		return "⏸️"
	}

	// Idle or running states - select emoji based on session type
	switch sessionType {
	case "short_break":
		return "☕"
	case "long_break":
		return "🌟"
	default:
		// Default to tomato (for work and idle)
		return "🍅"
	}
}

// FormatMinutes converts seconds to a formatted minute string
func (t *Tray) FormatMinutes(seconds int) string {
	minutes := seconds / 60
	return fmt.Sprintf("%dm", minutes)
}

// UpdateDisplay updates both the tray icon and title text with session info
func (t *Tray) UpdateDisplay(sessionType string, state timer.State, remainingSeconds int) {
	// Update icon
	iconData := t.getIconData(sessionType, state)
	systray.SetIcon(iconData)

	// Update title text with emoji and time
	emoji := t.GetEmojiForState(sessionType, state)
	timeStr := t.FormatMinutes(remainingSeconds)
	title := emoji + " " + timeStr
	systray.SetTitle(title)
}

// Initialize sets up the system tray icon
func Initialize() error {
	log.Println("[INFO] Tray initialization called")
	return nil
}

// SetIcon updates the tray icon image
func SetIcon(iconData []byte) error {
	systray.SetIcon(iconData)
	return nil
}

// LoadIconFromAssets loads the default icon from embedded assets
func LoadIconFromAssets() ([]byte, error) {
	return iconIdle, nil
}

// OnClick registers a callback for tray icon clicks
func OnClick(handler func()) {
	clickHandler = handler
}
