package ui

import "github.com/nsf/termbox-go"

func IsKey(e termbox.Event, keys ...termbox.Key) bool {
	if e.Type != termbox.EventKey {
		return false
	}
	for _, key := range keys {
		if key == e.Key {
			return true
		}
	}
	return false
}

func DrawString(x, y int, s string, fg, bg termbox.Attribute) {
	for _, c := range s {
		switch c {
		case '\n':
			x = 0
			y++
		default:
			termbox.SetCell(x, y, c, fg, bg)
			x++
		}
	}
}

func DrawBox(left, top, right, bottom int, fg, bg termbox.Attribute) {
	if left > right {
		left, right = right, left
	}
	if top > bottom {
		top, bottom = bottom, top
	}
	for x := left; x < right; x++ {
		termbox.SetCell(x, top, '─', fg, bg)
		termbox.SetCell(x, bottom, '─', fg, bg)
	}
	for y := top; y < bottom; y++ {
		termbox.SetCell(left, y, '│', fg, bg)
		termbox.SetCell(right, y, '│', fg, bg)
	}
	termbox.SetCell(left, top, '┌', fg, bg)
	termbox.SetCell(right, top, '┐', fg, bg)
	termbox.SetCell(left, bottom, '└', fg, bg)
	termbox.SetCell(right, bottom, '┘', fg, bg)
}

type Writer struct {
	X, Y   int
	Fg, Bg termbox.Attribute
}

func (w *Writer) Write(p []byte) (n int, err error) {
	for _, c := range string(p) {
		switch c {
		case '\n':
			w.X = 0
			w.Y++
		default:
			termbox.SetCell(w.X, w.Y, c, w.Fg, w.Bg)
			w.X++
		}
	}
	return len(p), nil
}
