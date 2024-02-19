package utils

import "fmt"

func ParseArgs(args []string) (string, bool) {
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
