package main

import (
	"fmt"
	"os"

	"github.com/nsf/termbox-go"
)

var drv midiDriver
var lp *Launchpad
var view LaunchpadView
var model *Model

func main() {
	quit := make(chan struct{})

	err := termbox.Init()
	must(err)
	defer termbox.Close()

	lp, err = ConnectLaunchpad()
	must(err)
	defer lp.Close()
	lp.Reset()

	model = NewModel()
	view = NewLaunchpadSessionView()

	render()

	go func() {
		for {
			select {
			case m := <-lp.In():
				view.Handle(lp, model, m)
				render()

			case <-quit:
				return
			}
		}
	}()

	var done bool
	for !done {
		e := termbox.PollEvent()
		switch e.Type {
		case termbox.EventKey:
			switch e.Key {
			case termbox.KeyEsc:
				done = true

			case termbox.KeyPgup:
				model.DecPage()
				render()

			case termbox.KeyPgdn:
				model.IncPage()
				render()
			}
		}
	}
	close(quit)
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
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

func render() {
	view.Update(model)

	var (
		page = 1 + model.Page()
		bar  = fmt.Sprintf("%d:4", 1+model.Cursor()/16)
		step = 1 + (model.Cursor() % 16)
	)

	write(0, 0, fmt.Sprintf("SEQ %03d PAGE %s TRIG %02d", page, bar, step))
	termbox.Flush()

	view.Render(lp, model)
	lp.Flush()
}
