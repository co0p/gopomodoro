package testing

import gopomodoro "github.com/co0p/gopomodoro/pkg"

// CompleteCycle advances the timer through the full duration of the current state.
func CompleteCycle(c *gopomodoro.Cycle) {
	duration := int(c.State)
	for i := 0; i < duration; i++ {
		c.AdvanceMinute()
	}
}
