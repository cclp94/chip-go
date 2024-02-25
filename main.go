package main

import (
	"log"
	"os"
	"sync/atomic"

	"github.com/cclp94/chip-go/chip8"
	"github.com/cclp94/chip-go/display"
	"github.com/cclp94/chip-go/io/keyboard"
	"github.com/cclp94/chip-go/timer"
)


func main() {
	var delayTimer *atomic.Int64 = timer.Timer()
	var soundTimer *atomic.Int64 = timer.SoundTimer()
	var kb keyboard.KeyboardInteface = keyboard.Create()
	displayChan := make(chan [][]byte)

  start, onSelectRomChan := display.Init(&displayChan, kb)

  go func() {
    var kill chan<- bool
    for {
      romPath := <-onSelectRomChan
      if kill != nil {
        kill <- true
      }
      rom := readFile(romPath)
      c8 := chip8.Init(delayTimer, soundTimer, displayChan, kb, false)
      kill = c8.Start(rom)
    }
  }()
  start()
}


func readFile(path string) []byte {
	file, err := os.ReadFile(path)
	if err != nil {
		log.Println("File:", path, "could not be read:", err.Error())
		panic(1)
	}

  return file
}
