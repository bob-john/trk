package main

import (
	"fmt"
	"os"
	"strconv"

	"trk/track"

	"github.com/gomidi/midi/midimessage/channel"
	"github.com/gomidi/midi/midimessage/realtime"
	"github.com/nsf/termbox-go"
)

var (
	driver, _ = NewDriver()
	ui        = NewUI()
	model     = NewModel()
	player    = NewPlayer()
	recorder  = NewRecorder()
)

func main() {
	var err error

	if len(os.Args) != 2 {
		fmt.Println("usage: trk <path>")
		fmt.Println("trk: invalid command")
		os.Exit(1)
	}

	err = termbox.Init()
	must(err)
	defer termbox.Close()

	model.Track, err = track.Open(os.Args[1])
	must(err)
	defer model.Track.Close()

	player.Play(model.Track, 0)
	recorder.Listen(track.InputPorts(model.Track))

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
					//FIXME
				case termbox.KeyEsc:
					done = true

				case termbox.KeyPgup:
					model.SetPage(model.Page() - 1)
				case termbox.KeyHome:
					model.SetPage(0)
				case termbox.KeyEnd:
					model.SetPage(model.LastPage())

				case termbox.KeyArrowRight:
					model.SetX(model.X() + 1)
				case termbox.KeyArrowLeft:
					model.SetX(model.X() - 1)
				case termbox.KeyArrowDown:
					model.SetY(model.Y() + 1)
				case termbox.KeyArrowUp:
					model.SetY(model.Y() - 1)
				case termbox.KeyPgdn:
					model.SetPage(model.Page() + 1)

				case termbox.KeyDelete:
					track.Clear(model.Track, model.Head)
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
				parts, err := track.Parts(model.Track)
				must(err)
				switch mm := m.Message.(type) {
				case channel.ProgramChange:
					for _, part := range parts {
						if Contains(part.ProgChgPortIn, m.Port) && int(mm.Channel()) == part.ProgChgInCh {
							err = track.SetPattern(model.Track, part, model.Head, int(mm.Program()))
							must(err)
						}
					}

				case channel.ControlChange:
					if mm.Controller() != 94 {
						break
					}
					for _, part := range parts {
						n := part.TrackOf(int(mm.Channel()))
						if Contains(part.MutePortIn, m.Port) && n != -1 {
							err = track.SetMuted(model.Track, part, model.Head, n, mm.Value() != 0)
							must(err)
						}
					}
				}
			}
		}
		recorder.Listen(track.InputPorts(model.Track))
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
		} else if model.Head != oldHead {
			player.Play(model.Track, model.Head)
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

func render() {
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
	var y int
	var fg, bg termbox.Attribute
	var recording = model.State == Recording
	parts, err := track.Parts(model.Track)
	must(err)
	for _, part := range parts {
		var (
			pattern = track.Pattern(model.Track, part, model.Head)
			mute    = track.Mute(model.Track, part, model.Head)
		)
		fg, bg = colors(false, recording, track.IsPartModified(model.Track, part, model.Head))
		DrawString(4, y, part.ShortName, fg, bg)
		fg, bg = colors(false, recording, track.IsPatternModified(model.Track, part, model.Head))
		DrawString(4+len(part.ShortName)+1, y, track.FormatPattern(pattern), fg, bg)
		fg, bg = colors(false, recording, track.IsMuteModified(model.Track, part, model.Head))
		DrawString(8+len(part.ShortName)+1, y, track.FormatMute(mute, part), fg, bg)
		y++
	}
	y++
	DrawString(0, y, fmt.Sprintf("%03d", 1+model.Page()), termbox.ColorDefault, termbox.ColorDefault)
	for n := 0; n < 16; n++ {
		var (
			tick   = model.HeadForTrig(n)
			fg, bg = colors(n == model.Head%16, recording, track.IsModified(model.Track, tick))
		)
		DrawString(4+(n%8)*3, y+3*(n/16)+(n/8)%2, fmt.Sprintf("%02d", 1+n%16), fg, bg)
	}
	ui.Render()
	termbox.Flush()
}

func colors(highlighted, recording, modified bool) (termbox.Attribute, termbox.Attribute) {
	if highlighted {
		if recording {
			return termbox.ColorDefault, termbox.ColorRed
		}
		if modified {
			return termbox.ColorRed, termbox.ColorWhite
		}
		return termbox.ColorBlack, termbox.ColorWhite
	}
	if modified {
		return termbox.ColorRed, termbox.ColorDefault
	}
	return termbox.ColorDefault, termbox.ColorDefault
}

func options() *OptionPage {
	var (
		inputs, _  = driver.Ins()
		outputs, _ = driver.Outs()
	)
	addPartOptions := func(page *OptionPage, part *track.Part) {
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
					must(track.SetPart(model.Track, part))
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
					must(track.SetPart(model.Track, part))
				})
			}
		}
		channels := map[int]string{-1: "OFF"}
		for ch := 0; ch < 16; ch++ {
			channels[ch] = strconv.Itoa(ch + 1)
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
			for n, ch := range part.TrackCh {
				n := n
				page.AddPicker(track.FormatTrackName(part.Name, n)+" CH", channels, ch, func(ch int) {
					part.TrackCh[n] = ch
					must(track.SetPart(model.Track, part))
				})
			}
			page.AddPicker("PROG CHG IN CH", channels, part.ProgChgInCh, func(selected int) {
				part.ProgChgInCh = selected
				must(track.SetPart(model.Track, part))
			})
			page.AddPicker("PROG CHG OUT CH", channels, part.ProgChgOutCh, func(selected int) {
				part.ProgChgOutCh = selected
				must(track.SetPart(model.Track, part))
			})
		})
	}
	parts, err := track.Parts(model.Track)
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
