package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gomidi/midi/midimessage/realtime"
	"github.com/nsf/termbox-go"
)

const (
	pageSize = 16
)

var (
	ui       = new(UI)
	model    = new(Model)
	digitakt *Device
	digitone *Device
)

func main() {
	var err error

	if len(os.Args) != 2 {
		fmt.Println("usage: trk <path> [<device> ...]")
		fmt.Println("trk: invalid command line")
		os.Exit(1)
	}

	err = model.LoadSeq(os.Args[1])
	must(err)

	err = termbox.Init()
	must(err)
	defer termbox.Close()
	termbox.SetInputMode(termbox.InputCurrent | termbox.InputMouse)

	digitakt, _ = OpenDevice("Digitakt", "Elektron Digitakt", "Elektron Digitakt")
	digitone, _ = OpenDevice("Digitone", "Elektron Digitone", "Elektron Digitone")

	model.Seq.ConsolidatedRow(0).Play(digitone, digitakt)

	var (
		eventC = make(chan termbox.Event)
		midiC  = make(chan Message)
	)

	if digitakt != nil {
		go func() {
			for m := range digitakt.In() {
				midiC <- Message{m, digitakt}
			}
		}()
	}
	if digitone != nil {
		go func() {
			for m := range digitone.In() {
				midiC <- Message{m, digitone}
			}
		}()
	}

	go func() {
		for {
			e := termbox.PollEvent()
			if e.Type == termbox.EventInterrupt {
				break
			}
			eventC <- e
		}
	}()

	var (
		done bool
		tick int
	)
	for !done {
		var (
			oldHead = model.Head
		)
		select {
		case e := <-eventC:
			switch e.Type {
			case termbox.EventKey:
				switch e.Key {
				case termbox.KeyCtrlS:
					err := model.Seq.Write(os.Args[1])
					if err != nil {
						log.Fatal(err)
					}
				case termbox.KeyCtrlX:
					err := model.Seq.Write(os.Args[1])
					if err != nil {
						log.Fatal(err)
					}
					done = true

				case termbox.KeyArrowRight:
					model.SetTrig(model.Trig() + 1)
				case termbox.KeyArrowLeft:
					model.SetTrig(model.Trig() - 1)
				case termbox.KeyArrowUp:
					model.SetTrig(model.Trig() - 8)
				case termbox.KeyArrowDown:
					model.SetTrig(model.Trig() + 8)
				case termbox.KeyPgup:
					model.SetPattern(model.Pattern() - 1)
				case termbox.KeyPgdn:
					model.SetPattern(model.Pattern() + 1)
				case termbox.KeyHome:
					model.SetPattern(0)
				case termbox.KeyEnd:
					model.SetPattern(model.LastPattern())

				case termbox.KeyDelete, termbox.KeyBackspace:
					model.ClearStep()
				case termbox.KeyEnter:
					model.ToggleRecording()
				}

			case termbox.EventMouse:
				switch e.Key {
				case termbox.MouseWheelUp:
					model.SetPattern(model.Pattern() + 1)
				case termbox.MouseWheelDown:
					model.SetPattern(model.Pattern() - 1)

				case termbox.MouseLeft:
					ui.Click(e.MouseX, e.MouseY)
				}
			}

		case m := <-midiC:
			switch m.Message {
			case realtime.TimingClock:
				if model.State == Playing {
					tick++
				}
			case realtime.Start:
				model.State = Playing
				tick = 0
			case realtime.Continue:
				model.State = Playing
			case realtime.Stop:
				model.State = Viewing
			}
			if model.State == Recording {
				model.Seq.Insert(m.Device.Name(), model.Head, m.Message)
			}
		}
		if model.State == Playing {
			switch tick {
			case 12:
				row := model.Seq.ConsolidatedRow(model.Head + 2)
				row.Digitone.Pattern.Play(digitone, 15)
				row.Digitakt.Pattern.Play(digitakt, 15)

			case 18:
				row := model.Seq.ConsolidatedRow(model.Head + 1)
				row.Digitone.Mute.Play(digitone, row.Digitone.Channels)
				row.Digitakt.Mute.Play(digitakt, row.Digitakt.Channels)
				row.Digitone.Mute.Play(digitakt, row.Digitone.Channels) //HACK

			case 24:
				model.Head++
				tick = 0
			}
		} else {
			if model.Head != oldHead {
				model.Seq.ConsolidatedRow(model.Head).Play(digitone, digitakt)
			}
		}
		render()
	}
	termbox.Interrupt()
}

func clamp(val, min, max int) int {
	if val < min {
		return min
	}
	if val > max {
		return max
	}
	return val
}

func color(on, ch bool) (fg termbox.Attribute) {
	if ch {
		fg = termbox.ColorRed
	}
	if !on {
		return
	}
	switch model.State {
	case Playing:
		fg = termbox.ColorGreen
	case Recording:
		fg = termbox.ColorRed
	}
	fg = fg | termbox.AttrReverse
	return
}

func render() {
	ui.Clear()
	row, org := model.Seq.ConsolidatedRow(model.Head), model.Seq.Row(model.Head)
	SetString(4, 0, row.Digitakt.Pattern.String(), color(false, org.Digitakt.Pattern != -1), termbox.ColorDefault)
	SetString(8, 0, row.Digitakt.Mute.Format(row.Digitakt.Channels), color(false, len(org.Digitakt.Mute) != 0), termbox.ColorDefault)
	SetString(8+row.Digitakt.Channels.Len+1, 0, row.Digitone.Pattern.String(), color(false, org.Digitone.Pattern != -1), termbox.ColorDefault)
	SetString(12+row.Digitakt.Channels.Len+1, 0, row.Digitone.Mute.Format(row.Digitone.Channels), color(false, len(org.Digitone.Mute) != 0), termbox.ColorDefault)
	SetString(0, 2, fmt.Sprintf("%03d", 1+model.Pattern()), termbox.ColorDefault, termbox.ColorDefault)
	for n := 0; n < 16; n++ {
		n := n
		ch := model.HeadForTrig(n) == 0 || model.Seq.Row(model.HeadForTrig(n)).HasChanges(model.Seq.Row(model.HeadForTrig(n-1)))
		ui.Print(4+(n%8)*3, 2+3*(n/16)+(n/8)%2, fmt.Sprintf("%02d", 1+n%16), color(n == model.Head%16, ch), termbox.ColorDefault, func(x, y int) {
			model.SetTrig(n)
		})
	}
	ui.Flush()
}
