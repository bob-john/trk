package main

import (
	"github.com/nsf/termbox-go"
)

type Box struct {
	x, y, width, height int
}

func MakeBox(x, y, width, height int) Box {
	if width < 0 {
		x, width = x+width, -width
	}
	if height < 0 {
		y, height = y+height, -height
	}
	return Box{x, y, width, height}
}

func (b Box) Left() int {
	return b.x
}

func (b Box) Right() int {
	return b.x + b.width
}

func (b Box) Top() int {
	return b.y
}

func (b Box) Bottom() int {
	return b.y + b.height
}

func (b Box) Width() int {
	return b.width
}

func (b Box) Height() int {
	return b.height
}

func (b Box) Render() {
	x0, x1, y0, y1 := b.Left(), b.Right(), b.Top(), b.Bottom()
	fg, bg := termbox.ColorDefault, termbox.ColorDefault
	for x := x0; x < x1; x++ {
		termbox.SetCell(x, y0, '─', fg, bg)
		termbox.SetCell(x, y1, '─', fg, bg)
	}
	for y := y0; y < y1; y++ {
		termbox.SetCell(x0, y, '│', fg, bg)
		termbox.SetCell(x1, y, '│', fg, bg)
	}
	termbox.SetCell(x0, y0, '┌', fg, bg)
	termbox.SetCell(x1, y0, '┐', fg, bg)
	termbox.SetCell(x0, y1, '└', fg, bg)
	termbox.SetCell(x1, y1, '┘', fg, bg)
}

type Dialog struct {
	Box
	model        *Settings
	selectedItem int
}

func NewDialog(x, y, width, height int, model *Settings) *Dialog {
	return &Dialog{MakeBox(x, y, width, height), model, 0}
}

func (d *Dialog) Handle(e termbox.Event) bool {
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

	default:
		d.model.Items[d.selectedItem].Handle(e)
	}
	d.selectedItem = clamp(d.selectedItem, 0, len(d.model.Items)-1)
	return false
}

func (d *Dialog) Render() {
	d.Box.Render()
	SetString(d.Left()+1, d.Top(), " "+d.model.Title+" ", termbox.ColorDefault, termbox.ColorDefault)
	for n, item := range d.model.Items {
		fg, bg := termbox.ColorDefault, termbox.ColorDefault
		if n == d.selectedItem {
			fg = fg | termbox.AttrReverse
		}
		SetString(d.Left()+1, d.Top()+1+n, item.String(d.Width()), fg, bg)
	}
}
