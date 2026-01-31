package main

import (
	"flag"
	"log"

	gopomodoro "github.com/co0p/gopomodoro/pkg"
	"github.com/co0p/gopomodoro/pkg/sound"
	"github.com/co0p/gopomodoro/pkg/ticker"
	"github.com/co0p/gopomodoro/pkg/tray"
)

func main() {
	silent := flag.Bool("silent", false, "disable sound notifications")
	flag.Parse()

	t := ticker.New()

	var notifier gopomodoro.Notifier
	if !*silent {
		notifier = sound.NewNotifier()
	}

	c := &gopomodoro.Cycle{
		Ticker:   t,
		Notifier: notifier,
	}
	tr := tray.New(c)
	c.Observer = tr

	if err := tr.Run(); err != nil {
		log.Fatal(err)
	}
}
