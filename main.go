package main

import (
	"fmt"

	"github.com/nsf/termbox-go"
)

var (
	currentStep = 0
	editing     = false
	currentCell = 0
	editor      = &LineEditor{}
)

func main() {
	err := termbox.Init()
	must(err)
	defer termbox.Close()

	// SetString(0, 0, "01", termbox.ColorDefault, termbox.ColorDefault)
	// SetString(3, 0, "A01", termbox.ColorDefault, termbox.ColorDefault)
	// SetString(7, 0, strings.Repeat("\u258E", 8), termbox.ColorDefault, termbox.ColorDefault)
	// SetString(16, 0, "A01", termbox.ColorDefault, termbox.ColorDefault)
	// SetString(20, 0, strings.Repeat("\u258E", 4), termbox.ColorDefault, termbox.ColorDefault)
	// SetString(0, 1, "02", termbox.ColorDefault, termbox.ColorDefault)
	// SetString(3, 1, strings.Repeat(".", 3), termbox.ColorDefault, termbox.ColorDefault|termbox.AttrReverse)
	// SetString(7, 1, strings.Repeat(".", 8), termbox.ColorDefault, termbox.ColorDefault|termbox.AttrReverse)
	// SetString(16, 1, strings.Repeat(".", 3), termbox.ColorDefault, termbox.ColorDefault|termbox.AttrReverse)
	// SetString(20, 1, strings.Repeat("\u258e", 4), termbox.ColorDefault, termbox.ColorDefault|termbox.AttrReverse)

	// SetString(0, 1, fmt.Sprintf("%02d %s %s", 2+step, pc.String(), muted.String()), termbox.ColorDefault, termbox.ColorRed)
	// SetString(0, 2, fmt.Sprintf("%02d %s %s", 3+step, pc.String(), muted.String()), termbox.ColorDefault, termbox.ColorDefault)

	render()

	var done bool
	for !done {
		e := termbox.PollEvent()
		switch e.Type {
		case termbox.EventKey:
			switch e.Key {
			case termbox.KeyEsc:
				done = true

			case termbox.KeyArrowUp:
				if editing {
					editor.Cell(currentCell).Inc()
				} else if currentStep > 0 {
					currentStep--
				}

			case termbox.KeyArrowDown:
				if editing {
					editor.Cell(currentCell).Dec()
				} else if currentStep < 0xfff {
					currentStep++
				}

			case termbox.KeyDelete, termbox.KeyBackspace:
				if editing {
					editor.Cell(currentCell).Clear()
				}

			case termbox.KeyArrowLeft:
				if currentCell > 0 {
					currentCell--
				}

			case termbox.KeyArrowRight:
				if currentCell < editor.CellCount()-1 {
					currentCell++
				}

			case termbox.KeyEnter:
				editing = !editing
				currentCell = 0
				if editing {
					editor.Reset(line(currentStep), "*** A01 ++++++++ A01 ++++")
				}
			}
		}
		render()
	}
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
}

func render() {
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
	for i := 0; i < 16; i++ {
		step := currentStep - 8 + i
		if step < 0 || step > 0xfff {
			continue
		}
		fg, bg := termbox.ColorBlue, termbox.ColorDefault
		if (step/16)%2 == 1 {
			fg = termbox.ColorGreen
		}
		if i == 8 && !editing {
			fg = fg | termbox.AttrReverse
		}
		line := line(step)
		if i == 8 && editing {
			line = editor.Line()
		}
		SetString(0, i, line, fg, bg)
		if i == 8 && editing {
			cell := editor.Cell(currentCell)
			SetString(cell.Index(), i, cell.String(), fg|termbox.AttrReverse, bg)
		}
	}
	termbox.Flush()
}

func line(step int) string {
	return fmt.Sprintf("%03X ... ........ ... ....", step)
}

// var drv midiDriver
// var lp *Launchpad
// var view LaunchpadView
// var model *Model

// func main() {
// 	quit := make(chan struct{})

// 	err := termbox.Init()
// 	must(err)
// 	defer termbox.Close()

// 	lp, err = ConnectLaunchpad()
// 	must(err)
// 	defer lp.Close()
// 	lp.Reset()

// 	model = NewModel()
// 	view = &LaunchpadMainView{}

// 	render()

// 	go func() {
// 		for {
// 			select {
// 			case m := <-lp.In():
// 				view.Handle(lp, model, m)
// 				render()

// 			case <-quit:
// 				return
// 			}
// 		}
// 	}()

// var done bool
// for !done {
// 	e := termbox.PollEvent()
// 	switch e.Type {
// 	case termbox.EventKey:
// 		switch e.Key {
// 		case termbox.KeyEsc:
// 			done = true

// 		case termbox.KeyPgup:
// 			model.DecPage()
// 			render()

// 		case termbox.KeyPgdn:
// 			model.IncPage()
// 			render()
// 		}
// 	}
// }
// close(quit)
// termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
// }

// func write(x, y int, s string) {
// 	for i, c := range s {
// 		termbox.SetCell(x+i, y, c, termbox.ColorDefault, termbox.ColorDefault)
// 	}
// }

// func render() {
// 	view.Update(model)

// 	var (
// 		page = 1 + model.Page()
// 		bar  = fmt.Sprintf("%d:4", 1+model.Cursor()/16)
// 		step = 1 + (model.Cursor() % 16)
// 	)

// 	write(0, 0, fmt.Sprintf("SEQ %03d PAGE %s TRIG %02d", page, bar, step))
// 	termbox.Flush()

// 	view.Render(lp, model)
// 	lp.Flush()
// }
