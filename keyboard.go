package main

import (
	"log"
	"slices"
	"sync"
)

type keyboard struct {
	pressedKeys []byte
	mu          sync.Mutex
}

type keyboardInteface interface {
	PressKey(k byte)
	ReleaseKey(k byte)
	IsKeyPressed(k byte) bool
	GetTopKeyPressed() (byte, bool)
}

func (kb *keyboard) PressKey(key byte) {
	kb.mu.Lock()
	if !slices.Contains(kb.pressedKeys, key) {
		log.Println("keyboard:", kb.pressedKeys)
		kb.pressedKeys = append(kb.pressedKeys, key)
	}
	kb.mu.Unlock()
}

func (kb *keyboard) ReleaseKey(key byte) {
	kb.mu.Lock()
	var newKeys []byte
	for _, k := range kb.pressedKeys {
		if k != key {
			newKeys = append(newKeys, k)
		}
	}
	kb.pressedKeys = newKeys
	log.Println("keyboard:", kb.pressedKeys)
	kb.mu.Unlock()
}

func (kb *keyboard) IsKeyPressed(key byte) bool {
	kb.mu.Lock()
	hasKey := slices.Contains(kb.pressedKeys, key)
	kb.mu.Unlock()
	return hasKey
}

func (kb *keyboard) hasKeyPressed() bool {
	return len(kb.pressedKeys) > 0
}

func (kb *keyboard) GetTopKeyPressed() (byte, bool) {
	if !kb.hasKeyPressed() {
		return 0xf, false
	}
	kb.mu.Lock()
	topKey := kb.pressedKeys[len(kb.pressedKeys)-1]
	kb.mu.Unlock()
	return topKey, true
}
