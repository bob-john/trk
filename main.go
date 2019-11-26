package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/asdine/storm"
	"github.com/gomidi/midi/midimessage/realtime"
	"github.com/nsf/termbox-go"
)

var (
	ui       = NewUI()
	model    = NewModel()
	player   = NewPlayer()
	recorder = NewRecorder()
	trk      *storm.DB
	digitakt *Input
	digitone *Input
)

func main() {
	var err error

	if len(os.Args) != 2 {
		fmt.Println("usage: trk <path>")
		fmt.Println("trk: invalid command")
		os.Exit(1)
	}

	trk, err = OpenTrack(os.Args[1])
	must(err)
	defer trk.Close()

	err = model.LoadTrack(os.Args[1])
	must(err)

	err = termbox.Init()
	must(err)
	defer termbox.Close()

	player.Play(model.Track, 0)
	recorder.Listen(model.Track.Settings.InputPortNames())

	var (
		eventC = make(chan termbox.Event)
		midiC  = recorder.C()
	)

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
				model.Track.Seq.Insert(m.Port, model.Head, m.Message, model.Track.Settings)
			}
		}
		recorder.Listen(model.Track.Settings.InputPortNames())
		if model.State == Playing {
			switch tick {
			case 12:
				player.PlayPattern(model.Track, model.Head+2)

			case 18:
				player.PlayMute(model.Track, model.Head+1)

			case 24:
				model.Head++
				tick = 0
			}
		} else {
			if model.Head != oldHead {
				player.Play(model.Track, model.Head)
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
	var y int
	for name, part := range row.Parts {
		DrawString(4, y, name+":", termbox.ColorDefault, termbox.ColorDefault)
		DrawString(4+len(name)+2, y, part.Pattern.String(), color(false, org.Parts[name].Pattern != -1), termbox.ColorDefault)
		DrawString(8+len(name)+2, y, part.Mute.Format(part.Channels), color(false, len(org.Parts[name].Mute) != 0), termbox.ColorDefault)
		y++
	}
	y++
	DrawString(0, y, fmt.Sprintf("%03d", 1+model.Pattern()), termbox.ColorDefault, termbox.ColorDefault)
	for n := 0; n < 16; n++ {
		n := n
		ch := model.HeadForTrig(n) == 0 || model.Track.Seq.Row(model.HeadForTrig(n)).HasChanges(model.Track.Seq.Row(model.HeadForTrig(n-1)))
		DrawString(4+(n%8)*3, y+3*(n/16)+(n/8)%2, fmt.Sprintf("%02d", 1+n%16), color(n == model.Head%16, ch), termbox.ColorDefault)
	}
	ui.Render()
	termbox.Flush()
}

func options() *OptionPage {
	var (
		inputs, _  = driver.Ins()
		outputs, _ = driver.Outs()
	)
	addPartOptions := func(page *OptionPage, part *Part1) {
		addInputs := func(page *OptionPage, ports *[]string) {
			for _, port := range inputs {
				name := port.String()
				on := Contains(*ports, name)
				page.AddCheckbox(" "+name, on, func(on bool) {
					if on {
						*ports = Insert(*ports, name)
					} else {
						*ports = Remove(*ports, name)
					}
					must(trk.Save(part))
				})
			}
		}
		addOutputs := func(page *OptionPage, ports *[]string) {
			for _, port := range outputs {
				name := port.String()
				on := Contains(*ports, name)
				page.AddCheckbox(" "+name, on, func(on bool) {
					if on {
						*ports = Insert(*ports, name)
					} else {
						*ports = Remove(*ports, name)
					}
					must(trk.Save(part))
				})
			}
		}
		var channels = []string{"OFF"}
		for n := 0; n < 16; n++ {
			channels = append(channels, strconv.Itoa(1+n))
		}
		page.AddMenu("PORT CONFIG", func(page *OptionPage) {
			page.AddLabel("PROG CHG PORT IN")
			addInputs(page, &part.ProgChgPortIn)
			page.AddLabel("PROG CHG PORT OUT")
			addOutputs(page, &part.ProgChgPortOut)
			page.AddLabel("MUTE PORT IN")
			addInputs(page, &part.MutePortIn)
			page.AddLabel("MUTE PORT OUT")
			addOutputs(page, &part.MutePortOut)
		})
		page.AddMenu("CHANNELS", func(page *OptionPage) {
			for n, ch := range part.Track {
				n := n
				page.AddPicker(FormatTrackName(part.Name, n)+" CH", channels, ch, func(ch int) {
					part.Track[n] = ch
					must(trk.Save(part))
				})
			}
			page.AddPicker("PROG CHG IN CH", channels, part.ProgChgInCh, func(selected int) {
				part.ProgChgInCh = selected
				must(trk.Save(part))
			})
			page.AddPicker("PROG CHG OUT CH", channels, part.ProgChgOutCh, func(selected int) {
				part.ProgChgOutCh = selected
				must(trk.Save(part))
			})
		})
	}
	var parts []*Part1
	err := trk.All(&parts)
	must(err)
	options := NewOptionPage("MIDI CONFIG")
	for _, part := range parts {
		options.AddMenu(part.Name, func(page *OptionPage) {
			addPartOptions(page, part)
		})
	}
	return options
}

func Contains(a []string, x string) bool {
	for _, n := range a {
		if x == n {
			return true
		}
	}
	return false
}

func IndexOf(a []string, x string) int {
	for i, n := range a {
		if n == x {
			return i
		}
	}
	return -1
}

func Insert(a []string, x string) []string {
	if Contains(a, x) {
		return a
	}
	return append(a, x)
}

func Remove(a []string, x string) []string {
	i := IndexOf(a, x)
	if i == -1 {
		return a
	}
	return append(a[:i], a[i+1:]...)
}
