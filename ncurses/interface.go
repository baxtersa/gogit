package gogit

import (
	gh "github.com/baxtersa/gogit/github"
	common "github.com/baxtersa/gogit/internal"
	gc "github.com/rthornton128/goncurses"
)

var Quit = make(chan byte)
var windows = []View{}
var activeW View

func handleInput(c gc.Char, reqs *gh.ReqChannels) bool {
	switch c {
	case 'q':
		Quit <- byte(c)
		return false
	case 'r':
		reqs.Repo <- true
		return true
	case 'u':
		reqs.User <- true
		return true
	case 'i':
		reqs.Issue <- true
		return true
	default:
		return activeW.HandleInput(c, reqs)
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
	// Initialize ncurses terminal
	stdscr, err := gc.Init()
	common.Check(err)

	gc.Raw(true)
	gc.Echo(false)
	gc.Cursor(0)
	stdscr.Clear()
	stdscr.Keypad(true)

	// Create ncurses window
	rows, cols := stdscr.MaxYX()
	height, width := rows-5, cols-10
	y, x := (rows-height)/2, (cols-width)/2

	win, err := gc.NewWindow(height, width, y, x)
	common.Check(err)
	repos := Menu(win, []string{"foo", "bar"})
	activeW = &repos
	reposIdx := AddView(&repos)
	defer DeleteView(reposIdx)

	// Initialize stdin handling
	in := make(chan gc.Char)
	ready := make(chan bool)
	go readIn(stdscr, in, ready)

loop:
	for {
		stdscr.Move(0, 0)
		select {
		case s := <-resps.Repo:
			stdscr.Println("repos returned")
			repos.SetItems(s)
			UpdateViews()
			stdscr.Refresh()
		case s := <-resps.User:
			for _, str := range s {
				stdscr.Println(str)
			}
			UpdateViews()
			stdscr.Refresh()
		case s := <-resps.Issue:
			for _, str := range s {
				stdscr.Println(str)
			}
			UpdateViews()
			stdscr.Refresh()
		case c := <-in:
			if !handleInput(c, reqs) {
				break loop
			}
			stdscr.MovePrintf(rows-2, 0, "Character pressed: %c", rune(c))
			stdscr.ClearToEOL()
			UpdateViews()
			stdscr.Refresh()
		case ready <- true:
		}
	}
}

func End() {
	gc.End()
}
