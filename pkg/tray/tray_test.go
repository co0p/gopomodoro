package tray_test

import (
	"testing"

	gopomodoro "github.com/co0p/gopomodoro/pkg"
	"github.com/co0p/gopomodoro/pkg/tray"
)

func TestTrayImplementsCycleObserver(t *testing.T) {
	c := &gopomodoro.Cycle{}
	tr := tray.New(c)

	// Verify Tray implements CycleObserver interface
	var _ gopomodoro.CycleObserver = tr

	if tr == nil {
		t.Error("expected non-nil tray")
	}
}
