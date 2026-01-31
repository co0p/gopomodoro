package sound_test

import (
	"context"
	"testing"
	"time"

	"github.com/co0p/gopomodoro/pkg/sound"
)

func TestSoundNotifier_Notify_DoesNotBlock(t *testing.T) {
	notifier := sound.NewNotifier()

	// Create a context with short timeout
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	done := make(chan struct{})
	go func() {
		notifier.Notify()
		close(done)
	}()

	select {
	case <-done:
		// Success - Notify returned quickly
	case <-ctx.Done():
		t.Fatal("Notify() blocked for more than 100ms")
	}
}

func TestSoundNotifier_Notify_PlaysSound(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping sound playback test in short mode")
	}

	notifier := sound.NewNotifier()

	// Call Notify and wait for sound to complete
	notifier.Notify()

	// Sound is 200ms, wait a bit longer to ensure it completes
	time.Sleep(300 * time.Millisecond)

	// If we get here without panic/crash, sound played successfully
	// (We can't assert that audio was actually heard, but we verify no errors)
}
