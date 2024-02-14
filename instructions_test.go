package main

import (
	"sync/atomic"
	"testing"
)

func testOp(opcode uint16, r1 uint8, r2 uint8) [16]uint8 {
  var mockTimer atomic.Int64
  mockDisplayChan := make(chan [][]byte) 

  c := Chip8{}
  c.V[0] = r1
  c.V[1] = r2

  c.decodeInstruction(opcode, &mockTimer, &mockTimer, &mockDisplayChan)

  return c.V
}

func Test_8XY1(t *testing.T) {
  V := testOp(0x8011, 0x2d, 0x4b)
  expected := uint8(0x6F)
  if  V[0] != expected {
    t.Fatalf("Failed expected: %X, Got %X", expected, V[0])
  }
}

func Test_8XY2(t *testing.T) {
  V := testOp(0x8012, 0x2d, 0x4b)
  expected := uint8(0x09)
  if  V[0] != expected {
    t.Fatalf("Failed expected: %X, Got %X", expected, V[0])
  }
}
func Test_8XY3(t *testing.T) {
  V := testOp(0x8013, 0x2d, 0x4b)
  expected := uint8(0x66)
  if  V[0] != expected {
    t.Fatalf("Failed expected: %X, Got %X", expected, V[0])
  }
}
func Test_8XY4(t *testing.T) {
  V := testOp(0x8014, 0x2d, 0x4b)
  expected := uint8(0x78)
  if  V[0] != expected {
    t.Fatalf("Failed expected: %X, Got %X", expected, V[0])
  }
  if V[15] != 0x0 {
    t.Fatalf("VF incorrectly set")
  }
}
func Test_8XY4_with_carry(t *testing.T) {
  V := testOp(0x8014, 0xed, 0x4b)
  expected := uint8(0x38)
  if  V[0] != expected {
    t.Fatalf("Failed expected: %X, Got %X", expected, V[0])
  }
  if V[15] != 0x1 {
    t.Fatalf("VF incorrectly set")
  }
}
