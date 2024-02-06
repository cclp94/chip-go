package main

import (
  "fmt"
  "os"
  "time"
)


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


func chip8(memory []byte) (func(), func()) {

  // var display [64][32]byte
  var pc int = 0x200
  //var l *uint8
  //var stack []uint16
  //var delayTimer *Timer = CreateTimer()
  //var soundTimer *Timer = CreateTimer() 
  // var registers [16]uint8

  fetchInstruction := func () uint16 {
    instruction := uint16(memory[pc])<<8 | uint16(memory[pc+1])
    fmt.Printf("fetch instruction =0x%04X\n", instruction)
    pc += 2
    return instruction
  }

  decodeInstruction := func () {
    // Binary mask first hex nibble
    switch instruction & 0xf000 {
    case 0x0:
      println("decoded: 0")
      switch instruction & 0x00ff {
      case 0xE0:
        // TODO Clear screen
      case 0xEE:
        // TODO pop from stack and set PC
      }
    case 0x1:
      println("decoded: 1")
      addr := instruction & 0x0fff
      // TODO Set pc to addr
    case 0x2:
      println("decoded: 2")
      addr := instruction & 0x0fff
      // TODO Push PC to stack

      // TODO Set PC to addr
    case 0x3:
      println("decoded: 3")
    case 0x4:
      println("decoded: 4")
    case 0x5:
      println("decoded: 5")
    case 0x6:
      println("decoded: 6")
    case 0x7:
      println("decoded: 7")
    case 0x8:
      println("decoded: 8")
    case 0x9:
      println("decoded: 9")
    case 0xa:
      println("decoded: a")
    case 0xb:
      println("decoded: b")
    case 0xc:
      println("decoded: c")
    case 0xd:
      println("decoded: d")
    case 0xe:
      println("decoded: e")
    case 0xf:
      println("decoded: f")
    }
  }

  return fetchInstruction, decodeIsntruction

}
func main() {
  romArg := os.Args[1]

  var memory [4096]byte


  // Font runs from addr 050 to 09F
  registerFont(memory[0x50:])
  registerRom(romArg, memory[0x200:])
  tick, _ := time.ParseDuration("1.5ms")

  for pc < len(memory) - 2 { 
    instruction := fetchInstruction(&pc, memory[:])
    decode(instruction)
    pc += 2
    time.Sleep(tick) 
  }
}

