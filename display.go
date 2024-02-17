package main

import (
	"fmt"
	"time"

	"github.com/gopxl/pixel"
	"github.com/gopxl/pixel/imdraw"
	"github.com/gopxl/pixel/pixelgl"
	"golang.org/x/image/colornames"
)

const (
	scale float64 = 20
)

func display(drawingChan *chan [][]byte, keyboardChan *chan byte) func() {
	return func() {
		cfg := pixelgl.WindowConfig{
			Title:  "Pixel Rocks!",
			Bounds: pixel.R(0, 0, 64*scale, 32*scale),
			VSync:  true,
		}
		win, err := pixelgl.NewWindow(cfg)
		if err != nil {
			panic(err)
		}
		win.Smooth()
		win.Clear(colornames.Skyblue)

		var (
			frames = 0
			second = time.Tick(time.Second)
		)

		imd := imdraw.New(nil)
		imd.Color = colornames.Darkgreen
		for !win.Closed() {
			if win.Pressed(pixelgl.Key1) {
				*keyboardChan <- 0x1
			} else if win.Repeated(pixelgl.Key2) {
				*keyboardChan <- 0x2
			} else if win.Pressed(pixelgl.Key3) {
				*keyboardChan <- 0x3
			} else if win.Pressed(pixelgl.Key4) {
				*keyboardChan <- 0xc
			} else if win.Pressed(pixelgl.KeyQ) {
				*keyboardChan <- 0x4
			} else if win.Pressed(pixelgl.KeyW) {
				*keyboardChan <- 0x5
			} else if win.Pressed(pixelgl.KeyE) {
				*keyboardChan <- 0x6
			} else if win.Pressed(pixelgl.KeyR) {
				*keyboardChan <- 0xd
			} else if win.Pressed(pixelgl.KeyA) {
				*keyboardChan <- 0x7
			} else if win.Pressed(pixelgl.KeyS) {
				*keyboardChan <- 0x8
			} else if win.Pressed(pixelgl.KeyD) {
				*keyboardChan <- 0x9
			} else if win.Pressed(pixelgl.KeyF) {
				*keyboardChan <- 0xe
			} else if win.Pressed(pixelgl.KeyZ) {
				*keyboardChan <- 0xa
			} else if win.Pressed(pixelgl.KeyX) {
				*keyboardChan <- 0x0
			} else if win.Pressed(pixelgl.KeyC) {
				*keyboardChan <- 0xb
			} else if win.Pressed(pixelgl.KeyV) {
				*keyboardChan <- 0xf
			}

			win.Clear(colornames.Black)
			select {
			case d := <-*drawingChan:
				imd.Clear()
				toImdraw(d, imd, win)
			default:
			}
			imd.Draw(win)
			win.Update()

			frames++
			select {
			case <-second:
				win.SetTitle(fmt.Sprintf("%s | FPS: %d", cfg.Title, frames))
				frames = 0
			default:
			}
		}
	}
}

func toImdraw(pixelDisplay [][]byte, imd *imdraw.IMDraw, win *pixelgl.Window) {
	for y := 31; y >= 0; y-- {
		for x := 63; x >= 0; x-- {
			if pixelDisplay[x][y] == 1 {
				nY := 31 - y
				imd.Push(pixel.V(float64(x)*scale, float64(nY)*scale), pixel.V(float64(x+1)*scale, float64(nY+1)*scale))
				imd.Rectangle(0)
				// time.Sleep(100 * time.Millisecond)
			}
		}
	}
}
