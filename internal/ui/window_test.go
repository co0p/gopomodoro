package ui

import "testing"

func TestCreateWindow(t *testing.T) {
	win, err := CreateWindow()
	if err != nil {
		t.Fatalf("CreateWindow() failed: %v", err)
	}
	if win == nil {
		t.Fatal("CreateWindow() returned nil window")
	}
}

func TestShowHide(t *testing.T) {
	win, _ := CreateWindow()

	err := win.Show(100, 100)
	if err != nil {
		t.Fatalf("Show() failed: %v", err)
	}
	if !win.IsVisible() {
		t.Error("Window should be visible after Show()")
	}

	err = win.Hide()
	if err != nil {
		t.Fatalf("Hide() failed: %v", err)
	}
	if win.IsVisible() {
		t.Error("Window should not be visible after Hide()")
	}
}
