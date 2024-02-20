package display

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/gopxl/pixel"
	"github.com/gopxl/pixel/pixelgl"
	"github.com/gopxl/pixel/text"
	"golang.org/x/image/font/basicfont"
)

type selectScene struct{
  availableSelection []string
  basicTxt *text.Text
}

type SelectScene interface {
  Draw(win *pixelgl.Window)
}

func createSelectScene() SelectScene {
  roms := getAllFiles("./roms")
  basicAtlas := text.NewAtlas(basicfont.Face7x13, text.ASCII)
  basicTxt := text.New(pixel.V(0, 0), basicAtlas)

  for i, rom := range roms {
    fmt.Fprintln(basicTxt, fmt.Sprintf("%d. %s", i, rom))
  }
  return &selectScene{
    availableSelection: roms,
    basicTxt: basicTxt,
  }
}

func (ss *selectScene) Draw(win *pixelgl.Window) {
  ss.basicTxt.Draw(win, pixel.IM)

}

func getAllFiles(path string) []string {
  var files []string
  roms, err := os.ReadDir(path)
  if err != nil {
    log.Println("Error reading directory", path)
    return files 
  }

  for _, rom := range roms {
    if rom.IsDir() {
      files = append(files, getAllFiles(path + "/" + rom.Name())...)
      continue
    }

    if !strings.HasSuffix(rom.Name(), ".ch8") {
    //  continue
     }

    files = append(files, path + "/" + rom.Name())
  }
  return files
}
