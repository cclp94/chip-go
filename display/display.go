package display

import (
	"fmt"
	"log"
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


func Init(drawingChan *chan [][]byte, kb keyboard.KeyboardInteface) (func(), chan string) {
  _e := external{
    drawingChan: drawingChan,
    kb: kb,
  }
  selectionChan := make(chan string, 2)

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
			fps    = time.NewTicker(16 * time.Millisecond).C
			second = time.NewTicker(time.Second).C
    )

    gs := createGameScene(win, &_e) 
    ss := createSelectScene(win, selectionChan)
    currentScene := ss
    
    for !win.Closed() {
      win.Clear(colornames.Black)
      currentScene.Draw()
      win.Update()

      if currentScene == ss {
        select {
        case selection, done := <- selectionChan:
          log.Println(selection, done)
          // TODO code a way to make backspace go back to selection
          currentScene = gs
        default:
        }

      }

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

