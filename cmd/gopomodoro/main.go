package main

import (
	"flag"
	"log"

	"github.com/benbjohnson/clock"
	"github.com/co0p/gopomodoro/internal/adapters/ui"
	"github.com/co0p/gopomodoro/internal/pomodoro"
	"github.com/co0p/gopomodoro/internal/storage"
	"github.com/co0p/gopomodoro/internal/tray"
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
	systray.SetTitle("🍅 25m")
	log.Println("[INFO] Tray icon initialized successfully")

	// Create infrastructure adapters
	clk := clock.New()
	fileStorage := storage.NewFileStorage()
	log.Println("[INFO] Infrastructure adapters created")

	// Create core service
	service := pomodoro.NewService(clk, fileStorage, nil)
	log.Println("[INFO] Pomodoro service created")

	// Create tray instance for UI updates
	trayInstance := tray.New()
	log.Println("[INFO] Tray instance created")

	// Create UI adapter
	uiAdapter := ui.NewSystrayAdapter(service, trayInstance)
	log.Println("[INFO] UI adapter created")

	// Wire UI as notifier
	service.SetNotifier(uiAdapter)
	log.Println("[INFO] UI adapter wired as notifier")

	// Initialize menu
	uiAdapter.InitializeMenu()
	log.Println("[INFO] Systray menu initialized")

	if *smokeTest {
		log.Println("[INFO] Smoke test - quitting after initialization")
		systray.Quit()
	}
}

func onExit() {
	log.Println("[INFO] Application shutting down")
}
