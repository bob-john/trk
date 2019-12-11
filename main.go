package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"reflect"
	"text/tabwriter"
	"trk/rtmididrv"
	"trk/track"

	"github.com/nsf/termbox-go"
	"gitlab.com/gomidi/midi/midimessage/realtime"
)

var (
	midiDriver, _ = rtmididrv.New()
	ui            = NewUI()
	model         = NewModel()
	player        = NewPlayer()
	recorder      = NewRecorder()
	console       = NewConsole()
)

func main() {
	var err error

	defer player.Close()
	defer recorder.Close()

	flag.Parse()

	if len(flag.Args()) != 1 {
		fmt.Println("usage: trk <path>")
		fmt.Println("trk: invalid command")
		os.Exit(1)
	}

	model.Track, err = track.Open(flag.Arg(0))
	must(err)
	defer model.Track.Close()

	err = termbox.Init()
	must(err)
	defer termbox.Close()

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

	for !model.Done {
		render()

		recorder.Listen(model.Track.InPorts())

		select {
		case e := <-eventC:
			if ui.Handle(e) {
				break
			}
			if e.Type == termbox.EventKey {
				switch e.Key {
				case termbox.KeyEsc:
					model.Done = true

				case termbox.KeyCtrlO:
					ui.Show(NewDialog(options()))
				case termbox.KeyCtrlD:
					console.Enabled = !console.Enabled

				case termbox.KeyEnter, termbox.KeySpace:
					if model.Recording {
						must(model.Track.Save())
					}
					model.Recording = !model.Recording
				}
			}

		case m := <-midiC:
			log.Print(reflect.TypeOf(m.Message), m.Message == realtime.TimingClock)
			switch m.Message {
			case realtime.TimingClock:
				if model.Playing {
					model.Tick++
					// player.Play(model.Track.OutPorts(m.Port, m.Message.Raw()), m.Message)
				}
			case realtime.Start:
				model.Playing = true
				model.Tick = 0
			case realtime.Continue:
				model.Playing = true
			case realtime.Stop:
				model.Playing = false
			}
			if model.Recording {
				port, message := m.Port, m.Message.Raw()
				if len(message) > 0 && message[0] >= 0x80 && message[0] <= 0xf0 {
					model.Track.Insert(model.Tick, port, message)
				}
			}
		}
	}

	termbox.Interrupt()

	must(model.Track.Save())
}

func render() {
	_, pageSize := termbox.Size()
	pageSize -= 2
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
	w := tabwriter.NewWriter(&Writer{}, 0, 0, 3, ' ', 0)
	fmt.Fprintln(w, "#\tBar.16th\tPort\tCh\tType\tSubtype/Note\tValue")
	events, base := model.Track.Events, 1
	if len(events) > pageSize {
		events, base = events[len(events)-pageSize:], len(events)-pageSize
	}
	for i, e := range events {
		fmt.Fprintf(w, "%o\t%s\t%s\t%d\t%s\t%s\t%d\n", base+i, e.Beat(), e.Port, e.Channel(), e.Type(), e.Subtype(), e.Value())
	}
	w.Flush()
	console.Render()
	ui.Render()
	termbox.Flush()
}

func options() (p *OptionPage) {
	inputs, err := NewInputPortList()
	must(err)
	outputs, err := NewOutputPortList()
	must(err)

	p = NewOptionPage("DEVICES CONFIG")
	for _, src := range model.Track.Devices {
		src := src
		p.Page(src.Name, func(p *OptionPage) {
			p.Picker("INPUT PORT", inputs, inputs.IndexOf(src.Input), func(val int) {
				src.Input = inputs[val]
				must(model.Track.Save())
			})
			p.Picker("OUTPUT PORT", outputs, outputs.IndexOf(src.Output), func(val int) {
				src.Output = outputs[val]
				must(model.Track.Save())
			})
			for _, dst := range model.Track.Devices {
				if dst.Name == src.Name {
					continue
				}
				r, err := model.Track.CreateRouteIfNotExist(src.Name, dst.Name)
				must(err)
				p.Page(r.String(), func(p *OptionPage) {
					p.Checkbox("PROG CH", r.ProgCh, func(val bool) {
						r.ProgCh = val
						must(model.Track.Save())
					})
					p.Checkbox("NOTES", r.Notes, func(val bool) {
						r.Notes = val
						must(model.Track.Save())
					})
					p.Checkbox("CC/NRPN", r.CC, func(val bool) {
						r.CC = val
						must(model.Track.Save())
					})
				})
			}
		})
	}
	return
}

type PortList map[int]string

func NewInputPortList() (PortList, error) {
	ports, err := midiDriver.Ins()
	if err != nil {
		return nil, err
	}
	values := map[int]string{-1: "OFF"}
	for i, port := range ports {
		values[i] = port.String()
	}
	return PortList(values), nil
}

func NewOutputPortList() (PortList, error) {
	ports, err := midiDriver.Outs()
	if err != nil {
		return nil, err
	}
	values := map[int]string{-1: "OFF"}
	for i, port := range ports {
		values[i] = port.String()
	}
	return PortList(values), nil
}

func (pl PortList) IndexOf(port string) int {
	for i, name := range pl {
		if i != -1 && name == port {
			return i
		}
	}
	return -1
}

// var (
// 	midiDriver, _ = rtmididrv.New()
// 	ui            = NewUI()
// 	model         = NewModel()
// 	player        = NewPlayer()
// 	recorder      = NewRecorder()
// 	console       = NewConsole()
// )

// func main() {
// 	var err error

// 	defer player.Close()
// 	defer recorder.Close()

// 	flag.Parse()

// 	if len(flag.Args()) != 1 {
// 		fmt.Println("usage: trk <path>")
// 		fmt.Println("trk: invalid command")
// 		os.Exit(1)
// 	}

// 	model.Track, err = track.Open(flag.Arg(0))
// 	must(err)
// 	defer model.Track.Close()

// 	err = termbox.Init()
// 	must(err)
// 	defer termbox.Close()

// 	player.Play(model.Track, 0)
// 	recorder.Listen(model.Track.Inputs())

// 	var (
// 		eventC = make(chan termbox.Event)
// 		midiC  = recorder.C()
// 	)

// 	go func() {
// 		for {
// 			e := termbox.PollEvent()
// 			if e.Type == termbox.EventInterrupt {
// 				break
// 			}
// 			eventC <- e
// 		}
// 	}()

// 	var (
// 		done bool
// 		tick int
// 	)
// 	for !done {
// 		render()

// 		var (
// 			oldHead = model.Head
// 		)
// 		select {
// 		case e := <-eventC:
// 			if ui.Handle(e) {
// 				break
// 			}
// 			switch e.Type {
// 			case termbox.EventKey:
// 				switch e.Key {
// 				case termbox.KeyCtrlO:
// 					ui.Show(NewDialog(5, 5, options()))

// 				case termbox.KeyCtrlS:
// 					//FIXME
// 				case termbox.KeyEsc:
// 					done = true

// 				case termbox.KeyCtrlD:
// 					console.Enabled = !console.Enabled

// 				case termbox.KeyPgup:
// 					model.SetPage(model.Page() - 1)
// 				case termbox.KeyHome:
// 					model.SetHead(0)
// 				case termbox.KeyEnd:
// 					model.SetHead(model.LastStep())

// 				case termbox.KeyArrowRight:
// 					model.SetX(model.X() + 1)
// 				case termbox.KeyArrowLeft:
// 					model.SetX(model.X() - 1)
// 				case termbox.KeyArrowDown:
// 					model.SetY(model.Y() + 1)
// 				case termbox.KeyArrowUp:
// 					model.SetY(model.Y() - 1)
// 				case termbox.KeyPgdn:
// 					model.SetPage(model.Page() + 1)

// 				case termbox.KeyDelete:
// 					model.Track.Clear(model.Head)
// 				case termbox.KeyEnter:
// 					model.ToggleRecording()
// 				}
// 			}

// 		case m := <-midiC:
// 			switch m.Message {
// 			case realtime.TimingClock:
// 				if model.State == Playing {
// 					tick++
// 				}
// 			case realtime.Start:
// 				model.State = Playing
// 				tick = 0
// 			case realtime.Continue:
// 				model.State = Playing
// 			case realtime.Stop:
// 				model.State = Viewing
// 			}
// 			if model.State == Recording {
// 				switch mm := m.Message.(type) {
// 				case channel.ProgramChange:
// 					for _, part := range model.Track.Parts() {
// 						if Contains(part.ProgChgPortIn, m.Port) && int(mm.Channel()) == part.ProgChgInCh {
// 							err = model.Track.SetPattern(part, model.Head, int(mm.Program()))
// 							must(err)
// 						}
// 					}

// 				case channel.ControlChange:
// 					if mm.Controller() != 94 {
// 						break
// 					}
// 					for _, part := range model.Track.Parts() {
// 						n := part.TrackOf(int(mm.Channel()))
// 						if Contains(part.MutePortIn, m.Port) && n != -1 {
// 							err = model.Track.SetMuted(part, model.Head, n, mm.Value() != 0)
// 							must(err)
// 						}
// 					}
// 				}
// 			}
// 		}
// 		recorder.Listen(model.Track.Inputs())
// 		if model.State == Playing {
// 			switch tick {
// 			case 12:
// 				player.PlayPattern(model.Track, model.Head+2)

// 			case 18:
// 				player.PlayMute(model.Track, model.Head+1)

// 			case 24:
// 				model.Head++
// 				tick = 0
// 			}
// 		} else if model.Head != oldHead {
// 			player.Play(model.Track, model.Head)
// 		}
// 	}
// 	termbox.Interrupt()
// }

// func render() {
// 	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
// 	var y int
// 	var fg, bg termbox.Attribute
// 	var recording = model.State == Recording
// 	for _, part := range model.Track.Parts() {
// 		var (
// 			pattern = model.Track.Pattern(part, model.Head)
// 			mute    = model.Track.Mute(part, model.Head)
// 		)
// 		fg, bg = colors(false, recording, model.Track.IsPartModified(part, model.Head))
// 		DrawString(4, y, part.ShortName, fg, bg)
// 		fg, bg = colors(false, recording, model.Track.IsPatternModified(part, model.Head))
// 		DrawString(4+len(part.ShortName)+1, y, track.FormatPattern(pattern), fg, bg)
// 		fg, bg = colors(false, recording, model.Track.IsMuteModified(part, model.Head))
// 		DrawString(8+len(part.ShortName)+1, y, track.FormatMute(mute, part), fg, bg)
// 		y++
// 	}
// 	y++
// 	DrawString(0, y, fmt.Sprintf("%03d", 1+model.Page()), termbox.ColorDefault, termbox.ColorDefault)
// 	for n := 0; n < 16; n++ {
// 		var (
// 			tick   = model.HeadForTrig(n)
// 			fg, bg = colors(n == model.Head%16, recording, model.Track.IsModified(tick))
// 		)
// 		DrawString(4+(n%8)*3, y+3*(n/16)+(n/8)%2, fmt.Sprintf("%02d", 1+n%16), fg, bg)
// 	}
// 	console.Render()
// 	ui.Render()
// 	termbox.Flush()
// }

// func colors(highlighted, recording, modified bool) (termbox.Attribute, termbox.Attribute) {
// 	if highlighted {
// 		if recording {
// 			return termbox.ColorDefault, termbox.ColorRed
// 		}
// 		if modified {
// 			return termbox.ColorRed, termbox.ColorWhite
// 		}
// 		return termbox.ColorBlack, termbox.ColorWhite
// 	}
// 	if modified {
// 		return termbox.ColorRed, termbox.ColorDefault
// 	}
// 	return termbox.ColorDefault, termbox.ColorDefault
// }

// func options() *OptionPage {
// 	var (
// 		inputs, _  = midiDriver.Ins()
// 		outputs, _ = midiDriver.Outs()
// 	)
// 	addPartOptions := func(page *OptionPage, part *track.Part) {
// 		addInputs := func(page *OptionPage, ports *[]string) {
// 			for _, port := range inputs {
// 				name := port.String()
// 				on := Contains(*ports, name)
// 				page.Checkbox(" "+name, on, func(on bool) {
// 					if on {
// 						*ports = Insert(*ports, name)
// 					} else {
// 						*ports = Remove(*ports, name)
// 					}
// 					must(model.Track.SetPart(part))
// 				})
// 			}
// 		}
// 		addOutputs := func(page *OptionPage, ports *[]string) {
// 			for _, port := range outputs {
// 				name := port.String()
// 				on := Contains(*ports, name)
// 				page.Checkbox(" "+name, on, func(on bool) {
// 					if on {
// 						*ports = Insert(*ports, name)
// 					} else {
// 						*ports = Remove(*ports, name)
// 					}
// 					must(model.Track.SetPart(part))
// 				})
// 			}
// 		}
// 		channels := map[int]string{-1: "OFF"}
// 		for ch := 0; ch < 16; ch++ {
// 			channels[ch] = strconv.Itoa(ch + 1)
// 		}
// 		page.Page("PORT CONFIG", func(page *OptionPage) {
// 			page.Label("PROG CHG PORT IN")
// 			addInputs(page, &part.ProgChgPortIn)
// 			page.Label("PROG CHG PORT OUT")
// 			addOutputs(page, &part.ProgChgPortOut)
// 			page.Label("MUTE PORT IN")
// 			addInputs(page, &part.MutePortIn)
// 			page.Label("MUTE PORT OUT")
// 			addOutputs(page, &part.MutePortOut)
// 		})
// 		page.Page("CHANNELS", func(page *OptionPage) {
// 			for n, ch := range part.TrackCh {
// 				n := n
// 				page.Picker(track.FormatTrackName(part.Name, n)+" CH", channels, ch, func(ch int) {
// 					part.TrackCh[n] = ch
// 					must(model.Track.SetPart(part))
// 				})
// 			}
// 			page.Picker("PROG CHG IN CH", channels, part.ProgChgInCh, func(selected int) {
// 				part.ProgChgInCh = selected
// 				must(model.Track.SetPart(part))
// 			})
// 			page.Picker("PROG CHG OUT CH", channels, part.ProgChgOutCh, func(selected int) {
// 				part.ProgChgOutCh = selected
// 				must(model.Track.SetPart(part))
// 			})
// 		})
// 	}
// 	fillFilterInPage := func(f *track.Filter, page *OptionPage) {
// 		for _, port := range inputs {
// 			name := port.String()
// 			on := Contains(f.Inputs, name)
// 			page.Checkbox(" "+name, on, func(on bool) {
// 				if on {
// 					f.Inputs = Insert(f.Inputs, name)
// 				} else {
// 					f.Inputs = Remove(f.Inputs, name)
// 				}
// 				must(f.Save(model.Track))
// 			})
// 		}
// 	}
// 	fillFilterOutPage := func(f *track.Filter, page *OptionPage) {
// 		for _, port := range outputs {
// 			name := port.String()
// 			on := Contains(f.Outputs, name)
// 			page.Checkbox(" "+name, on, func(on bool) {
// 				if on {
// 					f.Outputs = Insert(f.Outputs, name)
// 				} else {
// 					f.Outputs = Remove(f.Outputs, name)
// 				}
// 				must(f.Save(model.Track))
// 			})
// 		}
// 	}
// 	fillFilterPage := func(f *track.Filter, page *OptionPage) {
// 		page.Page("In", func(page *OptionPage) {
// 			fillFilterInPage(f, page)
// 		})
// 		page.Page("Out", func(page *OptionPage) {
// 			fillFilterOutPage(f, page)
// 		})
// 		page.Checkbox("Note", f.Note, func(val bool) {
// 			f.Note = val
// 			must(f.Save(model.Track))
// 		})
// 		page.Checkbox("PolyphonicAftertouch", f.PolyphonicAftertouch, func(val bool) {
// 			f.PolyphonicAftertouch = val
// 			must(f.Save(model.Track))
// 		})
// 		page.Checkbox("ControlChange", f.ControlChange, func(val bool) {
// 			f.ControlChange = val
// 			must(f.Save(model.Track))
// 		})
// 		page.Checkbox("ProgramChange", f.ProgramChange, func(val bool) {
// 			f.ProgramChange = val
// 			must(f.Save(model.Track))
// 		})
// 		page.Checkbox("ChannelAftertouch", f.ChannelAftertouch, func(val bool) {
// 			f.ChannelAftertouch = val
// 			must(f.Save(model.Track))
// 		})
// 		page.Checkbox("PitchBendChange", f.PitchBendChange, func(val bool) {
// 			f.PitchBendChange = val
// 			must(f.Save(model.Track))
// 		})
// 	}
// 	options := NewOptionPage("OPTIONS")
// 	options.Page("MIDI ROUTING", func(page *OptionPage) {
// 		for _, f := range model.Track.Filters() {
// 			page.Page(fmt.Sprintf("%s -> %s", f.Inputs, f.Outputs), func(page *OptionPage) {
// 				fillFilterPage(f, page)
// 			})
// 		}
// 		page.Page("+ NEW ROUTE", func(page *OptionPage) {
// 			fillFilterPage(&track.Filter{}, page)
// 		})
// 	})
// 	options.Page("MIDI CONFIG", func(page *OptionPage) {
// 		for _, part := range model.Track.Parts() {
// 			page.Page(part.Name, func(page *OptionPage) {
// 				addPartOptions(page, part)
// 			})
// 		}
// 	})
// 	return options
// }
