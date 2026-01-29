package gopomodoro

import "time"

// CycleState represents the state of the pomodoro cycle.
// When used as duration, the value represents minutes.
type CycleState int

const (
	Idle       CycleState = 0
	ShortBreak CycleState = 5
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
					c.Tick()
				}
			}()
		}
	}
}

func (c *Cycle) Stop() {
	c.State = Idle
	c.TimeLeft = 0
	c.notifyStateChanged()
	if c.Ticker != nil {
		c.Ticker.Stop()
	}
}

func (c *Cycle) Remaining() time.Duration {
	return c.TimeLeft
}

func (c *Cycle) Tick() {
	switch c.State {
	case Pomodoro:
		c.TimeLeft -= time.Minute
		if c.TimeLeft <= 0 {
			c.State = ShortBreak
			c.TimeLeft = time.Duration(ShortBreak) * time.Minute
		}
		c.notifyStateChanged()
	case ShortBreak:
		c.TimeLeft -= time.Minute
		if c.TimeLeft <= 0 {
			c.Stop()
			return
		}
		c.notifyStateChanged()
	}
}
