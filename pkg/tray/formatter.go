package tray

import (
	"fmt"
	"time"

	gopomodoro "github.com/co0p/gopomodoro/pkg"
)

type Formatter struct{}

func (f *Formatter) Format(state gopomodoro.CycleState, remaining time.Duration) string {
	const tomatoIcon = "🍅"
	const coffeeIcon = "☕"
	const longBreakIcon = "🌴"

	minutes := fmt.Sprintf("%d", int(remaining.Minutes()))

	switch state {
	case gopomodoro.Pomodoro:
		return tomatoIcon + " " + minutes + "m"
	case gopomodoro.ShortBreak:
		return coffeeIcon + " " + minutes + "m"
	case gopomodoro.LongBreak:
		return longBreakIcon + " " + minutes + "m"
	case gopomodoro.Idle:
		fallthrough
	default:
		return tomatoIcon
	}
}
