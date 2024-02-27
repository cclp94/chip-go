package chip8

import (
	"sync/atomic"
	"testing"

	"github.com/cclp94/chip-go/io/keyboard"
)

func testOp(opcode uint16, r1 uint8, r2 uint8) (chip8, []uint8) {
	var mockTimer atomic.Int64
	mockDisplayChan := make(chan [][]byte)
	mockKeyboard := keyboard.Create()

  var memory []byte = make([]byte, 4096)
  c := chip8{
    memory:   memory,
    pc:       0x200,
    registers: make([]uint8, 16),
    displayChan: mockDisplayChan,
    kb: mockKeyboard,
    delayTimer: &mockTimer,
    soundTimer: &mockTimer,
    isLegacy: true,
  }
	c.registers[0] = r1
	c.registers[1] = r2

	c.decodeOpcode(opcode)

	return c, c.registers
}

func Test_8XY1(t *testing.T) {
	_, V := testOp(0x8011, 0x2d, 0x4b)
	expected := uint8(0x6F)
	if V[0] != expected {
		t.Fatalf("Failed expected: %X, Got %X", expected, V[0])
	}
}

func Test_8XY2(t *testing.T) {
	_, V := testOp(0x8012, 0x2d, 0x4b)
	expected := uint8(0x09)
	if V[0] != expected {
		t.Fatalf("Failed expected: %X, Got %X", expected, V[0])
	}
}
func Test_8XY3(t *testing.T) {
	_, V := testOp(0x8013, 0x2d, 0x4b)
	expected := uint8(0x66)
	if V[0] != expected {
		t.Fatalf("Failed expected: %X, Got %X", expected, V[0])
	}
}
func Test_8XY4(t *testing.T) {
	_, V := testOp(0x8014, 0x2d, 0x4b)
	expected := uint8(0x78)
	if V[0] != expected {
		t.Fatalf("Failed expected: %X, Got %X", expected, V[0])
	}
	if V[15] != 0x0 {
		t.Fatalf("VF incorrectly set")
	}
}
func Test_8XY4_with_carry(t *testing.T) {
	_, V := testOp(0x8014, 0xed, 0x4b)
	expected := uint8(0x38)
	if V[0] != expected {
		t.Fatalf("Failed expected: %X, Got %X", expected, V[0])
	}
	if V[15] != 0x1 {
		t.Fatalf("VF incorrectly set")
	}
}
func Test_8XY5_no_borrow(t *testing.T) {
	_, V := testOp(0x8015, 0x4b, 0x2d)
	expected := uint8(0x1e)
	if V[0] != expected {
		t.Fatalf("Failed expected: %X, Got %X", expected, V[0])
	}
	if V[15] != 0x1 {
		t.Fatalf("VF incorrectly set")
	}
}
func Test_8XY5_with_borrow(t *testing.T) {
	_, V := testOp(0x8015, 0x2d, 0x4b)
	expected := uint8(0xe2)
	if V[0] != expected {
		t.Fatalf("Failed expected: %X, Got %X", expected, V[0])
	}
	if V[15] != 0x0 {
		t.Fatalf("VF incorrectly set")
	}
}
func Test_8XY6(t *testing.T) {
	_, V := testOp(0x8016, 0x00, 0x2c)
	expected := uint8(0x16)
	if V[0] != expected {
		t.Fatalf("Failed expected: %X, Got %X", expected, V[0])
	}
	if V[15] != 0x0 {
		t.Fatalf("VF incorrectly set")
	}
}

func Test_FXNN(t *testing.T) {
	c, _ := testOp(0xF033, 0x68, 0x0)
	if c.memory[c.i] != 1 {
		t.Fatalf("Failed first digit. Expected %d, got %d", 1, c.memory[c.i])
	}
	if c.memory[c.i+1] != 0 {
		t.Fatalf("Failed first digit. Expected %d, got %d", 0, c.memory[c.i+1])
	}
	if c.memory[c.i+2] != 4 {
		t.Fatalf("Failed first digit. Expected %d, got %d", 4, c.memory[c.i+2])
	}

}
