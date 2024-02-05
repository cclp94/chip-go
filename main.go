package main

import (
  "fmt"
  "math"
  "os"
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


func registerRom(romPath string, memory []byte) {
  file, err := os.ReadFile(romPath)
  if err != nil {
    fmt.Println("File:", romPath, "could not be read:", err.Error())
    panic(1)
  }
  for i, b := range file {
    memory[i] = b
  }
}

func registerFont(memory []byte) {
  font := []byte{
    0xF0, 0x90, 0x90, 0x90, 0xF0, // 0
    0x20, 0x60, 0x20, 0x20, 0x70, // 1
    0xF0, 0x10, 0xF0, 0x80, 0xF0, // 2
    0xF0, 0x10, 0xF0, 0x10, 0xF0, // 3
    0x90, 0x90, 0xF0, 0x10, 0x10, // 4
    0xF0, 0x80, 0xF0, 0x10, 0xF0, // 5
    0xF0, 0x80, 0xF0, 0x90, 0xF0, // 6
    0xF0, 0x10, 0x20, 0x40, 0x40, // 7
    0xF0, 0x90, 0xF0, 0x90, 0xF0, // 8
    0xF0, 0x90, 0xF0, 0x10, 0xF0, // 9
    0xF0, 0x90, 0xF0, 0x90, 0x90, // A
    0xE0, 0x90, 0xE0, 0x90, 0xE0, // B
    0xF0, 0x80, 0x80, 0x80, 0xF0, // C
    0xE0, 0x90, 0x90, 0x90, 0xE0, // D
    0xF0, 0x80, 0xF0, 0x80, 0xF0, // E
    0xF0, 0x80, 0xF0, 0x80, 0x80,  // F
  }

  for i, f := range font {
    memory[i] = f
  }
}



func main() {
  romArg := os.Args[1]

  var memory [4096]byte
  // var display [64][32]byte
  var pc int = 0x200
  //var l *uint8
  //var stack []uint16
  //var delayTimer *Timer = CreateTimer()
  //var soundTimer *Timer = CreateTimer() 
  // var registers [16]uint8

  // Font runs from addr 050 to 09F
  registerFont(memory[0x50:])
  registerRom(romArg, memory[0x200:])
  tick, _ := time.ParseDuration("1.5ms")
  
  for pc < len(memory) - 2 { 
    instruction := uint16(memory[pc])<<8 | uint16(memory[pc+1])
    fmt.Printf("n =0x%04X = %d\n", instruction, instruction)
    pc += 2
    time.Sleep(tick) 
  }
}

