package main

import (
	"strings"

	"github.com/nsf/termbox-go"
)

type UI struct {
	dialog *Dialog
}

type Size struct {
	Width, Height int
}

func MakeSize(width, height int) Size {
	if width < 0 {
		width = -width
	}
	if height < 0 {
		height = -height
	}
	return Size{width, height}
}

func (s Size) Union(o Size) Size {
	u := s
	if u.Width < o.Width {
		u.Width = o.Width
	}
	if u.Height < o.Height {
		u.Height = o.Height
	}
	return u
}

func NewUI() *UI {
	return new(UI)
}

func (ui *UI) Show(dialog *Dialog) {
	ui.dialog = dialog
}

func (ui *UI) Dismiss() {
	ui.dialog = nil
}

func (ui *UI) Handle(e termbox.Event) bool {
	if ui.dialog == nil {
		return false
	}
	if IsKey(e, termbox.KeyEsc) {
		ui.Dismiss()
	} else if ui.dialog != nil {
		ui.dialog.Handle(ui, e)
	}
	return true
}

func (ui *UI) Render() {
	if ui.dialog != nil {
		ui.dialog.Render()
	}
}

type Dialog struct {
	x, y         int
	stack        []*Settings
	selectedItem int
}

func NewDialog(x, y int, model *Settings) *Dialog {
	return &Dialog{x, y, []*Settings{model}, 0}
}

func (d *Dialog) Page() *Settings {
	return d.stack[len(d.stack)-1]
}

func (d *Dialog) Breadcrumb() string {
	if len(d.stack) < 2 {
		return d.Page().Title
	}
	var titles []string
	for _, page := range d.stack[len(d.stack)-2:] {
		titles = append(titles, page.Title)
	}
	return strings.Join(titles, " > ")
}

func (d *Dialog) Enter(page *Settings) {
	d.stack = append(d.stack, page)
	d.selectedItem = 0 //TODO stack to restore in Back()
}

func (d *Dialog) Back() {
	if len(d.stack) > 1 {
		d.stack = d.stack[:len(d.stack)-1]
		d.selectedItem = 0
	}
}

func (d *Dialog) Handle(ui *UI, e termbox.Event) bool {
	if e.Type != termbox.EventKey {
		return false
	}
	switch e.Key {
	case termbox.KeyEsc:
		return true

	case termbox.KeyArrowDown:
		d.selectedItem++
	case termbox.KeyArrowUp:
		d.selectedItem--
	case termbox.KeyArrowLeft, termbox.KeyBackspace:
		if len(d.stack) > 1 {
			d.Back()
		} else {
			ui.Dismiss()
		}

	default:
		d.Page().Items[d.selectedItem].Handle(d, e)
	}
	d.selectedItem = clamp(d.selectedItem, 0, len(d.Page().Items)-1)
	return false
}

func (d *Dialog) Render() {
	page := d.Page()
	title := d.Breadcrumb()
	sz := page.PreferredSize().Union(MakeSize(len(title)+4, 0))
	DrawBox(d.x, d.y, d.x+sz.Width, d.y+sz.Height, termbox.ColorDefault, termbox.ColorDefault)
	DrawString(d.x+1, d.y, " "+title+" ", termbox.ColorDefault, termbox.ColorDefault)
	for n, item := range page.Items {
		fg, bg := termbox.ColorDefault, termbox.ColorDefault
		if n == d.selectedItem {
			fg = fg | termbox.AttrReverse
		}
		DrawString(d.x+1, d.y+1+n, item.String(sz.Width), fg, bg)
	}
}

type Settings struct {
	Title string
	Items []settingsItem
}

func NewSettings(title string) *Settings {
	return &Settings{title, nil}
}

func (s *Settings) AddMenu(title string, build func(page *Settings)) {
	page := &Settings{title, nil}
	s.Items = append(s.Items, &menuItem{title, page})
	build(page)
}

func (s *Settings) AddCheckbox(title string, on bool) {
	s.Items = append(s.Items, &checkboxItem{title, on})
}

func (s *Settings) PreferredSize() Size {
	sz := MakeSize(len(s.Title)+4, len(s.Items)+1)
	for _, item := range s.Items {
		minWidth := item.MinWidth() + 2
		if sz.Width < minWidth {
			sz.Width = minWidth
		}
	}
	return sz
}

type settingsItem interface {
	Handle(*Dialog, termbox.Event)
	String(int) string
	MinWidth() int
}

type menuItem struct {
	label string
	page  *Settings
}

func (m *menuItem) Handle(dialog *Dialog, e termbox.Event) {
	if IsKey(e, termbox.KeyArrowRight, termbox.KeyEnter) {
		dialog.Enter(m.page)
	}
}

func (m *menuItem) String(width int) string {
	return LayoutString(m.label, "", width)
}

func (m *menuItem) MinWidth() int {
	return len(m.label)
}

type checkboxItem struct {
	label string
	on    bool
}

func (c *checkboxItem) Handle(dialog *Dialog, e termbox.Event) {
	if IsKey(e, termbox.KeyArrowRight, termbox.KeyEnter) {
		c.on = !c.on
	}
}

func (c *checkboxItem) String(width int) string {
	box := "[ ]"
	if c.on {
		box = "[x]"
	}
	return LayoutString(c.label, box, width)
}

func (c *checkboxItem) MinWidth() int {
	return len(c.label) + 4
}

func LayoutString(lhs, rhs string, width int) string {
	if len(lhs)+1+len(rhs) > width {
		lhs = strings.TrimSpace(lhs[:width-len(rhs)-5]) + "..."
	}
	return lhs + strings.Repeat(" ", width-len(rhs)-1-len(lhs)) + rhs
}
