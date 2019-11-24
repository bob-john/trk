package main

import (
	"strings"

	"github.com/nsf/termbox-go"
)

type UI struct {
	dialog *Dialog
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
	x, y  int
	stack []*OptionPage
}

func NewDialog(x, y int, model *OptionPage) *Dialog {
	return &Dialog{x, y, []*OptionPage{model}}
}

func (d *Dialog) Page() *OptionPage {
	return d.stack[len(d.stack)-1]
}

func (d *Dialog) Breadcrumb() string {
	if len(d.stack) < 2 {
		return d.Page().title
	}
	var titles []string
	for _, page := range d.stack[len(d.stack)-2:] {
		titles = append(titles, page.title)
	}
	return strings.Join(titles, " > ")
}

func (d *Dialog) Enter(page *OptionPage) {
	d.stack = append(d.stack, page)
	page.selectedItem = 0 //HACK
}

func (d *Dialog) Back() {
	if len(d.stack) > 1 {
		d.stack = d.stack[:len(d.stack)-1]
	}
}

func (d *Dialog) Handle(ui *UI, e termbox.Event) {
	if e.Type != termbox.EventKey {
		return
	}
	switch e.Key {
	case termbox.KeyEsc:
		ui.Dismiss()
	case termbox.KeyArrowLeft, termbox.KeyBackspace:
		if len(d.stack) > 1 {
			d.Back()
		} else {
			ui.Dismiss()
		}
	default:
		d.Page().Handle(d, e)
	}
	return
}

func (d *Dialog) Render() {
	page := d.Page()
	title := d.Breadcrumb()
	sz := page.PreferredSize().Union(MakeSize(len(title)+4, 0))
	DrawBox(d.x, d.y, d.x+sz.Width, d.y+sz.Height, termbox.ColorDefault, termbox.ColorDefault)
	DrawString(d.x+1, d.y, " "+title+" ", termbox.ColorDefault, termbox.ColorDefault)
	page.Render(d.x+1, d.y+1, sz.Width)
}

type OptionPage struct {
	title        string
	items        []OptionItem
	selectedItem int
}

func NewOptionPage(title string) *OptionPage {
	return &OptionPage{title, nil, 0}
}

func (p *OptionPage) AddMenu(title string, build func(page *OptionPage)) {
	page := &OptionPage{title, nil, 0}
	p.items = append(p.items, &Menu{title, page})
	build(page)
}

func (p *OptionPage) AddCheckbox(title string, on bool) {
	p.items = append(p.items, &Checkbox{title, on})
}

func (p *OptionPage) PreferredSize() Size {
	sz := MakeSize(len(p.title)+4, len(p.items)+1)
	for _, item := range p.items {
		minWidth := item.MinWidth() + 2
		if sz.Width < minWidth {
			sz.Width = minWidth
		}
	}
	return sz
}

func (p *OptionPage) Handle(d *Dialog, e termbox.Event) {
	if e.Type != termbox.EventKey {
		return
	}
	switch e.Key {
	case termbox.KeyArrowDown:
		p.selectedItem++
	case termbox.KeyArrowUp:
		p.selectedItem--
	default:
		p.items[p.selectedItem].Handle(d, e)
	}
	p.selectedItem = clamp(p.selectedItem, 0, len(p.items)-1)
}

func (p *OptionPage) Render(x, y, width int) {
	for n, item := range p.items {
		fg, bg := termbox.ColorDefault, termbox.ColorDefault
		if n == p.selectedItem {
			fg = fg | termbox.AttrReverse
		}
		DrawString(x, y+n, item.String(width), fg, bg)
	}
}

type OptionItem interface {
	Handle(*Dialog, termbox.Event)
	String(int) string
	MinWidth() int
}

type Menu struct {
	label string
	page  *OptionPage
}

func (m *Menu) Handle(dialog *Dialog, e termbox.Event) {
	if IsKey(e, termbox.KeyArrowRight, termbox.KeyEnter) {
		dialog.Enter(m.page)
	}
}

func (m *Menu) String(width int) string {
	return LayoutString(m.label, "", width)
}

func (m *Menu) MinWidth() int {
	return len(m.label)
}

type Checkbox struct {
	label string
	on    bool
}

func (c *Checkbox) Handle(dialog *Dialog, e termbox.Event) {
	if IsKey(e, termbox.KeyArrowRight, termbox.KeyEnter) {
		c.on = !c.on
	}
}

func (c *Checkbox) String(width int) string {
	box := "[ ]"
	if c.on {
		box = "[x]"
	}
	return LayoutString(c.label, box, width)
}

func (c *Checkbox) MinWidth() int {
	return len(c.label) + 4
}

func LayoutString(lhs, rhs string, width int) string {
	if len(lhs)+1+len(rhs) > width {
		lhs = strings.TrimSpace(lhs[:width-len(rhs)-5]) + "..."
	}
	return lhs + strings.Repeat(" ", width-len(rhs)-1-len(lhs)) + rhs
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
