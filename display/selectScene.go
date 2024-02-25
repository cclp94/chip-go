package display

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/gopxl/pixel"
	"github.com/gopxl/pixel/pixelgl"
	"github.com/gopxl/pixel/text"
	"github.com/pkg/errors"
	"golang.org/x/image/font/basicfont"
)

const showRange int = 9

type selectScene struct{
  win *pixelgl.Window
  
  availableSelection []string

  basicTxt *text.Text
  
  _selectionOption string
  _shownPage int


  callbackFn func(args ...string)
}

func createSelectScene(win *pixelgl.Window) Scene {
  roms := getAllFiles("./roms")
  basicAtlas := text.NewAtlas(basicfont.Face7x13, text.ASCII)

  basicTxt := text.New(pixel.V(0, win.Bounds().Max.Y), basicAtlas)

  ss := &selectScene{
    win: win,
    availableSelection: roms,
    basicTxt: basicTxt,
    _shownPage: -1,
  }
  ss.nextPage()

  return ss 
}

func (ss * selectScene) RegisterCallback(fn func(args ...string)) {
  ss.callbackFn = fn
}

func (ss *selectScene) Draw() {
  ss.basicTxt.Draw(ss.win, pixel.IM.Moved(pixel.V(0, -ss.basicTxt.Bounds().H())))
  ss.checkKeyboardInput();
}

func (ss *selectScene) Done(selection int) {
  log.Printf("Selected %d: %s", selection, ss.availableSelection[selection])
  ss.callbackFn(ss.availableSelection[selection])
}

func (ss *selectScene) checkKeyboardInput() {
  typed := ss.win.Typed()
  if len(typed) > 0 {
    selection, err := ss.toSelection(typed)
    if err != nil {
      log.Println(err.Error())
      return
    }

    ss.Done(selection-1)
  }

  if ss.win.JustPressed(pixelgl.KeyEnter) {
    ss.nextPage()
  }

  if ss.win.JustPressed(pixelgl.KeyBackspace) {
    ss.prevPage()
  }
}

func (ss *selectScene) refreshPage() {
  ss.basicTxt.Clear()
  low := showRange*ss._shownPage
  high := showRange + showRange*ss._shownPage
  for i, rom := range ss.availableSelection[low:high] {
    fmt.Fprintln(ss.basicTxt, fmt.Sprintf("%d. %s", i+1, rom))
  }
  fmt.Fprintln(ss.basicTxt, fmt.Sprint("\n\nBACKSPACE ...Back"))
  fmt.Fprintln(ss.basicTxt, fmt.Sprint("ENTER Next..."))
}

func (ss *selectScene) prevPage() {
  if ss._shownPage > 0 {
    ss._shownPage--
    ss.refreshPage()
  }
}

func (ss *selectScene) nextPage() {
  low := showRange*ss._shownPage
  if low < len(ss.availableSelection) {
    ss._shownPage++
    ss.refreshPage()
  }
}

func (ss *selectScene) toSelection(typed string) (int, error) {
  option, err := strconv.Atoi(typed)
  if err != nil {
    log.Println("Invalid option. Try again")
    return 0, errors.New("Invalid Option")
  }

  selection := option + ss._shownPage*10 - (1* ss._shownPage)
  isValid := selection >= 0 && selection < len(ss.availableSelection) 
  if isValid {
    return selection, nil
  }

  return 0, errors.New("Invalid Option")
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
      continue
    }

    files = append(files, path + "/" + rom.Name())
  }
  return files
}
