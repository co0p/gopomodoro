package gopomodoro

import "time"

// CycleState represents the state of the pomodoro cycle.
// When used as duration, the value represents minutes.
type CycleState int

const (
	Idle       CycleState = 0
	ShortBreak CycleState = 5
	LongBreak  CycleState = 15
	Pomodoro   CycleState = 25
)

// Ticker provides time ticks for the pomodoro countdown.
type Ticker interface {
	Start()
	Stop()
	OnTick() <-chan struct{}
}

// CycleObserver receives notifications of cycle state changes.
type CycleObserver interface {
	OnStateChanged(state CycleState)
}

type Cycle struct {
	State    CycleState
	TimeLeft time.Duration
	Ticker   Ticker
	Observer CycleObserver

	// pomodoroCount tracks completed pomodoros to determine break type.
	// Increments when a pomodoro completes, persists across short breaks,
	// resets to 0 when Stop() is called or after a long break completes.
	pomodoroCount int
}

func (c *Cycle) Is(s CycleState) bool {
	return c.State == s
}

func (c *Cycle) notifyStateChanged() {
	if c.Observer != nil {
		c.Observer.OnStateChanged(c.State)
	}
}

func (c *Cycle) Start() {
	if c.State == Idle {
		c.State = Pomodoro
		c.TimeLeft = time.Duration(Pomodoro) * time.Minute
		c.notifyStateChanged()
		if c.Ticker != nil {
			c.Ticker.Start()
			go func() {
				for range c.Ticker.OnTick() {
					c.AdvanceMinute()
				}
			}()
		}
	}
}

func (c *Cycle) Stop() {
	c.State = Idle
	c.TimeLeft = 0
	c.pomodoroCount = 0
	c.notifyStateChanged()
	if c.Ticker != nil {
		c.Ticker.Stop()
	}
}

func (c *Cycle) Remaining() time.Duration {
	return c.TimeLeft
}

// AdvanceMinute decrements the timer by one minute and may transition state.
func (c *Cycle) AdvanceMinute() {
	switch c.State {
	case Pomodoro:
		c.advancePomodoro()
	case ShortBreak:
		c.advanceShortBreak()
	case LongBreak:
		c.advanceLongBreak()
		return
	}
	c.notifyStateChanged()
}

func (c *Cycle) advancePomodoro() {
	c.TimeLeft -= time.Minute
	if c.TimeLeft <= 0 {
		c.pomodoroCount++
		if c.pomodoroCount >= 4 {
			c.State = LongBreak
			c.TimeLeft = time.Duration(LongBreak) * time.Minute
		} else {
			c.State = ShortBreak
			c.TimeLeft = time.Duration(ShortBreak) * time.Minute
		}
	}
}

func (c *Cycle) advanceShortBreak() {
	c.TimeLeft -= time.Minute
	if c.TimeLeft <= 0 {
		c.State = Pomodoro
		c.TimeLeft = time.Duration(Pomodoro) * time.Minute
	}
}

func (c *Cycle) advanceLongBreak() {
	c.TimeLeft -= time.Minute
	if c.TimeLeft <= 0 {
		c.Stop()
	}
}
