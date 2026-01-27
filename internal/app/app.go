package app

// TrayUI defines what App needs from a tray implementation.
type TrayUI interface {
	Run() error
}

// App is the application composition root.
type App struct {
	tray TrayUI
}

// New creates a new App instance.
func New(tray TrayUI) *App {
	return &App{tray: tray}
}

// Run starts the application.
func (a *App) Run() error {
	return a.tray.Run()
}
