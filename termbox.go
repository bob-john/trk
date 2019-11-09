package main

import "github.com/nsf/termbox-go"

func SetString(x, y int, s string, fg, bg termbox.Attribute) {
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
