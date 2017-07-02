package gogit

import (
	gh "github.com/baxtersa/gogit/github"
	gc "github.com/rthornton128/goncurses"
)

type menu struct {
	w *gc.Window

	items []*gc.MenuItem
	menu  *gc.Menu
}

func Menu(w *gc.Window, items []string) menu {
	menu_items := make([]*gc.MenuItem, len(items))
	for i, val := range items {
		menu_items[i], _ = gc.NewItem(val, "")
	}
	ncmenu, _ := gc.NewMenu(menu_items)
	ncmenu.SetWindow(w)

	rows, cols := w.MaxYX()
	height, width := rows-5, cols-10
	y, x := (rows-height)/2, (cols-width)/2
	dwin := w.Derived(height, width, y, x)
	ncmenu.SubWindow(dwin)
	ncmenu.Format(10, 1)

	w.Box(0, 0)

	ncmenu.Post()
	return menu{
		w:     w,
		items: menu_items,
		menu:  ncmenu,
	}
}

func (m *menu) Draw() {
}

func (m *menu) Free() {
	defer m.menu.UnPost()
	for i, _ := range m.items {
		defer m.items[i].Free()
	}
	defer m.menu.Free()
}

func (m *menu) HandleInput(c gc.Char, reqs *gh.ReqChannels) bool {
	PAGE_NUM := 5
	switch c {
	case 'j':
		fallthrough
	case gc.KEY_DOWN:
		m.menu.Driver(gc.REQ_DOWN)
	case 'k':
		fallthrough
	case gc.KEY_UP:
		m.menu.Driver(gc.REQ_UP)
	case '':
		for i := 0; i < PAGE_NUM; i++ {
			m.menu.Driver(gc.REQ_DOWN)
		}
	case '':
		for i := 0; i < PAGE_NUM; i++ {
			m.menu.Driver(gc.REQ_UP)
		}
	}
	return true
}

func (m *menu) Update() {
	m.w.Refresh()
}

func (m *menu) SetItems(is []string) {
	for i, _ := range m.items {
		defer m.items[i].Free()
	}
	m.items = make([]*gc.MenuItem, len(is))
	for i, val := range is {
		m.items[i], _ = gc.NewItem(val, "")
	}
	defer m.menu.Free()
	m.menu, _ = gc.NewMenu(m.items)
	m.menu.SetWindow(m.w)

	rows, cols := m.w.MaxYX()
	height, width := rows-5, cols-10
	y, x := (rows-height)/2, (cols-width)/2
	dwin := m.w.Derived(height, width, y, x)
	m.menu.SubWindow(dwin)
	m.menu.Format(height, 1)

	m.menu.Post()
}
