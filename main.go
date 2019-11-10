package main

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"

	"github.com/gomidi/midi/midimessage/realtime"
	"github.com/nsf/termbox-go"
)

var (
	currentStep        = 0
	editing            = false
	editor             = &LineEditor{}
	seq                = &Seq{}
	digitakt, digitone *Device
)

func main() {
	err := termbox.Init()
	must(err)
	defer termbox.Close()

	user, err := user.Current()
	must(err)
	home := user.HomeDir

	appDir := filepath.Join(home, ".trk")
	err = os.MkdirAll(appDir, 0700)
	must(err)
	tmpFilePath := filepath.Join(appDir, "tmp.trk")

	seq.ReadFile(tmpFilePath)

	digitakt, _ = ConnectDevice("Elektron Digitakt", "Elektron Digitakt")
	digitone, _ = ConnectDevice("Elektron Digitone", "Elektron Digitone")

	var (
		quit       = make(chan struct{})
		transportC = make(chan realtime.Message)
		renderC    = make(chan struct{})
		playing    = false
		tick       = 0
	)
	if digitakt != nil {
		defer digitakt.Close()
		go listenTransport(digitakt, transportC, quit)
	}
	if digitone != nil {
		defer digitone.Close()
	}
	go func() {
		for {
			select {
			case <-renderC:
				render()
			case <-quit:
				return
			}
		}
	}()
	go func() {
		for {
			select {
			case m := <-transportC:
				switch m {
				case realtime.TimingClock:
					if playing && !editing {
						tick++
						if tick == 12 {
							seq.Play(currentStep+1, digitakt, digitone)
						} else if tick == 24 {
							currentStep++
							tick = 0
							renderC <- struct{}{}
						}
					}

				case realtime.Start:
					playing = true
					tick = 0

				case realtime.Continue:
					playing = true

				case realtime.Stop:
					playing = false
				}

			case <-quit:
				return
			}
		}
	}()

	renderC <- struct{}{}

	seq.Play(currentStep, digitakt, digitone)

	var done bool
	for !done {
		e := termbox.PollEvent()
		switch e.Type {
		case termbox.EventKey:
			switch e.Key {
			case termbox.KeyEsc:
				if editing {
					editing = false
				} else {
					done = true
				}

			case termbox.KeyArrowUp:
				if editing {
					editor.ActiveCell().Inc()
				} else if currentStep > 0 {
					currentStep--
					seq.Play(currentStep, digitakt, digitone)
				}

			case termbox.KeyArrowDown:
				if editing {
					editor.ActiveCell().Dec()
				} else if currentStep < 0xfff {
					currentStep++
					seq.Play(currentStep, digitakt, digitone)
				}

			case termbox.KeyDelete, termbox.KeyBackspace:
				if editing {
					editor.ActiveCell().Clear()
				} else {
					seq.Insert(seq.emptyLine(currentStep))
					seq.WriteFile(tmpFilePath)
					seq.Play(currentStep, digitakt, digitone)
				}

			case termbox.KeyArrowLeft:
				if editing {
					editor.MoveToPreviousCell()
				}

			case termbox.KeyArrowRight:
				if editing {
					editor.MoveToNextCell()
				}

			case termbox.KeyEnter:
				editing = !editing
				if editing {
					editor.Reset(seq.Line(currentStep), seq.ConsolidatedLine(currentStep))
				} else {
					seq.Insert(editor.Line())
					seq.WriteFile(tmpFilePath)
					seq.Play(currentStep, digitakt, digitone)
				}

			case termbox.KeyPgup:
				if editing {
					editor.ActiveCell().PageInc()
				} else {
					currentStep -= 16
					if currentStep < 0 {
						currentStep = 0
					}
					seq.Play(currentStep, digitakt, digitone)
				}

			case termbox.KeyPgdn:
				if editing {
					editor.ActiveCell().PageDec()
				} else {
					currentStep += 16
					if currentStep > 0xfff {
						currentStep = 0xfff
					}
					seq.Play(currentStep, digitakt, digitone)
				}
			}
		}
		renderC <- struct{}{}
	}
	close(quit)
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
}

func listenTransport(synth *Device, transportC chan<- realtime.Message, quit <-chan struct{}) {
	for {
		select {
		case m := <-synth.In():
			if m, ok := m.(realtime.Message); ok {
				transportC <- m
			}

		case <-quit:
			return
		}
	}
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
		line := seq.Line(step)
		if i == 8 && editing {
			line = editor.Line()
		}
		SetString(0, i, line, fg, bg)
		if i == 8 && editing {
			cell := editor.ActiveCell()
			SetString(cell.Index(), i, cell.String(), fg|termbox.AttrReverse, bg)
		}
	}
	SetString(30, 0, fmt.Sprintf("DT: %v", digitakt != nil), termbox.ColorDefault, termbox.ColorDefault)
	SetString(30, 1, fmt.Sprintf("DN: %v", digitone != nil), termbox.ColorDefault, termbox.ColorDefault)
	termbox.Flush()
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
