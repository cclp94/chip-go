package main

import (
	"fmt"
	"math/rand"
	"sync/atomic"
	"time"
)

type Chip8 struct {
  memory []byte
  pc uint16 
  l uint16
  stack []uint16
  V [16]uint8
  vDisplay [][]byte
  isLegacy bool
}

func chip8(
  memory []byte,
  delayTimer *atomic.Int64,
  soundTimer *atomic.Int64,
  displayChan *chan [][]byte,
  isLegacy bool,
) {
  c := Chip8 {
    memory: memory,
    pc: 0x200,
    isLegacy: isLegacy,
  }
  c.clearDisplay()

  tick := time.Tick(2 * time.Millisecond)
  for { 
    instruction := c.fetchInstruction()
    c.decodeInstruction(instruction, delayTimer, soundTimer, displayChan)
    <- tick
  }
}

func (c *Chip8) clearDisplay() {
  c.vDisplay = make([][]byte, 64)
  for r := range c.vDisplay {
    c.vDisplay[r] = make([]byte, 32)
  }
}

func (c *Chip8) refreshDisplay (xCoord uint8, yCoord uint8, spriteCoord int, N int) {
  c.V[15] = 0
  for i := 0; i < int(N); i++ {
    sprite := c.memory[spriteCoord + i]  
    // for each bit in the sprite with 1 byte so start at 10000000 and shift right until 1 
    for j, x := 0, xCoord; j < 8; j, x = j+1, x+1 {
      pixel := getBitAt(uint8(sprite), j)
      vPixel := c.vDisplay[x][yCoord]
      c.vDisplay[x][yCoord] = vPixel ^ pixel
      // fmt.Printf("p: %X, vp: %X, d: %X", pixel, vPixel, vDisplay[xCoord][yCoord])
      if vPixel == 1 {
        c.V[15] = 1
      } 
      // TODO stop if 63
      if x >= 63 {
        break
      }
    }
    yCoord += 1
    if yCoord >= 31 {
      return
    }
  }
}

func (c *Chip8) fetchInstruction  () uint16 {
  instruction := uint16(c.memory[c.pc])<<8 | uint16(c.memory[c.pc+1])
  c.pc += 2
  return instruction
}

func (c *Chip8) decodeInstruction (instruction uint16, delayTimer *atomic.Int64, soundTimer *atomic.Int64, displayChan *chan [][]byte) {
  // Binary mask first hex nibble
  fmt.Printf("Instruction: %X\t", instruction)
  switch getNibbleAt(instruction, 0) {
  case 0x0:
    fmt.Println("decoded: 0")
    switch instruction & 0x00ff {
    case 0xE0:
      c.clearDisplay()
    case 0xEE:
      fmt.Println("Exec 00EE")
      toPopAtIndex := len(c.stack) - 1
      c.pc = c.stack[toPopAtIndex]
      c.stack = c.stack[:toPopAtIndex]
    }
  case 0x1:
    fmt.Println("Exec 1NNN")
    addr := instruction & 0x0fff
    c.pc = uint16(addr)
  case 0x2:
    fmt.Println("Exec 2NNN")
    addr := instruction & 0x0fff
    c.stack = append(c.stack, uint16(c.pc))
    c.pc = uint16(addr)
  case 0x3:
    fmt.Println("Exec 3XNN")
    registerIndex := getNibbleAt(instruction, 1)
    NN := uint8(instruction & 0x00ff)
    if c.V[registerIndex] == NN {
      c.pc += 2
    }
  case 0x4:
    fmt.Println("Exec 4XNN")
    registerIndex := getNibbleAt(instruction, 1)
    NN := uint8(instruction & 0x00ff)
    if c.V[registerIndex] != NN {
      c.pc += 2
    }
  case 0x5:
    fmt.Println("Exec 5XY0")
    X := getNibbleAt(instruction, 1)
    Y := getNibbleAt(instruction, 2)
    if c.V[X] == c.V[Y] {
      c.pc += 2
    }
  case 0x6:
    fmt.Println("Exec 6XNN")
    registerIndex := getNibbleAt(instruction, 1)
    NN := instruction & 0x00ff
    c.V[registerIndex] = uint8(NN)
  case 0x7:
    fmt.Println("Exec 7XNN")
    registerIndex := getNibbleAt(instruction, 1)
    NN := instruction & 0x00ff
    c.V[registerIndex] += uint8(NN)
  case 0x8:
    fmt.Println("decoded: 8")
    switch instruction & 0x000f {
    case 0x0:
      fmt.Println("Exec 8XY0")
      X := getNibbleAt(instruction, 1)
      Y := getNibbleAt(instruction, 2)
      c.V[X] = c.V[Y]
    case 0x1:
      fmt.Println("Exec  8XY1")
      X := getNibbleAt(instruction, 1)
      Y := getNibbleAt(instruction, 2)
      c.V[X] =  c.V[X] | c.V[Y]
    case 0x2:
      fmt.Println("Exec  8XY2")
      X := getNibbleAt(instruction, 1)
      Y := getNibbleAt(instruction, 2)
      c.V[X] =  c.V[X] & c.V[Y]
    case 0x3:
      fmt.Println("Exec  8XY3")
      X := getNibbleAt(instruction, 1)
      Y := getNibbleAt(instruction, 2)
      c.V[X] =  c.V[X] ^ c.V[Y]
    case 0x4:
      fmt.Println("Exec  8XY4")
      X := getNibbleAt(instruction, 1)
      Y := getNibbleAt(instruction, 2)
      vx := c.V[X]
      vy := c.V[Y]

      c.V[X] = vx + vy
      if int(vx)+int(vy) > 255 {
        c.V[15] = 1 
      } else {
        c.V[15] = 0
      }
    case 0x5:
      fmt.Println("Exec  8XY5")
      X := getNibbleAt(instruction, 1)
      Y := getNibbleAt(instruction, 2)
      if c.V[X] > c.V[Y] {
        c.V[15] = 1 
      } else {
        c.V[15] = 0
      }
      c.V[X] =  c.V[X] - c.V[Y]
    case 0x6:
      fmt.Println("Exec  8XY6")
      X := getNibbleAt(instruction, 1)
      Y := getNibbleAt(instruction, 2)
      // Execution for original COSMAC c.VIP programs
      if c.isLegacy {
        c.V[X] = c.V[Y]
      }
      shiftedBit := c.V[X] & 1
      c.V[X] >>= 1
      c.V[15] = shiftedBit
    case 0x7:
      fmt.Println("Exec  8XY7")
      X := getNibbleAt(instruction, 1)
      Y := getNibbleAt(instruction, 2)
      if c.V[Y] > c.V[X] {
        c.V[15] = 1 
      } else {
        c.V[15] = 0
      }
      c.V[X] =  c.V[Y] - c.V[X]
    case 0xe:
      fmt.Println("Exec  8XYE")
      X := getNibbleAt(instruction, 1)
      Y := getNibbleAt(instruction, 2)
      // Execution for original COSMAC c.VIP programs
      if c.isLegacy {
        c.V[X] = c.V[Y]
      }
      shiftedBit := getBitAt(c.V[X], 0)
      c.V[X] <<= 1
      c.V[15] = shiftedBit
    }
  case 0x9:
    fmt.Println("Exec 9XY0")
    X := getNibbleAt(instruction, 1)
    Y := getNibbleAt(instruction, 2)
    if c.V[X] != c.V[Y] {
      c.pc += 2
    }
  case 0xa:
    NNN := instruction & 0x0fff
    fmt.Printf("Exec ANNN, point to %X", c.memory[NNN])
    c.l = uint16(NNN)
  case 0xb:
    fmt.Println("Exec BNNN")
    NNN := instruction & 0x0fff
    c.pc = NNN + uint16(c.V[0])
  case 0xc:
    fmt.Println("Exec CXNN")
    X := getNibbleAt(instruction, 1)
    NN := instruction & 0x00ff
    random := uint16(rand.Intn(int(NN))) & NN
    c.V[X] = uint8(random)
  case 0xd:
    fmt.Printf("Exec %04X\n", instruction)
    X := getNibbleAt(instruction, 1)
    Y := getNibbleAt(instruction, 2)
    N := getNibbleAt(instruction, 3)
    fmt.Println(X, Y, N)
    c.refreshDisplay(c.V[X] % 64, c.V[Y] % 32, int(c.l), int(N))
    *displayChan <- c.vDisplay
  case 0xe:
    switch instruction & 0x00ff {
    case 0x9E:
      fmt.Println("Exec EX9E")
      // TODO skip pc if key in c.V[X] is pressed
    case 0xA1:
      fmt.Println("Exec EXA1")
      //TODO skip pc if key in c.V[X] is not pressed
    }
  case 0xf:
    X := instruction & getNibbleAt(instruction, 1)
    switch instruction & 0x00ff {
      // Timers
    case 0x07:
      fmt.Println("Exec FX07")
      c.V[X] = uint8(delayTimer.Load())
    case 0x15:
      fmt.Println("Exec FX15")
      delayTimer.Add(int64(c.V[X]))
    case 0x18:
      fmt.Println("Exec FX18")
      soundTimer.Add(int64(c.V[X]))
    case 0x1E:
      fmt.Println("Exec FX1E")
      c.l += uint16(c.V[X])
    case 0x0A:
      fmt.Println("Exec FX0A")
      // TODO block until a key is pressed and then store key in c.V[X]
      // pc -= 2
    case 0x29:
      fmt.Println("Exec FX29")
      // set c.l to the c.last nibble of c.V[X] + the offset of the font stored in memory. c.l points to the address of the font character
      c.l = uint16(0x50 + c.V[X] & 0x0f)
    case 0x33:
      fmt.Println("Exec FX33")
      n := c.V[X]
      c3 := n % 10
      c2 := (c3 / 10) % 10
      c1 := (c2 / 10) % 10
      c.memory[c.l] = c1
      c.memory[c.l + 1] = c2
      c.memory[c.l + 2] = c3
    case 0x55:
      fmt.Println("Exec FX55")
      for i,v := range c.V[:X+1] {
        c.memory[uint16(c.l) + uint16(i)] = v
      }
    case 0x65:
      fmt.Println("Exec FX65")
      for i,v := range c.memory[c.l:uint16(c.l)+uint16(X)+1] {
        c.V[i] = v
      }
    }
  }
}


