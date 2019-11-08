package main

import (
	"fmt"
	"os"

	"github.com/nsf/termbox-go"
)

var tracks []*Track

func main() {
	err := termbox.Init()
	must(err)
	defer termbox.Close()

	tracks = append(tracks, NewTrack(8))
	tracks = append(tracks, NewTrack(4))

	var page int
	render(page)

	var done bool
	for !done {
		e := termbox.PollEvent()
		switch e.Type {
		case termbox.EventKey:
			switch e.Key {
			case termbox.KeyEsc:
				done = true

			case termbox.KeyPgup:
				if page > 0 {
					page--
					render(page)
				}

			case termbox.KeyPgdn:
				if page < 0xFF {
					page++
					render(page)
				}
			}
		}
	}
}

func must(err error) {
	if err != nil {
		fmt.Printf("trk: %v\n", err)
		os.Exit(2)
	}
}

func write(x, y int, s string) {
	for i, c := range s {
		termbox.SetCell(x+i, y, c, termbox.ColorDefault, termbox.ColorDefault)
	}
}

func render(page int) {
	for i := 0; i < 16; i++ {
		var (
			y    = i
			step = 16*page + i
		)

		write(0, y, fmt.Sprintf("%03X", step))

		x := 4
		for _, t := range tracks {
			var pat string
			if p, ch := t.Pattern(y); ch || step%16 == 0 {
				pat = Pattern(p).String()
			} else {
				pat = "---"
			}
			write(x, y, pat)
			x += 4
			for v := 0; v < t.VoiceCount(); v++ {
				m, ch := t.Muted(y, v)
				var str string
				if ch || step%16 == 0 {
					if m {
						str = "\u25c7"
					} else {
						str = "\u25c6"
					}
				} else {
					str = "-"
				}
				write(x, y, str)
				x++
			}
			x++
		}
	}
	termbox.Flush()
}
