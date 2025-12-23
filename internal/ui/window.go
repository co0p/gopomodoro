// Package ui manages the dropdown window/panel
package ui

import (
	"log"

	"github.com/getlantern/systray"
)

// Window represents the dropdown panel
type Window struct {
	visible  bool
	header   *systray.MenuItem
	timer    *systray.MenuItem
	btnStart *systray.MenuItem
	btnPause *systray.MenuItem
	btnReset *systray.MenuItem
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

	w.timer = systray.AddMenuItem("25:00", "Timer display")
	w.timer.Disable()

	systray.AddSeparator()

	w.btnStart = systray.AddMenuItem("Start", "Start timer")
	w.btnStart.Disable() // Non-functional for this increment

	w.btnPause = systray.AddMenuItem("Pause", "Pause timer")
	w.btnPause.Disable()

	w.btnReset = systray.AddMenuItem("Reset", "Reset timer")
	w.btnReset.Disable()

	log.Println("[INFO] Menu items initialized")
}
