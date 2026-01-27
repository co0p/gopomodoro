package app

import "testing"

// mockTray is a test double for ports.TrayUI
type mockTray struct {
	runCalled bool
}

func (m *mockTray) Run() error {
	m.runCalled = true
	return nil
}

func TestAppStartsWithoutError(t *testing.T) {
	tray := &mockTray{}
	application := New(tray)
	if application == nil {
		t.Fatal("expected app to be created, got nil")
	}
}

func TestAppRunCallsTrayRun(t *testing.T) {
	tray := &mockTray{}
	application := New(tray)

	err := application.Run()

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if !tray.runCalled {
		t.Fatal("expected tray.Run() to be called")
	}
}
