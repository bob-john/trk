package ui

import (
	"errors"
	"strings"
	"trk/rtmididrv"
	"trk/tracker"

	"github.com/nsf/termbox-go"
	"gitlab.com/gomidi/midi/mid"
)

var ErrCanceled = errors.New("ui: canceled")
var midiDriver mid.Driver

func init() {
	var err error
	midiDriver, err = rtmididrv.New()
	must(err)
	must(termbox.Init())
}

func Input(name string) *tracker.In {
	ins, err := midiDriver.Ins()
	must(err)
	var ports []string
	for _, port := range ins {
		ports = append(ports, port.String())
	}
	port := Port(name, "Input", ports)
	return tracker.OpenIn(port)
}

func Output(name string) (port string) {
	outs, err := midiDriver.Outs()
	must(err)
	var ports []string
	for _, port := range outs {
		ports = append(ports, port.String())
	}
	return Port(name, "Output", ports)
}

func Port(name, suffix string, ports []string) (port string) {
	Clear()
	p := NewOptionPage(name + " " + suffix)
	for i, port := range ports {
		if strings.Contains(port, name) {
			p.selected = i
		}
		p.Label(port)
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
