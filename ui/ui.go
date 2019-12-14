package ui

import (
	"errors"
	"strings"
	"trk/rtmididrv"
	"trk/tracker"

	"github.com/nsf/termbox-go"
)

var ErrCanceled = errors.New("ui: canceled")

type UI struct{}

func New() (*UI, error) {
	err := termbox.Init()
	if err != nil {
		return nil, err
	}
	return &UI{}, nil
}

// Output queries the user to pick an output port for the named device.
func (ui *UI) Out(tracker *tracker.Tracker, name string) (out *tracker.Out, err error) {
	ui.Clear()
	drv, err := rtmididrv.New()
	if err != nil {
		return
	}
	outs, err := drv.Outs()
	if err != nil {
		return
	}
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
			return nil, ErrCanceled

		case termbox.EventKey:
			switch e.Key {
			case termbox.KeyArrowUp, termbox.KeyArrowDown:
				d.Handle(ui, e)

			case termbox.KeyEnter:
				return tracker.Out(d.Page().Item().Value()), nil

			case termbox.KeyEsc:
				return nil, ErrCanceled
			}
		}
		d.Render()
		termbox.Flush()
	}
}

func (ui *UI) Clear() {
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
}

func (ui *UI) Close() {
	ui.Clear()
	// termbox.Interrupt()
	termbox.Close()
}
