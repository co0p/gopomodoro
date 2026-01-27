package ticker

import (
	"time"
)

// Ticker implements gopomodoro.Ticker using time.Ticker.
type Ticker struct {
	ticker   *time.Ticker
	tickChan chan struct{}
	stopChan chan struct{}
}

func New() *Ticker {
	return &Ticker{
		tickChan: make(chan struct{}),
		stopChan: make(chan struct{}),
	}
}

func (t *Ticker) Start() {
	t.ticker = time.NewTicker(1 * time.Minute)
	go func() {
		for {
			select {
			case <-t.ticker.C:
				t.tickChan <- struct{}{}
			case <-t.stopChan:
				return
			}
		}
	}()
}

func (t *Ticker) Stop() {
	if t.ticker != nil {
		t.ticker.Stop()
		t.stopChan <- struct{}{}
	}
}

func (t *Ticker) OnTick() <-chan struct{} {
	return t.tickChan
}
