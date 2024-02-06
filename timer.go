package main

import (
  "fmt"
  "math"
  "time"
)


type Timer struct {
  duration uint8
  tickHz int
}

func CreateTimer() *Timer {
  t := Timer{duration: math.MaxUint8, tickHz: 60}
  return &t
}

func (t *Timer) Reset () {
  t.duration = math.MaxUint8

}


func (t *Timer) Start() {
  go func() {
    tick := time.Duration(1000 / t.tickHz)
    for t.duration > uint8(0) {
      t.duration -= 1
      time.Sleep(tick * time.Millisecond)
      fmt.Println(t.duration)
    }
  }()
}
