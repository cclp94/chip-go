package display

import (
	"fmt"
	"time"

	keyboard "github.com/cclp94/chip-go/io/keyboard"
	"github.com/gopxl/pixel"
	"github.com/gopxl/pixel/pixelgl"
	"golang.org/x/image/colornames"
)

const (
	scale float64 = 10
)

type external struct {
  drawingChan *chan [][]byte
  kb keyboard.KeyboardInteface
}


func Init(
  drawingChan *chan [][]byte, 
  kb keyboard.KeyboardInteface,
) (func(), chan string) {
  _e := external{
    drawingChan: drawingChan,
    kb: kb,
  }
  selectionChan := make(chan string)

	display := func() {
		cfg := pixelgl.WindowConfig{
			Title:  "CHIP-8",
			Bounds: pixel.R(0, 0, 64*scale, 32*scale),
			VSync:  false,
		}
		win, err := pixelgl.NewWindow(cfg)
		if err != nil {
			panic(err)
		}
		win.Smooth()
		win.Clear(colornames.Skyblue)

		var (
			frames = 0
			fps    = time.NewTicker(16 * time.Millisecond).C
			second = time.NewTicker(time.Second).C
    )

    gs := createGameScene(win, &_e) 
    ss := createSelectScene(win)
    currentScene := ss
    
    gs.RegisterCallback(func(args ...string) {
      currentScene = ss
    })

    ss.RegisterCallback(func(args ...string) {
      currentScene = gs
      selection := args[0]
      selectionChan <- selection
    })

    for !win.Closed() {
      win.Clear(colornames.Black)
      currentScene.Draw()
      win.Update()

      // FPS control
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
  start := func () {
    pixelgl.Run(display)
  }

  return start, selectionChan
}

