package tray

import (
	"log"
	"os"
	"path/filepath"

	"github.com/co0p/gopomodoro/internal/timer"
	"github.com/getlantern/systray"
)

var clickHandler func()

// Tray manages the system tray icon
type Tray struct{}

// New creates a new Tray instance
func New() *Tray {
	return &Tray{}
}

// UpdateIcon updates the tray icon based on session type and timer state
func (t *Tray) UpdateIcon(sessionType string, state timer.State) {
	iconPath := t.getIconPath(sessionType, state)
	iconData, err := os.ReadFile(filepath.Join("assets", iconPath))
	if err != nil {
		log.Printf("[ERROR] Failed to load icon %s: %v", iconPath, err)
		return
	}

	systray.SetIcon(iconData)
	log.Printf("[INFO] Tray icon updated to %s", iconPath)
}

// getIconPath determines which icon file to use based on session type and state
func (t *Tray) getIconPath(sessionType string, state timer.State) string {
	if state == timer.StateRunning {
		switch sessionType {
		case "work":
			return "icon-work.png"
		case "short_break":
			return "icon-short-break.png"
		case "long_break":
			return "icon-long-break.png"
		}
	}

	if state == timer.StatePaused {
		return "icon-paused.png"
	}

	// Default to idle icon
	return "icon-idle.png"
}

// Initialize sets up the system tray icon
func Initialize() error {
	log.Println("[INFO] Tray initialization called")
	return nil
}

// SetIcon updates the tray icon image
func SetIcon(iconData []byte) error {
	systray.SetIcon(iconData)
	log.Printf("[INFO] Tray icon set (%d bytes)", len(iconData))
	return nil
}

// LoadIconFromAssets loads the default icon from the assets directory
func LoadIconFromAssets() ([]byte, error) {
	return os.ReadFile(filepath.Join("assets", "icon-idle.png"))
}

// OnClick registers a callback for tray icon clicks
func OnClick(handler func()) {
	clickHandler = handler
}
