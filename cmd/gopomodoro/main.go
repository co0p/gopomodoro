package main

import (
	"flag"
	"log"

	"github.com/co0p/gopomodoro/internal/storage"
	"github.com/co0p/gopomodoro/internal/timer"
	"github.com/co0p/gopomodoro/internal/tray"
	"github.com/co0p/gopomodoro/internal/ui"
	"github.com/getlantern/systray"
)

var smokeTest = flag.Bool("smoke", false, "Run smoke test (start and immediately exit)")

func main() {
	flag.Parse()

	log.Println("[INFO] GoPomodoro starting...")

	// Ensure storage directory exists before starting UI
	if err := storage.EnsureDataDir(); err != nil {
		log.Fatalf("[ERROR] Failed to initialize storage: %v", err)
	}
	log.Println("[INFO] Data directory ensured")

	// Note: To fully suppress dock icon on macOS, this app should be built as a .app bundle
	// with Info.plist containing LSUIElement=true
	// For development builds, dock icon may appear briefly - this is acceptable for lite mode
	// The systray library handles most of the integration automatically

	if *smokeTest {
		log.Println("[INFO] Smoke test mode - running minimal initialization")
		systray.Run(onReady, onExit)
		return
	}

	systray.Run(onReady, onExit)
}

func onReady() {
	log.Println("[INFO] Initializing tray icon...")

	if err := tray.Initialize(); err != nil {
		log.Fatalf("[ERROR] Failed to initialize tray: %v", err)
	}

	// Load and set icon
	iconData, err := tray.LoadIconFromAssets()
	if err != nil {
		log.Fatalf("[ERROR] Failed to load icon: %v", err)
	}

	if err := tray.SetIcon(iconData); err != nil {
		log.Fatalf("[ERROR] Failed to set icon: %v", err)
	}

	systray.SetTooltip("GoPomodoro")
	log.Println("[INFO] Tray icon initialized successfully")

	// Create UI window
	window, err := ui.CreateWindow()
	if err != nil {
		log.Fatalf("[ERROR] Failed to create window: %v", err)
	}

	// Initialize menu items (systray menu approach)
	window.InitializeMenu()
	log.Println("[INFO] Dropdown window created")

	// Create and wire timer
	tmr := timer.New()
	log.Println("[INFO] Timer created")

	// Create tray instance
	trayInstance := tray.New()
	log.Println("[INFO] Tray instance created")

	// Set timer in window (this registers event handlers and starts click handlers)
	window.SetTimer(tmr)
	log.Println("[INFO] Timer wired to UI")

	// Set tray in window for icon updates
	window.SetTray(trayInstance)
	log.Println("[INFO] Tray wired to UI")

	// Update button states to enable Start button
	window.UpdateButtonStates(tmr.GetState())
	log.Println("[INFO] Button states initialized")

	// For systray menu approach, menu is always "ready" to show
	// No explicit click handler needed - systray manages menu display

	if *smokeTest {
		log.Println("[INFO] Smoke test - quitting after initialization")
		systray.Quit()
	}
}

func onExit() {
	log.Println("[INFO] Application shutting down")
}
