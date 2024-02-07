package main

import (
	"sync/atomic"
	"time"
)


func timer() *atomic.Int64  {
  var timer atomic.Int64 
  
  go func() {
    tick := time.Tick(17 * time.Millisecond)
    for {
      if timer.Load() > 0 {
        timer.Add(-1)
      }
      <- tick
    }
  }()
  return &timer
}

