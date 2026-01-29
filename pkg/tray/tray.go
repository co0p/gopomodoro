package tray

import (
	gopomodoro "github.com/co0p/gopomodoro/pkg"
	"github.com/getlantern/systray"
)

// Tray implements the system tray using getlantern/systray.
type Tray struct {
	cycle *gopomodoro.Cycle
}

// New creates a new Tray with the given cycle.
func New(c *gopomodoro.Cycle) *Tray {
	return &Tray{cycle: c}
}

// OnStateChanged updates the tray display when the cycle state changes.
func (t *Tray) OnStateChanged(state gopomodoro.CycleState) {
	formatter := &Formatter{}
	systray.SetTitle(formatter.Format(state, t.cycle.Remaining()))
}

// Run starts the systray. Blocks until quit.
func (t *Tray) Run() error {
	systray.Run(t.onReady, t.onExit)
	return nil
}

func (t *Tray) onReady() {
	systray.SetTitle("🍅")
	systray.SetTooltip("GoPomodoro")

	mStart := systray.AddMenuItem("Start", "Start Pomodoro")
	mStop := systray.AddMenuItem("Stop", "Stop Pomodoro")
	systray.AddSeparator()
	mQuit := systray.AddMenuItem("Quit", "Quit GoPomodoro")

	go func() {
		for {
			select {
			case <-mStart.ClickedCh:
				t.cycle.Start()
			case <-mStop.ClickedCh:
				t.cycle.Stop()
			case <-mQuit.ClickedCh:
				systray.Quit()
				return
			}
		}
	}()
}

func (t *Tray) onExit() {
	// Cleanup if needed
}
