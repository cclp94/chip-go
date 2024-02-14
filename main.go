package main

import (
	"fmt"
	"os"
	"sync/atomic"

	"github.com/gopxl/pixel/pixelgl"
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


func getNibbleAt(value uint16, at int) uint16 {
 return uint16((value & (0xf000 >> (at * 4))) >> ((3 - at) * 4))
}

func getBitAt(value uint8, at int) uint8 {
  return (value & (0x80 >> at)) >> (7 - at) 
}

func parseArgs(args []string) (string, bool){
  var filename string
  var isLegacy bool

  if (len(args) < 2) {
    fmt.Println("USAGE: <rom-file> [--legacy]")
    panic(1)
  }
  filename = args[1]

  if len(args) == 3 && args[2] == "--legacy" {
    isLegacy = true
  }
  return filename, isLegacy
}

func main() {
  filename, _ := parseArgs(os.Args)

  var memory [4096]byte
  var delayTimer *atomic.Int64 = timer()
  var soundTimer *atomic.Int64 = timer() 

  displayChan := make(chan [][]byte)
  // Font runs from addr 050 to 09F
  registerFont(memory[0x50:])
  registerRom(filename, memory[0x200:])

  fmt.Println(memory)
  go chip8(memory[:], delayTimer, soundTimer, &displayChan, true)
  pixelgl.Run(display(&displayChan))
}

