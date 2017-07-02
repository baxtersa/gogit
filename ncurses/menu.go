package gogit

import (
	gc "github.com/rthornton128/goncurses"
)

// Wrapper type around ncurses menu
type menu struct {
	w *gc.Window

	items []*gc.MenuItem
	menu  *gc.Menu
}

// Construct a new menu
// params:
//   `w`: Window within which to create a menu
//   `items`: Initial contents of the new menu
//
// returns:
//   A menu `View`
func Menu(w *gc.Window, items []string) menu {
	// Initialize ncurses menu items out of string contents
	menu_items := make([]*gc.MenuItem, len(items))
	for i, val := range items {
		menu_items[i], _ = gc.NewItem(val, "")
	}

	// Create the ncurses menu
	ncmenu, _ := gc.NewMenu(menu_items)
	ncmenu.SetWindow(w)

	// Size and create subwindow to display contents
	rows, cols := w.MaxYX()
	height, width := rows-5, cols-10
	y, x := (rows-height)/2, (cols-width)/2
	dwin := w.Derived(height, width, y, x)
	ncmenu.SubWindow(dwin)
	ncmenu.Format(10, 1)

	// Box the menu
	w.Box(0, 0)

	// Post the menu, making it visible
	ncmenu.Post()
	return menu{
		w:     w,
		items: menu_items,
		menu:  ncmenu,
	}
}

// Drawing is managed by ncurses menu driver
func (m *menu) Draw() {
}

// Free allocated resources of ncurses menu and items
func (m *menu) Free() {
	defer m.menu.UnPost()
	for i, _ := range m.items {
		defer m.items[i].Free()
	}
	defer m.menu.Free()
}

// Handle character input forwarded to this menu
// params:
//   `c`: Character read from input
//
// returns:
//   `true` on successful handling, `false` on fatal error
func (m *menu) HandleInput(c gc.Char) bool {
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

// Update the menu's window
func (m *menu) Update() {
	m.w.Refresh()
}

// Update the contents of the menu
// params:
//   `is`: New contents to be displayed
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
