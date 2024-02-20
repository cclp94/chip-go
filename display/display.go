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


func Start(drawingChan *chan [][]byte, kb keyboard.KeyboardInteface) {
  _e := external{
    drawingChan: drawingChan,
    kb: kb,
  }
  log.Println(_e)
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

    // gs := initGameScene() 
    ss := createSelectScene()
		for !win.Closed() {
			win.Clear(colornames.Black)
      // gs.Draw(win,&_e)
      ss.Draw(win)
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
	pixelgl.Run(display)
}

