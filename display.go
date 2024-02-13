package main

import (

	"github.com/gopxl/pixel"
	"github.com/gopxl/pixel/imdraw"
	"github.com/gopxl/pixel/pixelgl"
	"golang.org/x/image/colornames"
)

const (
  scale float64 = 20
)

func display(drawingChan *chan [][]byte) func() {
  return func() {
    cfg := pixelgl.WindowConfig{
      Title:  "Pixel Rocks!",
      Bounds: pixel.R(0, 0, 64 * scale, 32 * scale),
      VSync:  true,
    }
    win, err := pixelgl.NewWindow(cfg)
    if err != nil {
      panic(err)
    }

    win.Clear(colornames.Skyblue)

    imd := imdraw.New(nil)
    imd.Color = colornames.Orange
    for !win.Closed() {
      win.Clear(colornames.Skyblue)
      select {
      case d:= <- *drawingChan:
        // imd.Reset()
        imd.Clear()
        toImdraw(d, imd, win) 
      default:
      }
      imd.Draw(win)
      win.Update()

    }
  }
}

func toImdraw(pixelDisplay [][]byte, imd *imdraw.IMDraw, win *pixelgl.Window) {
  for y:= 31; y >= 0; y-- {
    for x := 63; x >= 0; x-- {
      if pixelDisplay[x][y] == 1 {
        nY := 31 - y
        imd.Push(pixel.V(float64(x)*scale, float64(nY)*scale), pixel.V(float64(x+1) * scale, float64(nY+1) * scale))
        imd.Rectangle(0)
        // time.Sleep(100 * time.Millisecond)
      }
    }
  }
}
