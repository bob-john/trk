package ui

import (
	"errors"
	"strings"
	"trk/rtmididrv"

	"github.com/nsf/termbox-go"
)

var ErrCanceled = errors.New("ui: canceled")

func init() {
	must(termbox.Init())
}

// Out queries the user to pick an output port for the named device.
func Out(name string) (port string) {
	Clear()
	drv, err := rtmididrv.New()
	must(err)
	outs, err := drv.Outs()
	must(err)
	p := NewOptionPage(name + " Output")
	for i, port := range outs {
		if strings.Contains(port.String(), name) {
			p.selected = i
		}
		p.Label(port.String())
	}
	d := NewDialog(p)

	d.Render()
	termbox.Flush()

	for {
		e := termbox.PollEvent()
		switch e.Type {
		case termbox.EventInterrupt:
			panic(ErrCanceled)

		case termbox.EventKey:
			switch e.Key {
			case termbox.KeyArrowUp, termbox.KeyArrowDown:
				d.Handle(e)

			case termbox.KeyEnter:
				return d.Page().Item().Value()

			case termbox.KeyEsc:
				panic(ErrCanceled)
			}
		}
		d.Render()
		termbox.Flush()
	}
}

func Clear() {
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
}

func Close() {
	Clear()
	// termbox.Interrupt()
	termbox.Close()
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
