package main

import (
	"log"

	gopomodoro "github.com/co0p/gopomodoro/pkg"
	"github.com/co0p/gopomodoro/pkg/ticker"
	"github.com/co0p/gopomodoro/pkg/tray"
)

func main() {
	t := ticker.New()
	c := &gopomodoro.Cycle{Ticker: t}
	tr := tray.New(c)
	c.Observer = tr

	if err := tr.Run(); err != nil {
		log.Fatal(err)
	}
}
