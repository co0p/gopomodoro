package sound

import (
	"sync"
	"time"

	"github.com/faiface/beep"
	"github.com/faiface/beep/generators"
	"github.com/faiface/beep/speaker"
)

const sampleRate = beep.SampleRate(48000)

var (
	speakerInit sync.Once
	initErr     error
)

type Notifier struct {
}

func NewNotifier() *Notifier {
	// Initialize speaker once for the entire process
	speakerInit.Do(func() {
		initErr = speaker.Init(sampleRate, sampleRate.N(time.Second/10))
	})
	return &Notifier{}
}

func (n *Notifier) Notify() {
	// Non-blocking: play sound in goroutine
	go n.playSound()
}

func (n *Notifier) playSound() {
	// If speaker init failed, graceful degradation
	if initErr != nil {
		return
	}

	// Generate a simple beep tone (440 Hz for 200ms)
	tone, err := generators.SinTone(sampleRate, 440)
	if err != nil {
		return
	}
	duration := sampleRate.N(200 * time.Millisecond)
	sound := beep.Take(duration, tone)

	// Play the sound
	done := make(chan bool)
	speaker.Play(beep.Seq(sound, beep.Callback(func() {
		done <- true
	})))

	// Wait for sound to finish (in goroutine, so doesn't block Notify())
	<-done
}
