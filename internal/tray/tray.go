package tray

import (
"log"
"os"
"path/filepath"

"github.com/getlantern/systray"
)

var clickHandler func()

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
