package gogit

import (
	gc "github.com/rthornton128/goncurses"
)

// Generic interface for various types of views to display user content
// and process user interaction
type View interface {
	Draw()
	Free()
	HandleInput(c gc.Char) bool
	Update()
}

// Global views
var views = make([]View, 0, 16)

// Iterate through and update all global views, publishing changes
// to the ui
func UpdateViews() {
	for _, vw := range views {
		vw.Update()
	}
	gc.Update()
}

// Iterate through and draw all global views
func DrawViews() {
	for _, vw := range views {
		vw.Draw()
	}
}

// Add a view `v` to the global views context
// params:
//   `v`: View to be added
//
// returns:
//   Index of the view `v` added into the global views context
func AddView(v View) int {
	views = append(views, v)
	return len(views) - 1
}

// Delete the view at index `n` of the global views context
// params:
//   `n`: Index of the view to be deleted
func DeleteView(n int) {
	views[n].Free()
	views = append(views[:n], views[n:]...)
}
