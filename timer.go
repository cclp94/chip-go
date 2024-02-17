package main

import (
	"fmt"
	"os"
	"sync/atomic"
	"time"

	"github.com/gopxl/beep"
	"github.com/gopxl/beep/mp3"
	"github.com/gopxl/beep/speaker"
)

func timer() *atomic.Int64 {
	var timer atomic.Int64

	go func() {
		tick := time.Tick(17 * time.Millisecond)
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
func soundTimer() *atomic.Int64 {
	t := timer()
	f, err := os.Open("./assets/beep.mp3")
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
		// done := make(chan bool)
		defer streamer.Close()
		ctrl := &beep.Ctrl{Streamer: beep.Loop(-1, streamer), Paused: true}
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
