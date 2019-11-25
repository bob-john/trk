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

func (d *Dialog) Enter(page *OptionPage) {
	d.stack = append(d.stack, page)
	page.OnEnter()
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
	case termbox.KeyBackspace:
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
	title := page.title
	sz := page.PreferredSize().Union(MakeSize(len(title)+4, 0))
	DrawBox(d.x, d.y, d.x+sz.Width, d.y+sz.Height, termbox.ColorDefault, termbox.ColorDefault)
	DrawString(d.x+1, d.y, " "+title+" ", termbox.ColorDefault, termbox.ColorDefault)
	page.Render(d.x+1, d.y+1, sz.Width)
}

type OptionPage struct {
	title    string
	items    []OptionItem
	selected int
	offset   int
}

func NewOptionPage(title string) *OptionPage {
	return &OptionPage{title, nil, 0, 0}
}

func (p *OptionPage) AddMenu(title string, build func(page *OptionPage)) {
	page := &OptionPage{title, nil, 0, 0}
	p.items = append(p.items, &Menu{title, page})
	build(page)
}

func (p *OptionPage) AddCheckbox(title string, on bool, onchange func(bool)) {
	p.items = append(p.items, &Checkbox{title, on, onchange})
}

func (p *OptionPage) AddPicker(label string, values []string, selected int, onchange func(int)) {
	p.items = append(p.items, &Picker{label, values, selected, onchange})
}

func (p *OptionPage) AddLabel(label string) {
	p.items = append(p.items, &Label{label})
}

func (p *OptionPage) PreferredSize() Size {
	sz := MakeSize(len(p.title)+4, len(p.items)+1)
	for _, item := range p.items {
		minWidth := item.MinWidth() + 2
		if sz.Width < minWidth {
			sz.Width = minWidth
		}
	}
	if sz.Height > 6 {
		sz.Height = 6
	}
	return sz
}

func (p *OptionPage) OnEnter() {
	p.selected = 0
	p.offset = 0
}

func (p *OptionPage) Handle(d *Dialog, e termbox.Event) {
	if e.Type != termbox.EventKey {
		return
	}
	switch e.Key {
	case termbox.KeyArrowDown:
		p.selected++
		if p.selected >= p.offset+5 {
			p.offset++
		}
	case termbox.KeyArrowUp:
		p.selected--
		if p.selected < p.offset {
			p.offset--
		}
	default:
		p.items[p.selected].Handle(d, e)
	}
	p.selected = clamp(p.selected, 0, len(p.items)-1)
	if len(p.items) > 5 {
		p.offset = clamp(p.offset, 0, len(p.items)-5)
	} else {
		p.offset = 0
	}
}

func (p *OptionPage) Render(x, y, width int) {
	for n := 0; n < 5; n++ {
		if p.offset+n >= len(p.items) {
			return
		}
		item := p.items[p.offset+n]
		fg, bg := termbox.ColorDefault, termbox.ColorDefault
		if p.offset+n == p.selected {
			fg = fg | termbox.AttrReverse
		}
		DrawString(x, y+n, item.String(width), fg, bg)
	}
	if len(p.items) > 5 {
		DrawString(x+width-1, y+5*p.selected/len(p.items), "\u2590", termbox.ColorDefault, termbox.ColorDefault)
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
	label    string
	on       bool
	onchange func(bool)
}

func (c *Checkbox) Handle(dialog *Dialog, e termbox.Event) {
	if IsKey(e, termbox.KeyArrowRight, termbox.KeyEnter) {
		c.on = !c.on
		c.onchange(c.on)
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

type Picker struct {
	label    string
	values   []string
	selected int
	onchange func(int)
}

func (p *Picker) Handle(dialog *Dialog, e termbox.Event) {
	if e.Type != termbox.EventKey {
		return
	}
	switch e.Key {
	case termbox.KeyEnter:
		p.selected = (p.selected + 1) % len(p.values)
	case termbox.KeyArrowRight:
		p.selected++
	case termbox.KeyArrowLeft:
		p.selected--
	}
	p.selected = clamp(p.selected, 0, len(p.values)-1)
	p.onchange(p.selected)
}

func (p *Picker) String(width int) string {
	return LayoutString(p.label, p.values[p.selected], width)
}

func (p *Picker) MinWidth() int {
	w := 0
	for _, val := range p.values {
		if len(val) > w {
			w = len(val)
		}
	}
	return len(p.label) + 1 + w
}

type Label struct {
	text string
}

func (*Label) Handle(dialog *Dialog, e termbox.Event) {}

func (l *Label) String(width int) string {
	return LayoutString(l.text, "", width)
}

func (l *Label) MinWidth() int {
	return len(l.text)
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
