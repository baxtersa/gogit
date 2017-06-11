package gogit

import (
	gc "github.com/rthornton128/goncurses"
)

type View interface {
	Draw()
	Update()
	Free()
}

var views = make([]View, 0, 16)

func UpdateViews() {
	for _, vw := range views {
		vw.Update()
	}
	gc.Update()
}

func DrawViews() {
	for _, vw := range views {
		vw.Draw()
	}
}

func AddView(v View) int {
	views = append(views, v)
	return len(views) - 1
}

func DeleteView(n int) {
	views[n].Free()
	views = append(views[:n], views[n:]...)
}
