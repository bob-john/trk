package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/gomidi/midi/midimessage/realtime"
	"github.com/nsf/termbox-go"
)

var (
	ui       = NewUI()
	model    = new(Model)
	digitakt *Device
	digitone *Device
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: trk <path>")
		fmt.Println("trk: invalid command")
		os.Exit(1)
	}

	err := model.LoadTrack(os.Args[1])
	must(err)

	err = termbox.Init()
	must(err)
	defer termbox.Close()

	digitakt, _ = OpenDevice("Digitakt", "Elektron Digitakt", "Elektron Digitakt")
	digitone, _ = OpenDevice("Digitone", "Elektron Digitone", "Elektron Digitone")

	model.Track.Seq.ConsolidatedRow(0).Play(digitone, digitakt)

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
			if ui.Handle(e) {
				break
			}
			switch e.Type {
			case termbox.EventKey:
				switch e.Key {
				case termbox.KeyCtrlO:
					ui.Show(NewDialog(5, 5, options()))

				case termbox.KeyCtrlS:
					err := model.Track.Write(os.Args[1])
					if err != nil {
						log.Fatal(err)
					}
				case termbox.KeyEsc:
					err := model.Track.Write(os.Args[1])
					if err != nil {
						log.Fatal(err)
					}
					done = true

				case termbox.KeyPgup:
					model.SetPattern(model.Pattern() - 1)
				case termbox.KeyHome:
					model.SetPattern(0)
				case termbox.KeyEnd:
					model.SetPattern(model.LastPattern())

				case termbox.KeyArrowRight:
					model.SetX(model.X() + 1)
				case termbox.KeyArrowLeft:
					model.SetX(model.X() - 1)
				case termbox.KeyArrowDown:
					model.SetY(model.Y() + 1)
				case termbox.KeyArrowUp:
					model.SetY(model.Y() - 1)
				case termbox.KeyPgdn:
					model.SetPattern(model.Pattern() + 1)

				case termbox.KeyDelete:
					model.ClearStep()
				case termbox.KeyEnter:
					model.ToggleRecording()
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
				model.Track.Seq.Insert(m.Device.Name(), model.Head, m.Message)
			}
		}
		if model.State == Playing {
			switch tick {
			case 12:
				row := model.Track.Seq.ConsolidatedRow(model.Head + 2)
				row.Digitone.Pattern.Play(digitone, 15)
				row.Digitakt.Pattern.Play(digitakt, 15)

			case 18:
				row := model.Track.Seq.ConsolidatedRow(model.Head + 1)
				row.Digitone.Mute.Play(digitone, row.Digitone.Channels)
				row.Digitakt.Mute.Play(digitakt, row.Digitakt.Channels)
				row.Digitone.Mute.Play(digitakt, row.Digitone.Channels) //HACK

			case 24:
				model.Head++
				tick = 0
			}
		} else {
			if model.Head != oldHead {
				model.Track.Seq.ConsolidatedRow(model.Head).Play(digitone, digitakt)
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
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
	row, org := model.Track.Seq.ConsolidatedRow(model.Head), model.Track.Seq.Row(model.Head)
	DrawString(4, 0, row.Digitakt.Pattern.String(), color(false, org.Digitakt.Pattern != -1), termbox.ColorDefault)
	DrawString(8, 0, row.Digitakt.Mute.Format(row.Digitakt.Channels), color(false, len(org.Digitakt.Mute) != 0), termbox.ColorDefault)
	DrawString(8+row.Digitakt.Channels.Len+1, 0, row.Digitone.Pattern.String(), color(false, org.Digitone.Pattern != -1), termbox.ColorDefault)
	DrawString(12+row.Digitakt.Channels.Len+1, 0, row.Digitone.Mute.Format(row.Digitone.Channels), color(false, len(org.Digitone.Mute) != 0), termbox.ColorDefault)
	DrawString(0, 2, fmt.Sprintf("%03d", 1+model.Pattern()), termbox.ColorDefault, termbox.ColorDefault)
	for n := 0; n < 16; n++ {
		n := n
		ch := model.HeadForTrig(n) == 0 || model.Track.Seq.Row(model.HeadForTrig(n)).HasChanges(model.Track.Seq.Row(model.HeadForTrig(n-1)))
		DrawString(4+(n%8)*3, 2+3*(n/16)+(n/8)%2, fmt.Sprintf("%02d", 1+n%16), color(n == model.Head%16, ch), termbox.ColorDefault)
	}
	ui.Render()
	termbox.Flush()
}

func options() *OptionPage {
	var (
		inputs, _  = driver.Ins()
		outputs, _ = driver.Outs()
	)
	addDeviceOptions := func(p *OptionPage, s *DeviceSettings) {
		addInputs := func(p *OptionPage, s *DeviceSettings) {
			for _, port := range inputs {
				name := port.String()
				_, on := s.Inputs[name]
				p.AddCheckbox(" "+port.String(), on, func(on bool) {
					if on {
						s.Inputs[name] = struct{}{}
					} else {
						delete(s.Inputs, name)
					}
				})
			}
		}
		addOutputs := func(p *OptionPage, s *DeviceSettings) {
			for _, port := range outputs {
				name := port.String()
				_, on := s.Outputs[name]
				p.AddCheckbox(" "+port.String(), on, func(on bool) {
					if on {
						s.Outputs[name] = struct{}{}
					} else {
						delete(s.Outputs, name)
					}
				})
			}
		}
		var channels = []string{"Off"}
		for n := 0; n < 16; n++ {
			channels = append(channels, strconv.Itoa(1+n))
		}
		p.AddMenu("Port config", func(page *OptionPage) {
			page.AddLabel("Input")
			addInputs(page, s)
			page.AddLabel("Output")
			addOutputs(page, s)
		})
		p.AddMenu("Program Change", func(page *OptionPage) {
			page.AddPicker("Record from", []string{"Digitatk", "Digitone", "Both"}, int(s.ProgChgSrc), func(selected int) {
				s.ProgChgSrc = DeviceSource(selected)
			})
			page.AddPicker("Input channel", channels, s.ProgChgInCh, func(selected int) {
				s.ProgChgInCh = selected
			})
			page.AddPicker("Output channel", channels, s.ProgChgOutCh, func(selected int) {
				s.ProgChgOutCh = selected
			})
		})
		p.AddMenu("Mute", func(page *OptionPage) {
			page.AddPicker("Record from", []string{"Digitatk", "Digitone", "Both"}, int(s.MuteSrc), func(selected int) {
				s.MuteSrc = DeviceSource(selected)
			})
			for n := 0; n < 16; n++ {
				ch := n
				_, on := s.Channels[ch]
				page.AddCheckbox(fmt.Sprintf("Channel %d", 1+n), on, func(on bool) {
					if on {
						s.Channels[ch] = struct{}{}
					} else {
						delete(s.Channels, ch)
					}
				})
			}
		})
	}
	options := NewOptionPage("MIDI config")
	options.AddMenu("Digitakt", func(page *OptionPage) {
		addDeviceOptions(page, model.Track.Settings.Digitakt)
	})
	options.AddMenu("Digitone", func(page *OptionPage) {
		addDeviceOptions(page, model.Track.Settings.Digitone)
	})
	return options
}
