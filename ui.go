package main

import "github.com/nsf/termbox-go"

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
