package timer

import (
	"fmt"
	"os"
	"sync/atomic"
	"time"

	"github.com/gopxl/beep"
	"github.com/gopxl/beep/mp3"
	"github.com/gopxl/beep/speaker"
)

func Timer() *atomic.Int64 {
	var timer atomic.Int64

	go func() {
		tick := time.NewTicker(17 * time.Millisecond).C
		for {
			if timer.Load() > 0 {
				timer.Add(-1)
			}
			<-tick
		}
	}()
	return &timer
}

// TODO wrap sound lib in separate component
func SoundTimer() *atomic.Int64 {
	t := Timer()
	f, err := os.Open("./assets/beep2.mp3")
	if err != nil {
		fmt.Println("Failed to open beep sound")
		panic(1)
	}
	streamer, format, err := mp3.Decode(f)
	if err != nil {
		fmt.Println("Failed to decode beep file")
		panic(1)
	}

	speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))

	go func() {
		defer streamer.Close()
		ctrl := &beep.Ctrl{Streamer: beep.Loop(-1, streamer), Paused: false}
		speaker.Play(ctrl)
		for {
			if t.Load() > 0 {
				speaker.Lock()
				ctrl.Paused = false
				speaker.Unlock()
			} else {
				speaker.Lock()
				ctrl.Paused = true
				speaker.Unlock()
			}
		}
	}()
	return t
}
