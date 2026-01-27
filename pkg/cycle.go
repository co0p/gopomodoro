package gopomodoro

// CycleState represents the state of the pomodoro cycle.
// When used as duration, the value represents minutes.
type CycleState int

const (
	Idle     CycleState = 0
	Pomodoro CycleState = 25
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
	State     CycleState
	Remaining int
	Ticker    Ticker
	Observer  CycleObserver
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
		c.Remaining = int(Pomodoro)
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
	c.Remaining = 0
	c.notifyStateChanged()
	if c.Ticker != nil {
		c.Ticker.Stop()
	}
}

func (c *Cycle) RemainingMinutes() int {
	return c.Remaining
}

func (c *Cycle) Tick() {
	if c.State != Pomodoro {
		return
	}
	c.Remaining--
	if c.Remaining == 0 {
		c.Stop()
	}
	c.notifyStateChanged()
}
