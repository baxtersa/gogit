package gogit

import (
	gh "github.com/baxtersa/gogit/github"
	common "github.com/baxtersa/gogit/internal"
	gc "github.com/rthornton128/goncurses"
)

var Quit = make(chan byte)
var windows = []View{}

func handleInput(c gc.Char, reqs *gh.ReqChannels) bool {
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
	// Initialize ncurses terminal
	stdscr, err := gc.Init()
	common.Check(err)

	gc.Echo(false)
	gc.Cursor(0)
	stdscr.Clear()
	stdscr.Keypad(true)

	// Create ncurses window
	rows, cols := stdscr.MaxYX()
	height, width := 5, 10
	y, x := (rows-height)/2, (cols-width)/2

	win, err := gc.NewWindow(height, width, y, x)
	common.Check(err)
	menu := Menu{w: win}
	menuIdx := AddView(&menu)
	menu.Draw()
	menu.Update()
	defer DeleteView(menuIdx)

	// Initialize stdin handling
	in := make(chan gc.Char)
	ready := make(chan bool)
	go readIn(win, in, ready)

loop:
	for {
		select {
		case s := <-resps.Repo:
			stdscr.Move(0, 0)
			for _, str := range s {
				stdscr.Println(str)
			}
			UpdateViews()
			stdscr.Refresh()
		case s := <-resps.User:
			stdscr.Move(0, 0)
			for _, str := range s {
				stdscr.Println(str)
			}
			UpdateViews()
			stdscr.Refresh()
		case s := <-resps.Issue:
			stdscr.Move(0, 0)
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
			DrawViews()
			UpdateViews()
			stdscr.Refresh()
		case ready <- true:
		}
	}
}

func End() {
	gc.End()
}
