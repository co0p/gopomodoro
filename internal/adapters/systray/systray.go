package systray

import "github.com/getlantern/systray"

// Adapter implements ports.TrayUI using getlantern/systray.
type Adapter struct{}

// New creates a new systray Adapter.
func New() *Adapter {
	return &Adapter{}
}

// Run starts the systray. Blocks until quit.
func (a *Adapter) Run() error {
	systray.Run(a.onReady, a.onExit)
	return nil
}

func (a *Adapter) onReady() {
	systray.SetTitle("🍅")
	systray.SetTooltip("GoPomodoro")

	mQuit := systray.AddMenuItem("Quit", "Quit GoPomodoro")

	go func() {
		<-mQuit.ClickedCh
		systray.Quit()
	}()
}

func (a *Adapter) onExit() {
	// Cleanup if needed
}
