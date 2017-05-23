package gogit

import (
	"time"

	common "github.com/baxtersa/gogit/internal"
	gc "github.com/rthornton128/goncurses"
)

const FPS_60HZ = time.Second / 60

type View interface {
	Draw(*gc.Window)
	Update()
}

var views = make([]View, 0, 16)

func updateViews(my int, mx int) {
}

func drawViews(s *gc.Window) {
	for _, vw := range views {
		vw.Draw(s)
	}
}

func handleInput(s *gc.Window) bool {
	k := s.GetChar()

	switch byte(k) {
	default:
		return false
	}
	return true
}

func Interface() {
	stdscr, err := gc.Init()
	common.Check(err)
	defer gc.End()

	gc.Echo(false)
	gc.Cursor(0)
	stdscr.Clear()
	stdscr.Keypad(true)

	frameRate := time.NewTicker(FPS_60HZ)

loop:
	for {
		select {
		case <-frameRate.C:
			updateViews(stdscr.MaxYX())
			drawViews(stdscr)
		default:
			if !handleInput(stdscr) {
				break loop
			}
		}
	}
}
