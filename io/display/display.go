package display

import (
	"fmt"
	"time"

	"github.com/cclp94/chip-go/io/keyboard"
	"github.com/gopxl/pixel"
	"github.com/gopxl/pixel/imdraw"
	"github.com/gopxl/pixel/pixelgl"
	"golang.org/x/image/colornames"
)

const (
	scale float64 = 10
)

func Start(drawingChan *chan [][]byte, kb keyboard.KeyboardInteface) {
	display := func() {
		cfg := pixelgl.WindowConfig{
			Title:  "CHIP-8",
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
			fps    = time.NewTicker(1 * time.Millisecond).C
			second = time.NewTicker(time.Second).C
		)

		imd := imdraw.New(nil)
		imd.Color = colornames.Darkgreen

		for !win.Closed() {
			checkKeyPress(win, kb)
			checkKeyRelease(win, kb)

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
			<-fps
		}
	}
	pixelgl.Run(display)
}

func checkKeyRelease(win *pixelgl.Window, kb keyboard.KeyboardInteface) {
	if win.JustReleased(pixelgl.Key1) {
		kb.ReleaseKey(0x1)
	}
	if win.JustReleased(pixelgl.Key2) {
		kb.ReleaseKey(0x2)
	}
	if win.JustReleased(pixelgl.Key3) {
		kb.ReleaseKey(0x3)
	}
	if win.JustReleased(pixelgl.Key4) {
		kb.ReleaseKey(0xc)
	}
	if win.JustReleased(pixelgl.KeyQ) {
		kb.ReleaseKey(0x4)
	}
	if win.JustReleased(pixelgl.KeyW) {
		kb.ReleaseKey(0x5)
	}
	if win.JustReleased(pixelgl.KeyE) {
		kb.ReleaseKey(0x6)
	}
	if win.JustReleased(pixelgl.KeyR) {
		kb.ReleaseKey(0xd)
	}
	if win.JustReleased(pixelgl.KeyA) {
		kb.ReleaseKey(0x7)
	}
	if win.JustReleased(pixelgl.KeyS) {
		kb.ReleaseKey(0x8)
	}
	if win.JustReleased(pixelgl.KeyD) {
		kb.ReleaseKey(0x9)
	}
	if win.JustReleased(pixelgl.KeyF) {
		kb.ReleaseKey(0xe)
	}
	if win.JustReleased(pixelgl.KeyZ) {
		kb.ReleaseKey(0xa)
	}
	if win.JustReleased(pixelgl.KeyX) {
		kb.ReleaseKey(0x0)
	}
	if win.JustReleased(pixelgl.KeyC) {
		kb.ReleaseKey(0xb)
	}
	if win.JustReleased(pixelgl.KeyV) {
		kb.ReleaseKey(0xf)
	}
}

func checkKeyPress(win *pixelgl.Window, kb keyboard.KeyboardInteface) {
	if win.JustPressed(pixelgl.Key1) {
		kb.PressKey(0x1)
	}
	if win.JustPressed(pixelgl.Key2) {
		kb.PressKey(0x2)
	}
	if win.JustPressed(pixelgl.Key3) {
		kb.PressKey(0x3)
	}
	if win.JustPressed(pixelgl.Key4) {
		kb.PressKey(0xc)
	}
	if win.JustPressed(pixelgl.KeyQ) {
		kb.PressKey(0x4)
	}
	if win.JustPressed(pixelgl.KeyW) {
		kb.PressKey(0x5)
	}
	if win.JustPressed(pixelgl.KeyE) {
		kb.PressKey(0x6)
	}
	if win.JustPressed(pixelgl.KeyR) {
		kb.PressKey(0xd)
	}
	if win.JustPressed(pixelgl.KeyA) {
		kb.PressKey(0x7)
	}
	if win.JustPressed(pixelgl.KeyS) {
		kb.PressKey(0x8)
	}
	if win.JustPressed(pixelgl.KeyD) {
		kb.PressKey(0x9)
	}
	if win.JustPressed(pixelgl.KeyF) {
		kb.PressKey(0xe)
	}
	if win.JustPressed(pixelgl.KeyZ) {
		kb.PressKey(0xa)
	}
	if win.JustPressed(pixelgl.KeyX) {
		kb.PressKey(0x0)
	}
	if win.JustPressed(pixelgl.KeyC) {
		kb.PressKey(0xb)
	}
	if win.JustPressed(pixelgl.KeyV) {
		kb.PressKey(0xf)
	}
}

func toImdraw(pixelDisplay [][]byte, imd *imdraw.IMDraw, win *pixelgl.Window) {
	for y := 31; y >= 0; y-- {
		for x := 63; x >= 0; x-- {
			if pixelDisplay[x][y] == 1 {
				nY := 31 - y
				imd.Push(pixel.V(float64(x)*scale, float64(nY)*scale), pixel.V(float64(x+1)*scale, float64(nY+1)*scale))
				imd.Rectangle(0)
			}
		}
	}
}
