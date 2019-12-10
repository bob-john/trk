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
	if IsKey(e, termbox.KeyEsc, termbox.KeyCtrlO) {
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
	stack []*OptionPage
}

func NewDialog(model *OptionPage) *Dialog {
	return &Dialog{[]*OptionPage{model}}
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

func (d *Dialog) Handle(ui *UI, e termbox.Event) (handled bool) {
	if e.Type != termbox.EventKey {
		return
	}
	handled = d.Page().Handle(d, e)
	if handled {
		return
	}
	switch e.Key {
	case termbox.KeyEsc:
		ui.Dismiss()
		handled = true
	case termbox.KeyBackspace, termbox.KeyBackspace2, termbox.KeyArrowLeft:
		if len(d.stack) > 1 {
			d.Back()
		} else {
			ui.Dismiss()
		}
		handled = true
	}
	return
}

func (d *Dialog) Render() {
	var (
		page  = d.Page()
		title = page.title
		sz    = page.PreferredSize().Union(MakeSize(len(title)+4, 0))
		w, h  = termbox.Size()
		x, y  = (w - sz.Width) / 2, (h - sz.Height) / 2
	)
	DrawBox(x, y, x+sz.Width, y+sz.Height, termbox.ColorDefault, termbox.ColorDefault)
	DrawString(x+1, y, " "+title+" ", termbox.ColorDefault, termbox.ColorDefault)
	page.Render(x+1, y+1, sz.Width, sz.Height)
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

func (p *OptionPage) Page(title string, build func(page *OptionPage)) {
	page := &OptionPage{title, nil, 0, 0}
	p.items = append(p.items, &Menu{title, page})
	build(page)
}

func (p *OptionPage) Checkbox(title string, on bool, onchange func(bool)) {
	p.items = append(p.items, &Checkbox{title, on, onchange})
}

func (p *OptionPage) Picker(label string, values map[int]string, selected int, onchange func(int)) {
	p.items = append(p.items, &Picker{label, values, selected, onchange})
}

func (p *OptionPage) Label(label string) {
	p.items = append(p.items, &Label{label})
}

func (p *OptionPage) Button(label string, onpush func()) {
	p.items = append(p.items, &Button{label, onpush})
}

func (p *OptionPage) PreferredSize() Size {
	sz := MakeSize(len(p.title)+4, len(p.items)+1)
	for _, item := range p.items {
		minWidth := item.MinWidth() + 2
		if sz.Width < minWidth {
			sz.Width = minWidth
		}
	}
	_, h := termbox.Size()
	if sz.Height > h-2 {
		sz.Height = h - 2
	}
	return sz
}

func (p *OptionPage) OnEnter() {
	p.selected = 0
	p.offset = 0
}

func (p *OptionPage) Handle(d *Dialog, e termbox.Event) (handled bool) {
	if e.Type != termbox.EventKey {
		return
	}
	switch e.Key {
	case termbox.KeyArrowDown:
		p.selected++
		if p.selected >= p.offset+5 {
			p.offset++
		}
		handled = true
	case termbox.KeyArrowUp:
		p.selected--
		if p.selected < p.offset {
			p.offset--
		}
		handled = true
	default:
		handled = p.items[p.selected].Handle(d, e)
	}
	p.selected = Clamp(p.selected, 0, len(p.items)-1)
	if len(p.items) > 5 {
		p.offset = Clamp(p.offset, 0, len(p.items)-5)
	} else {
		p.offset = 0
	}
	return
}

func (p *OptionPage) Render(x, y, width, height int) {
	for n := 0; n < height; n++ {
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
	if len(p.items) > height {
		DrawString(x+width-1, y+height*p.selected/len(p.items), "\u2590", termbox.ColorDefault, termbox.ColorDefault)
	}
}

type OptionItem interface {
	Handle(*Dialog, termbox.Event) bool
	String(int) string
	MinWidth() int
}

type Menu struct {
	label string
	page  *OptionPage
}

func (m *Menu) Handle(dialog *Dialog, e termbox.Event) bool {
	if IsKey(e, termbox.KeyArrowRight, termbox.KeyEnter) {
		dialog.Enter(m.page)
		return true
	}
	return false
}

func (m *Menu) String(width int) string {
	return LayoutString(m.label, ">", width)
}

func (m *Menu) MinWidth() int {
	return len(m.label) + 2
}

type Checkbox struct {
	label    string
	on       bool
	onchange func(bool)
}

func (c *Checkbox) Handle(dialog *Dialog, e termbox.Event) bool {
	if IsKey(e, termbox.KeyArrowRight, termbox.KeyEnter) {
		c.on = !c.on
		c.onchange(c.on)
		return true
	}
	return false
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
	values   map[int]string
	selected int
	onchange func(int)
}

func (p *Picker) Handle(dialog *Dialog, e termbox.Event) (handled bool) {
	if e.Type != termbox.EventKey {
		return
	}
	switch e.Key {
	case termbox.KeyArrowRight:
		if _, ok := p.values[p.selected+1]; ok {
			p.selected++
		}
		handled = true
	case termbox.KeyArrowLeft:
		if _, ok := p.values[p.selected-1]; ok {
			p.selected--
		}
		handled = true
	}
	p.onchange(p.selected)
	return
}

func (p *Picker) String(width int) string {
	return LayoutString(p.label, p.values[p.selected], width)
}

func (p *Picker) MinWidth() int {
	w := 0
	for _, str := range p.values {
		if len(str) > w {
			w = len(str)
		}
	}
	return len(p.label) + 1 + w
}

type Label struct {
	text string
}

func (*Label) Handle(dialog *Dialog, e termbox.Event) bool {
	return false
}

func (l *Label) String(width int) string {
	return LayoutString(l.text, "", width)
}

func (l *Label) MinWidth() int {
	return len(l.text)
}

type Button struct {
	label  string
	onpush func()
}

func (b *Button) Handle(dialog *Dialog, event termbox.Event) bool {
	if IsKey(event, termbox.KeyEnter) {
		b.onpush()
		return true
	}
	return false
}

func (b *Button) String(width int) string {
	return b.label
}

func (b *Button) MinWidth() int {
	return len(b.label)
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
