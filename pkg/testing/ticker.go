package testing

import gopomodoro "github.com/co0p/gopomodoro/pkg"

type MockTicker struct {
	tickChan chan struct{}
	started  bool
}

func NewMockTicker() *MockTicker {
	return &MockTicker{
		tickChan: make(chan struct{}),
	}
}

func (m *MockTicker) Start() {
	m.started = true
}

func (m *MockTicker) Stop() {
	m.started = false
}

func (m *MockTicker) OnTick() <-chan struct{} {
	return m.tickChan
}

func (m *MockTicker) Fire() {
	m.tickChan <- struct{}{}
}

var _ gopomodoro.Ticker = (*MockTicker)(nil)
