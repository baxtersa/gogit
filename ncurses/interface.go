package gogit

import (
	"fmt"
	"time"

	gh "github.com/baxtersa/gogit/github"
	common "github.com/baxtersa/gogit/internal"
	gc "github.com/rthornton128/goncurses"
)

var Quit = make(chan byte)

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

func handleInput(c gc.Char, reqs *gh.ReqChannels) bool {
	fmt.Printf(string(c))
	switch rune(c) {
	case 'q':
		Quit <- byte(c)
		return false
	case 'r':
		reqs.Repo <- true
	case 'u':
		reqs.User <- true
	case 'i':
		reqs.Issue <- true
	}
	return true
}

func readIn(w *gc.Window, ch chan<- gc.Char, ready <-chan bool) {
	for {
		// Block until all write operations are complete
		<-ready
		// Send typed character down the channel (which is blocking
		// in the main loop)
		ch <- gc.Char(w.GetChar())
	}
}

func Interface(reqs *gh.ReqChannels, resps *gh.RespChannels) {
	stdscr, err := gc.Init()
	common.Check(err)

	gc.Echo(false)
	gc.Cursor(0)
	stdscr.Clear()
	stdscr.Keypad(true)

	in := make(chan gc.Char)
	ready := make(chan bool)
	go readIn(stdscr, in, ready)

loop:
	for {
		select {
		case s := <-resps.Repo:
			for _, str := range s {
				fmt.Println(str)
			}
		case s := <-resps.User:
			for _, str := range s {
				fmt.Println(str)
			}
		case s := <-resps.Issue:
			for _, str := range s {
				fmt.Println(str)
			}
		case c := <-in:
			if !handleInput(c, reqs) {
				break loop
			}
		case ready <- true:
		}
	}
}

func End() {
	gc.End()
}
