package gogit

import (
	gh "github.com/baxtersa/gogit/github"
	common "github.com/baxtersa/gogit/internal"
	gc "github.com/rthornton128/goncurses"
)

// Channel to communicate with main cli function and terminate
// the program
var Quit = make(chan byte)

var activeW View

// How long to timeout waiting for navigation command input
var NAV_READ_TIMEOUT_MS int = 1000

// Navigate between different info views
// params:
//   `w`: Top-level window from which we read chars
//   `reqs`: Channels to make requests to GH client
//
// returns:
//   `true` on handling input
//   `false` on fatal error condition
func handleNavInput(w *gc.Window, reqs *gh.ReqChannels) bool {
	w.Timeout(1000)
	c := w.GetChar()
	w.Timeout(-1)
	switch c {
	case 0:
		// Nav command timed out
		return true
	case 'r':
		// Make GH request to get repos
		reqs.Repo <- true
		return true
	case 'u':
		// Make GH request to get authenticated user
		reqs.User <- true
		return true
	case 'i':
		// Make GH request to get issues
		reqs.Issue <- true
		return true
	default:
		return true
	}
}

// Handle reading a char from top-level stdscr
// params:
//   `w`: Top-level window from which we read chars
//   `c`: ncurses `Char` read from top-level window
//   `reqs`: Channels to make requests to GH client
//
// returns:
//   `true` on handling input
//   `false` on fatal error condition or input to quit
func handleInput(w *gc.Window, c gc.Char, reqs *gh.ReqChannels) bool {
	switch c {
	case 'q':
		// User prompted to quit the program
		Quit <- byte(c)
		return false
	case 'g':
		// User prompted to navigate between different views
		return handleNavInput(w, reqs)
	default:
		// Forward the command to the current view's handler
		return activeW.HandleInput(c)
	}
	return true
}

// Block and read a char from the top-level stdscr asynchronously
// params:
//   `w`: Top-level window from which we read chars
//   `ch`: Channel to write chars to for handling
//   `ready`: Channel to block until main routine is ready to continue reading
func readIn(w *gc.Window, ch chan<- gc.Char, ready <-chan bool) {
	for {
		// Block until all write operations are complete
		<-ready
		// Send typed character down the channel (which is blocking
		// in the main loop)
		ch <- gc.Char(w.GetChar())
	}
}

// Main TUI loop.
// params:
//   `reqs`: Channel to communicate with GH client to signal requests
//   `resps`: Channel to receive `[]string` responses from GH client
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

	// Create menu in `win`
	repos := Menu(win, []string{"foo", "bar"})
	// Add `repos` to global views
	reposIdx := AddView(&repos)
	defer DeleteView(reposIdx)

	// Initialize stdin handling
	in := make(chan gc.Char)
	ready := make(chan bool)
	go readIn(stdscr, in, ready)

	// Set `repos` as default active view
	activeW = &repos

loop:
	for {
		stdscr.Move(0, 0)
		select {
		case s := <-resps.Repo:
			// GH repositories request returned
			stdscr.Println("repos returned")
			repos.SetItems(s)
		case s := <-resps.User:
			// GH user request returned
			stdscr.Println(s[0])
		case s := <-resps.Issue:
			// GH issues request returned
			for _, str := range s {
				stdscr.Println(str)
			}
		case c := <-in:
			// Char read from stdscr
			if !handleInput(stdscr, c, reqs) {
				break loop
			}
			// Log character input
			stdscr.MovePrintf(rows-2, 0, "Character pressed: %c", rune(c))
			stdscr.ClearToEOL()
		case ready <- true:
			// Prime async stdscr reading to block
			continue
		}
		UpdateViews()
		stdscr.Refresh()
	}
}

// Terminate the ncurses environment
func End() {
	gc.End()
}
