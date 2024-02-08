package main

import (
	"fmt"
	"math/rand"
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
      fmt.Println("Exec BNNN")
      NNN := instruction & 0x0fff
      pc = NNN + uint16(V[0])
    case 0xc:
      fmt.Println("Exec CXNN")
      X := getNibbleAt(instruction, 1)
      NN := instruction & 0x00ff
      random := uint16(rand.Intn(int(NN))) & NN
      V[X] = uint8(random)
    case 0xd:
      fmt.Printf("Exec %04X\n", instruction)
    case 0xe:
      switch instruction & 0x00ff {
      case 0x9E:
        fmt.Println("Exec EX9E")
        // TODO skip pc if key in VX is pressed
      case 0xA1:
        fmt.Println("Exec EXA1")
        //TODO skip pc if key in VX is not pressed
      }
    case 0xf:
      X := instruction & getNibbleAt(instruction, 1)
      switch instruction & 0x00ff {
      // Timers
      case 0x07:
        fmt.Println("Exec FX07")
        V[X] = uint8(delayTimer.Load())
      case 0x15:
        fmt.Println("Exec FX15")
        delayTimer.Add(int64(V[X]))
      case 0x18:
        fmt.Println("Exec FX18")
        soundTimer.Add(int64(V[X]))
      case 0x1E:
        fmt.Println("Exec FX1E")
        l = V[X]
      case 0x0A:
        fmt.Println("Exec FX0A")
        // TODO block until a key is pressed and then store key in VX
        // pc -= 2
      case 0x29:
        fmt.Println("Exec FX29")
        // set l to the last nibble of VX + the offset of the font stored in memory. l points to the address of the font character
        l = 0x50 + V[X] & 0x0f
      case 0x33:
        fmt.Println("Exec FX33")
        n := V[X]
        c3 := n % 10
        c2 := (c3 / 10) % 10
        c1 := (c2 / 10) % 10
        memory[l] = c1
        memory[l + 1] = c2
        memory[l + 2] = c3
      case 0x55:
        fmt.Println("Exec FX55")
        for i,v := range V[:X+1] {
          memory[uint16(l) + uint16(i)] = v
        }
      case 0x65:
        fmt.Println("Exec FX65")
        for i,v := range memory[l:uint16(l)+uint16(X)+1] {
          V[i] = v
        }
      }
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
func main2() {
  filename, isLegacy := parseArgs(os.Args)

  var memory [4096]byte
  var delayTimer *atomic.Int64 = timer()
  var soundTimer *atomic.Int64 = timer() 

  // Font runs from addr 050 to 09F
  registerFont(memory[0x50:])
  registerRom(filename, memory[0x200:])
  
  chip8(memory[:], delayTimer, soundTimer, isLegacy)
}

