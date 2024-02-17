package main

import "fmt"

func getNibbleAt(value uint16, at int) uint16 {
	return uint16((value & (0xf000 >> (at * 4))) >> ((3 - at) * 4))
}

func getBitAt(value uint8, at int) uint8 {
	return (value & (0x80 >> at)) >> (7 - at)
}

func parseArgs(args []string) (string, bool) {
	var filename string
	var isLegacy bool

	if len(args) < 2 {
		fmt.Println("USAGE: <rom-file> [--legacy]")
		panic(1)
	}
	filename = args[1]

	if len(args) == 3 && args[2] == "--legacy" {
		isLegacy = true
	}
	return filename, isLegacy
}
