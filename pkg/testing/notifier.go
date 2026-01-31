package testing

type MockNotifier struct {
	NotifyCallCount int
}

func (m *MockNotifier) Notify() {
	m.NotifyCallCount++
}
