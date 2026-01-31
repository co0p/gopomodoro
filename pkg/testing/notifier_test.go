package testing_test

import (
	"testing"

	pomotest "github.com/co0p/gopomodoro/pkg/testing"
)

func TestMockNotifier_Notify_RecordsCall(t *testing.T) {
	mock := &pomotest.MockNotifier{}

	mock.Notify()

	if mock.NotifyCallCount != 1 {
		t.Errorf("expected NotifyCallCount = 1, got %d", mock.NotifyCallCount)
	}
}
