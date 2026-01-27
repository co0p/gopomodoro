package testing

import gopomodoro "github.com/co0p/gopomodoro/pkg"

type MockObserver struct {
	StateChanges []gopomodoro.CycleState
}

func (m *MockObserver) OnStateChanged(state gopomodoro.CycleState) {
	m.StateChanges = append(m.StateChanges, state)
}
