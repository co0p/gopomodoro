package tray

import (
	"fmt"
	"time"

	gopomodoro "github.com/co0p/gopomodoro/pkg"
)

type Formatter struct{}

func (f *Formatter) Format(state gopomodoro.CycleState, remaining time.Duration) string {
	const tomatoIcon = "🍅"

	switch state {
	case gopomodoro.Idle:
		return tomatoIcon
	case gopomodoro.Pomodoro:
		minutes := int(remaining.Minutes())
		return tomatoIcon + " " + formatMinutes(minutes) + "m"
	default:
		return tomatoIcon
	}
}

func formatMinutes(minutes int) string {
	return fmt.Sprintf("%02d", minutes)
}
