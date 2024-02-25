package chip8

import (
  "log"
  "math/rand"
  "sync/atomic"
  "time"

  "github.com/cclp94/chip-go/io/keyboard"
  "github.com/cclp94/chip-go/utils"
)


var font = []byte{
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
  0xF0, 0x80, 0xF0, 0x80, 0x80, // F
}

type chip8 struct {
  vDisplay [][]byte

  memory []byte

  pc        uint16
  i         uint16
  stack     []uint16
  registers []uint8


  displayChan chan [][]byte
  kb keyboard.KeyboardInteface
  delayTimer *atomic.Int64
  soundTimer *atomic.Int64

  isLegacy bool
  killSwitch chan bool
}

type Chip8 interface {
  Start(rom []byte) chan<- bool;
  Reset();
}


func Init(
  delayTimer *atomic.Int64,
  soundTimer *atomic.Int64,
  displayChan chan [][]byte,
  kb keyboard.KeyboardInteface,
  isLegacy bool,
) Chip8 {
  var memory []byte = make([]byte, 4096)
  c8 := chip8{
    memory:   memory,
    pc:       0x200,
    registers: make([]uint8, 16),
    displayChan: displayChan,
    kb: kb,
    delayTimer: delayTimer,
    soundTimer: soundTimer,
    isLegacy: isLegacy,
    killSwitch: make(chan bool, 1),
  }
  return &c8
}

func (c8 *chip8) Start(rom []byte) chan<- bool {
  c8.Reset()
  c8.clearDisplay()
  c8.registerMemory(rom, 0x200)

  go func() {
    tick := time.NewTicker(time.Second / 4000000).C
    for {
      opcode := c8.fecthOpcode()
      c8.decodeOpcode(opcode)
      <-tick
      select {
      case <- c8.killSwitch:
        return
      default:
      }
    } 
  }()

  return c8.killSwitch
}

func (c8 *chip8) Reset() {
  c8.pc = 0x200
  clear(c8.memory)
  clear(c8.stack)
  clear(c8.registers)
  c8.registerMemory(font, 0x50)
}

func (c8 *chip8) registerMemory(bytes []byte, memoryOffset int) {
  for i, b := range bytes {
    c8.memory[i+memoryOffset] = b
  }
}

func (c8 *chip8) clearDisplay() {
  c8.vDisplay = make([][]byte, 64)
  for r := range c8.vDisplay {
    c8.vDisplay[r] = make([]byte, 32)
  }
}

func (c8 *chip8) refreshDisplay(xCoord uint8, yCoord uint8, spriteCoord int, N int) {
  c8.registers[15] = 0
  for i := 0; i < int(N); i++ {
    sprite := c8.memory[spriteCoord+i]
    // for each bit in the sprite with 1 byte so start at 10000000 and shift right until 1
    for j, x := 0, xCoord; j < 8; j, x = j+1, x+1 {
      pixel := utils.GetBitAt(uint8(sprite), j)
      vPixel := c8.vDisplay[x][yCoord]
      c8.vDisplay[x][yCoord] = pixel ^ vPixel
      // log.Printf("p: %X, vp: %X, d: %X", pixel, vPixel, vDisplay[xCoord][yCoord])
      if vPixel == 1 {
        c8.registers[15] = 1
      }
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

func (c *chip8) fecthOpcode() uint16 {
  opcode := uint16(c.memory[c.pc])<<8 | uint16(c.memory[c.pc+1])
  c.pc += 2
  return opcode
}

func (c *chip8) decodeOpcode(
  opcode uint16,
) {
  // Binary mask first hex nibble
  log.Printf("Instruction: %X\t", opcode)
  switch utils.GetNibbleAt(opcode, 0) {
  case 0x0:
    switch opcode & 0x00ff {
    case 0xE0:
      c.clearDisplay()
    case 0xEE:
      log.Println("Exec 00EE")
      toPopAtIndex := len(c.stack) - 1
      c.pc = c.stack[toPopAtIndex]
      c.stack = c.stack[:toPopAtIndex]
    }
  case 0x1:
    log.Println("Exec 1NNN")
    addr := opcode & 0x0fff
    c.pc = uint16(addr)
  case 0x2:
    log.Println("Exec 2NNN")
    addr := opcode & 0x0fff
    c.stack = append(c.stack, uint16(c.pc))
    c.pc = uint16(addr)
  case 0x3:
    log.Println("Exec 3XNN")
    registerIndex := utils.GetNibbleAt(opcode, 1)
    NN := uint8(opcode & 0x00ff)
    if c.registers[registerIndex] == NN {
      c.pc += 2
    }
  case 0x4:
    log.Println("Exec 4XNN")
    registerIndex := utils.GetNibbleAt(opcode, 1)
    NN := uint8(opcode & 0x00ff)
    if c.registers[registerIndex] != NN {
      c.pc += 2
    }
  case 0x5:
    log.Println("Exec 5XY0")
    X := utils.GetNibbleAt(opcode, 1)
    Y := utils.GetNibbleAt(opcode, 2)
    if c.registers[X] == c.registers[Y] {
      c.pc += 2
    }
  case 0x6:
    log.Println("Exec 6XNN")
    registerIndex := utils.GetNibbleAt(opcode, 1)
    NN := opcode & 0x00ff
    c.registers[registerIndex] = uint8(NN)
  case 0x7:
    log.Println("Exec 7XNN")
    registerIndex := utils.GetNibbleAt(opcode, 1)
    NN := opcode & 0x00ff
    c.registers[registerIndex] += uint8(NN)
  case 0x8:
    switch opcode & 0x000f {
    case 0x0:
      log.Println("Exec 8XY0")
      X := utils.GetNibbleAt(opcode, 1)
      Y := utils.GetNibbleAt(opcode, 2)
      c.registers[X] = c.registers[Y]
    case 0x1:
      log.Println("Exec  8XY1")
      X := utils.GetNibbleAt(opcode, 1)
      Y := utils.GetNibbleAt(opcode, 2)
      c.registers[X] = c.registers[X] | c.registers[Y]
    case 0x2:
      log.Println("Exec  8XY2")
      X := utils.GetNibbleAt(opcode, 1)
      Y := utils.GetNibbleAt(opcode, 2)
      c.registers[X] = c.registers[X] & c.registers[Y]
    case 0x3:
      log.Println("Exec  8XY3")
      X := utils.GetNibbleAt(opcode, 1)
      Y := utils.GetNibbleAt(opcode, 2)
      c.registers[X] = c.registers[X] ^ c.registers[Y]
    case 0x4:
      log.Println("Exec  8XY4")
      X := utils.GetNibbleAt(opcode, 1)
      Y := utils.GetNibbleAt(opcode, 2)
      vx := c.registers[X]
      vy := c.registers[Y]
      c.registers[X] = vx + vy
      if int(vx)+int(vy) > 255 {
        c.registers[15] = 1
      } else {
        c.registers[15] = 0
      }
    case 0x5:
      log.Println("Exec  8XY5")
      X := utils.GetNibbleAt(opcode, 1)
      Y := utils.GetNibbleAt(opcode, 2)
      vx := c.registers[X]
      vy := c.registers[Y]
      c.registers[X] = vx - vy
      if vx >= vy {
        c.registers[15] = 1
      } else {
        c.registers[15] = 0
      }
    case 0x6:
      log.Println("Exec  8XY6")
      X := utils.GetNibbleAt(opcode, 1)
      Y := utils.GetNibbleAt(opcode, 2)
      // Execution for original COSMAC c.VIP programs
      if c.isLegacy {
        c.registers[X] = c.registers[Y]
      }
      shiftedBit := c.registers[X] & 1
      c.registers[X] >>= 1
      c.registers[15] = shiftedBit
    case 0x7:
      log.Println("Exec  8XY7")
      X := utils.GetNibbleAt(opcode, 1)
      Y := utils.GetNibbleAt(opcode, 2)
      vx := c.registers[X]
      vy := c.registers[Y]
      c.registers[X] = vy - vx
      if vy >= vx {
        c.registers[15] = 1
      } else {
        c.registers[15] = 0
      }
    case 0xe:
      log.Println("Exec  8XYE")
      X := utils.GetNibbleAt(opcode, 1)
      Y := utils.GetNibbleAt(opcode, 2)
      // Execution for original COSMAC c.VIP programs
      if c.isLegacy {
        c.registers[X] = c.registers[Y]
      }
      shiftedBit := utils.GetBitAt(c.registers[X], 0)
      c.registers[X] <<= 1
      c.registers[15] = shiftedBit
    }
  case 0x9:
    log.Println("Exec 9XY0")
    X := utils.GetNibbleAt(opcode, 1)
    Y := utils.GetNibbleAt(opcode, 2)
    if c.registers[X] != c.registers[Y] {
      c.pc += 2
    }
  case 0xa:
    NNN := opcode & 0x0fff
    log.Printf("Exec ANNN, point to %X\n", c.memory[NNN])
    c.i = uint16(NNN)
  case 0xb:
    log.Println("Exec BNNN")
    NNN := opcode & 0x0fff
    c.pc = NNN + uint16(c.registers[0])
  case 0xc:
    log.Println("Exec CXNN")
    X := utils.GetNibbleAt(opcode, 1)
    NN := opcode & 0x00ff
    random := uint16(rand.Intn(int(NN))) & NN
    c.registers[X] = uint8(random)
  case 0xd:
    log.Printf("Exec %04X\n", opcode)
    X := utils.GetNibbleAt(opcode, 1)
    Y := utils.GetNibbleAt(opcode, 2)
    N := utils.GetNibbleAt(opcode, 3)
    c.refreshDisplay(c.registers[X]%64, c.registers[Y]%32, int(c.i), int(N))
    c.displayChan <- c.vDisplay
  case 0xe:
    X := utils.GetNibbleAt(opcode, 1)
    switch opcode & 0x00ff {
    case 0x9E:
      log.Println("Exec EX9E")
      if c.kb.IsKeyPressed(c.registers[X]) {
        c.pc += 2
      }
    case 0xA1:
      log.Println("Exec EXA1")
      if !c.kb.IsKeyPressed(c.registers[X]) {
        c.pc += 2
      }
    }
  case 0xf:
    X := utils.GetNibbleAt(opcode, 1)
    switch opcode & 0x00ff {
      // Timers
    case 0x07:
      log.Println("Exec FX07")
      c.registers[X] = uint8(c.delayTimer.Load())
    case 0x15:
      log.Println("Exec FX15")
      c.delayTimer.Add(int64(c.registers[X]))
    case 0x18:
      log.Println("Exec FX18")
      c.soundTimer.Add(int64(c.registers[X]))
    case 0x1E:
      log.Println("Exec FX1E")
      c.i += uint16(c.registers[X])
    case 0x0A:
      log.Println("Exec FX0A")
      keyPressed, ok := c.kb.GetTopKeyPressed()
      if ok {
        c.registers[X] = keyPressed
      } else {
        c.pc -= 2
      }
    case 0x29:
      log.Println("Exec FX29")
      // set c.l to the c.last nibble of c.V[X] + the offset of the font stored in memory. c.l points to the address of the font character
      c.i = uint16(0x50 + (c.registers[X]&0xf)*5)
    case 0x33:
      log.Println("Exec FX33")
      n := c.registers[X]
      c3 := n % 10
      c2 := (n / 10) % 10
      c1 := (n / 100) % 10
      log.Println(n, c1, c2, c3)
      c.memory[c.i] = c1
      c.memory[c.i+1] = c2
      c.memory[c.i+2] = c3
    case 0x55:
      log.Println("Exec FX55")
      for i, v := range c.registers[:X+1] {
        c.memory[uint16(c.i)+uint16(i)] = v
      }
    case 0x65:
      log.Println("Exec FX65")
      for i, v := range c.memory[c.i : uint16(c.i)+uint16(X)+1] {
        c.registers[i] = v
      }
    }
  }
}
