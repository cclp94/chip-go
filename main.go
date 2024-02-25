package main

import (
	"fmt"
	"os"
	"sync/atomic"

	"github.com/cclp94/chip-go/chip8"
	"github.com/cclp94/chip-go/display"
	"github.com/cclp94/chip-go/io/keyboard"
	"github.com/cclp94/chip-go/timer"
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
		0xF0, 0x80, 0xF0, 0x80, 0x80, // F
	}

	for i, f := range font {
		memory[i] = f
	}
}

func main() {
	// filename, _ := utils.ParseArgs(os.Args)

	var memory [4096]byte
	var delayTimer *atomic.Int64 = timer.Timer()
	var soundTimer *atomic.Int64 = timer.SoundTimer()
	var kb keyboard.KeyboardInteface = keyboard.Create()
	displayChan := make(chan [][]byte)

	// Font runs from addr 050 to 09F
	registerFont(memory[0x50:])

  start, onSelectRomChan := display.Init(&displayChan, kb)

  go func() {
    rom := <-onSelectRomChan
    registerRom(rom, memory[0x200:])
    chip8.Start(memory[:], delayTimer, soundTimer, &displayChan, kb, false)
  }()
  start()
}
