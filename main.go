package main

import (
	"fmt"
	"os"
	"sync/atomic"
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


func getNibbleAt(value uint16, at int) uint16 {
 return uint16((value & (0xf000 >> (at))) >> ((at + 3) * 4))
}

func chip8(
  memory []byte,
  delayTimer *atomic.Int64,
  soundTimer *atomic.Int64,
  isLegacy bool,
) {

  var pc uint16 = 0x200
  var l uint8
  var stack []uint16

  var V [16]uint8

  fmt.Println(l)

  fetchInstruction := func () uint16 {
    instruction := uint16(memory[pc])<<8 | uint16(memory[pc+1])
    pc += 2
    return instruction
  }

  decodeInstruction := func (instruction uint16) {
    // Binary mask first hex nibble
    switch getNibbleAt(instruction, 0) {
    case 0x0:
      fmt.Println("decoded: 0")
      switch instruction & 0x00ff {
      case 0xE0:
        // TODO Clear screen
      case 0xEE:
        fmt.Println("Exec 00EE")
        toPopAtIndex := len(stack) - 1
        pc = stack[toPopAtIndex]
        stack = stack[:toPopAtIndex]
      }
    case 0x1:
      fmt.Println("Exec 1NNN")
      addr := instruction & 0x0fff
      pc = uint16(addr)
    case 0x2:
      fmt.Println("Exec 2NNN")
      addr := instruction & 0x0fff
      stack = append(stack, uint16(pc))
      pc = uint16(addr)
    case 0x3:
      fmt.Println("Exec 3XNN")
      registerIndex := getNibbleAt(instruction, 1)
      NN := uint8(instruction & 0x00ff)
      if V[registerIndex] == NN {
        pc += 2
      }
    case 0x4:
      fmt.Println("Exec 4XNN")
      registerIndex := getNibbleAt(instruction, 1)
      NN := uint8(instruction & 0x00ff)
      if V[registerIndex] != NN {
        pc += 2
      }
    case 0x5:
      fmt.Println("Exec 5XY0")
      X := getNibbleAt(instruction, 1)
      Y := getNibbleAt(instruction, 2)
      if V[X] == V[Y] {
        pc += 2
      }
    case 0x6:
      fmt.Println("Exec 6XNN")
      registerIndex := getNibbleAt(instruction, 1)
      NN := instruction & 0x00ff
      V[registerIndex] = uint8(NN)
    case 0x7:
      fmt.Println("Exec 7XNN")
      registerIndex := getNibbleAt(instruction, 1)
      NN := instruction & 0x00ff
      V[registerIndex] += uint8(NN)
    case 0x8:
      fmt.Println("decoded: 8")
      switch instruction & 0x000f {
      case 0x0:
        fmt.Println("Exec 8XY0")
        X := getNibbleAt(instruction, 1)
        Y := getNibbleAt(instruction, 2)
        V[X] = V[Y]
      case 0x1:
        fmt.Println("Exec  8XY1")
        X := getNibbleAt(instruction, 1)
        Y := getNibbleAt(instruction, 2)
        V[X] =  V[X] | V[Y]
      case 0x2:
        fmt.Println("Exec  8XY2")
        X := getNibbleAt(instruction, 1)
        Y := getNibbleAt(instruction, 2)
        V[X] =  V[X] & V[Y]
      case 0x3:
        fmt.Println("Exec  8XY3")
        X := getNibbleAt(instruction, 1)
        Y := getNibbleAt(instruction, 2)
        V[X] =  V[X] ^ V[Y]
      case 0x4:
        fmt.Println("Exec  8XY4")
        X := getNibbleAt(instruction, 1)
        Y := getNibbleAt(instruction, 2)
        if V[X]+V[Y] > 255 {
          V[15] = 1 
        } else {
          V[15] = 0
        }
        V[X] =  V[X] + V[Y]
      case 0x5:
        fmt.Println("Exec  8XY5")
        X := getNibbleAt(instruction, 1)
        Y := getNibbleAt(instruction, 2)
        if V[X] > V[Y] {
          V[15] = 1 
        } else {
          V[15] = 0
        }
        V[X] =  V[X] - V[Y]
      case 0x6:
        fmt.Println("Exec  8XY6")
        X := getNibbleAt(instruction, 1)
        Y := getNibbleAt(instruction, 2)
        // Execution for original COSMAC VIP programs
        if isLegacy {
          V[X] = V[Y]
        }
        shiftedBit := V[X] & 1
        V[X] >>= 1
        V[15] = shiftedBit
      case 0x7:
        fmt.Println("Exec  8XY7")
        X := getNibbleAt(instruction, 1)
        Y := getNibbleAt(instruction, 2)
        if V[Y] > V[X] {
          V[15] = 1 
        } else {
          V[15] = 0
        }
        V[X] =  V[Y] - V[X]
      case 0xe:
        fmt.Println("Exec  8XYE")
        X := getNibbleAt(instruction, 1)
        Y := getNibbleAt(instruction, 2)
        // Execution for original COSMAC VIP programs
        if isLegacy {
          V[X] = V[Y]
        }
        shiftedBit := V[X] & 128
        V[X] <<= 1
        V[15] = shiftedBit
      }
    case 0x9:
      fmt.Println("Exec 9XY0")
      X := getNibbleAt(instruction, 1)
      Y := getNibbleAt(instruction, 2)
      if V[X] != V[Y] {
        pc += 2
      }
    case 0xa:
      fmt.Println("Exec ANNN")
      NNN := instruction & 0x0fff
      l = uint8(NNN)
    case 0xb:
      fmt.Println("decoded: b")
    case 0xc:
      fmt.Println("decoded: c")
    case 0xd:
      fmt.Println("decoded: d")
    case 0xe:
      fmt.Println("decoded: e")
    case 0xf:
      fmt.Println("decoded: f")
    }
  }


  tick := time.Tick(2 * time.Millisecond)

  for { 
    instruction := fetchInstruction()
    decodeInstruction(instruction)
    <- tick
  }
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
  filename, isLegacy := parseArgs(os.Args)

  var memory [4096]byte
  var delayTimer *atomic.Int64 = timer()
  var soundTimer *atomic.Int64 = timer() 

  // Font runs from addr 050 to 09F
  registerFont(memory[0x50:])
  registerRom(filename, memory[0x200:])
  
  chip8(memory[:], delayTimer, soundTimer, isLegacy)
}

