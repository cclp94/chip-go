package display

import (
	"github.com/gopxl/pixel"
	"github.com/gopxl/pixel/imdraw"
	"github.com/gopxl/pixel/pixelgl"
	"golang.org/x/image/colornames"
)

const (
  
)

type gameScene struct {
  win *pixelgl.Window

  e *external
  imd *imdraw.IMDraw 

  doneCallback func(args ...string)
}

type Scene interface {
  Draw()
  RegisterCallback(callback func(args ...string))
}

func createGameScene(win *pixelgl.Window, e *external) Scene {
  imd := imdraw.New(nil)
  imd.Color = colornames.Darkgreen
  gs := gameScene{win: win, imd: imd, e: e}
  return &gs
}

func (gs *gameScene) RegisterCallback(fn func(args ...string)) {
  gs.doneCallback = fn
}

func (gs *gameScene) Draw() {
  win, e := gs.win, gs.e
  checkKeyPress(win, e)
  checkKeyRelease(win, e)
  
  if win.JustPressed(pixelgl.KeyEscape) {
   gs.doneCallback() 
  }

  select {
  case d := <- *e.drawingChan:
    gs.imd.Clear()
    toImdraw(d,gs.imd, win)
  default:
  }
  gs.imd.Draw(win)
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

func checkKeyPress(win *pixelgl.Window, e *external) {
  if win.JustPressed(pixelgl.Key1) {
    e.kb.PressKey(0x1)
  }
  if win.JustPressed(pixelgl.Key2) {
    e.kb.PressKey(0x2)
  }
  if win.JustPressed(pixelgl.Key3) {
    e.kb.PressKey(0x3)
  }
  if win.JustPressed(pixelgl.Key4) {
    e.kb.PressKey(0xc)
  }
  if win.JustPressed(pixelgl.KeyQ) {
    e.kb.PressKey(0x4)
  }
  if win.JustPressed(pixelgl.KeyW) {
    e.kb.PressKey(0x5)
  }
  if win.JustPressed(pixelgl.KeyE) {
    e.kb.PressKey(0x6)
  }
  if win.JustPressed(pixelgl.KeyR) {
    e.kb.PressKey(0xd)
  }
  if win.JustPressed(pixelgl.KeyA) {
    e.kb.PressKey(0x7)
  }
  if win.JustPressed(pixelgl.KeyS) {
    e.kb.PressKey(0x8)
  }
  if win.JustPressed(pixelgl.KeyD) {
    e.kb.PressKey(0x9)
  }
  if win.JustPressed(pixelgl.KeyF) {
    e.kb.PressKey(0xe)
  }
  if win.JustPressed(pixelgl.KeyZ) {
    e.kb.PressKey(0xa)
  }
  if win.JustPressed(pixelgl.KeyX) {
    e.kb.PressKey(0x0)
  }
  if win.JustPressed(pixelgl.KeyC) {
    e.kb.PressKey(0xb)
  }
  if win.JustPressed(pixelgl.KeyV) {
    e.kb.PressKey(0xf)
  }
}
func checkKeyRelease(win *pixelgl.Window, e *external) {
  if win.JustReleased(pixelgl.Key1) {
    e.kb.ReleaseKey(0x1)
  }
  if win.JustReleased(pixelgl.Key2) {
    e.kb.ReleaseKey(0x2)
  }
  if win.JustReleased(pixelgl.Key3) {
    e.kb.ReleaseKey(0x3)
  }
  if win.JustReleased(pixelgl.Key4) {
    e.kb.ReleaseKey(0xc)
  }
  if win.JustReleased(pixelgl.KeyQ) {
    e.kb.ReleaseKey(0x4)
  }
  if win.JustReleased(pixelgl.KeyW) {
    e.kb.ReleaseKey(0x5)
  }
  if win.JustReleased(pixelgl.KeyE) {
    e.kb.ReleaseKey(0x6)
  }
  if win.JustReleased(pixelgl.KeyR) {
    e.kb.ReleaseKey(0xd)
  }
  if win.JustReleased(pixelgl.KeyA) {
    e.kb.ReleaseKey(0x7)
  }
  if win.JustReleased(pixelgl.KeyS) {
    e.kb.ReleaseKey(0x8)
  }
  if win.JustReleased(pixelgl.KeyD) {
    e.kb.ReleaseKey(0x9)
  }
  if win.JustReleased(pixelgl.KeyF) {
    e.kb.ReleaseKey(0xe)
  }
  if win.JustReleased(pixelgl.KeyZ) {
    e.kb.ReleaseKey(0xa)
  }
  if win.JustReleased(pixelgl.KeyX) {
    e.kb.ReleaseKey(0x0)
  }
  if win.JustReleased(pixelgl.KeyC) {
    e.kb.ReleaseKey(0xb)
  }
  if win.JustReleased(pixelgl.KeyV) {
    e.kb.ReleaseKey(0xf)
  }
}
