package gogit

import (
	gc "github.com/rthornton128/goncurses"
)

type Menu struct {
	w *gc.Window
}

func (m *Menu) Draw() {
	m.w.Box(gc.ACS_VLINE, gc.ACS_HLINE)
}

func (m *Menu) Update() {
	gc.Update()
}

func (m *Menu) Free() {
	m.w.Delete()
}
