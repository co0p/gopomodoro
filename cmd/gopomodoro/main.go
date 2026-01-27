package main

import (
	"log"

	"github.com/co0p/gopomodoro/internal/adapters/systray"
	"github.com/co0p/gopomodoro/internal/app"
)

func main() {
	tray := systray.New()
	application := app.New(tray)
	if err := application.Run(); err != nil {
		log.Fatal(err)
	}
}
