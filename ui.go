package main

import (
	"github.com/nsf/termbox-go"
)

type UI struct {
	views []View
}

func (ui *UI) Clear() {
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
	ui.views = nil
}

func (ui *UI) Print(x, y int, text string, fg, bg termbox.Attribute, onClick func(x, y int)) {
	SetString(x, y, text, fg, bg)
	ui.views = append(ui.views, View{x, y, len(text), 1, onClick})
}

func (ui *UI) Click(x, y int) {
	for _, v := range ui.views {
		if v.Hit(x, y) && v.OnClick != nil {
			v.OnClick(x-v.X, y-v.Y)
		}
	}
}

func (ui *UI) Flush() {
	termbox.Flush()
}

type View struct {
	X, Y          int
	Width, Height int
	OnClick       func(x, y int)
}

func (v View) Hit(x, y int) bool {
	return x >= v.X && y >= v.Y && x <= v.X+v.Width && y <= v.Y+v.Height
}

type Box struct {
	X, Y, Width, Height int
}

func (b Box) Left() int {
	return b.X
}

func (b Box) Right() int {
	return b.X + b.Width
}

func (b Box) Top() int {
	return b.Y
}

func (b Box) Bottom() int {
	return b.Y + b.Height
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
	Items []string
}

func (d *Dialog) Render() {
	for n, item := range d.Items {
		SetString(d.Left()+1, d.Top()+1+n, item, termbox.ColorDefault, termbox.ColorDefault)
	}
	d.Box.Render()
}
