package main

import (
	"fmt"
	"io"
	"log"

	"github.com/gdamore/tcell"
	"github.com/gomidi/midi"
	"github.com/gomidi/midi/midimessage/channel"
	"github.com/gomidi/midi/midimessage/realtime"
	"github.com/gomidi/midi/midireader"
	"gitlab.com/gomidi/midi/mid"
)

func main() {
	fmt.Println("TRK")

	/*screen, err := tcell.NewScreen()
	must(err)
	defer screen.Fini()
	err = screen.Init()
	must(err)
	screen.SetStyle(tcell.StyleDefault.
		Foreground(tcell.ColorYellow).
		Background(tcell.ColorBlack))
	screen.Clear()

	s := tcell.StyleDefault.Foreground(tcell.ColorBlack).Background(tcell.ColorYellow)
	s = tcell.StyleDefault

	bar := 0
	for step := 0; step < 16; step++ {
		if step == 0 {
			SetString(screen, 1, 1+step, fmt.Sprintf("> %03d %02d", 1+bar, 1+step), s)
		} else {
			SetString(screen, 1, 1+step, fmt.Sprintf("  %03d %02d", 1+bar, 1+step), s)
		}
	}
	screen.Show()

	quit := make(chan struct{})
	go func() {
		for {
			ev := screen.PollEvent()
			switch ev := ev.(type) {
			case *tcell.EventKey:
				switch ev.Key() {
				case tcell.KeyEscape:
					close(quit)
					return
				case tcell.KeyCtrlL:
					screen.Sync()
				}
			case *tcell.EventResize:
				screen.Sync()
			}
		}
	}()

	<-quit*/

	// Scene
	scene := Scene{
		Name: "Scene",
		Children: []interface{}{
			&Group{
				Name: "Keys",
				Pads: []*Pad{
					&Pad{"C", 1, 1, 50, 48, true},
					&Pad{"D", 1, 2, 50, 48, false},
					&Pad{"E", 1, 3, 50, 48, false},
					&Pad{"F", 1, 4, 50, 48, false},
					&Pad{"G", 1, 5, 50, 48, false},
					&Pad{"A", 1, 6, 50, 48, false},
					&Pad{"B", 1, 7, 50, 48, false},
					&Pad{"c", 1, 8, 50, 48, false},
					&Pad{"C#", 2, 2, 50, 48, false},
					&Pad{"D#", 2, 3, 50, 48, false},
					&Pad{"F#", 2, 5, 50, 48, false},
					&Pad{"G#", 2, 6, 50, 48, false},
					&Pad{"A#", 2, 7, 50, 48, false},
				},
			},
			&Group{
				Name: "Chords",
				Pads: []*Pad{
					&Pad{"Major", 8, 1, 42, 40, true},
					&Pad{"Minor", 8, 2, 42, 40, false},
					&Pad{"Diminished", 8, 3, 42, 40, false},
					&Pad{"Major Seventh", 7, 1, 42, 40, false},
					&Pad{"Minor Seventh", 7, 2, 42, 40, false},
					&Pad{"Dominant Seventh", 7, 3, 42, 40, false},
					&Pad{"Sus2", 6, 1, 42, 40, false},
					&Pad{"Sus4", 6, 2, 42, 40, false},
				},
			},
		},
	}

	//

	var d midiDriver

	ins, err := d.Ins()
	must(err)

	outs, err := d.Outs()
	must(err)

	for _, port := range ins {
		fmt.Printf("[%v] %s\n", port.Number(), port.String())
	}
	for _, port := range outs {
		fmt.Printf("[%v] %s\n", port.Number(), port.String())
	}

	out, err := mid.OpenOut(midiDriver{}, -1, "MIDIOUT2 (LPMiniMK3 MIDI)")
	must(err)
	defer out.Close()

	wr := mid.ConnectOut(out)
	clear(wr)
	scene.Draw(wr)

	//in, err := mid.OpenIn(midiDriver{}, -1, "Elektron Digitone")
	in, err := mid.OpenIn(midiDriver{}, -1, "MIDIIN2 (LPMiniMK3 MIDI)")
	must(err)
	defer in.Close()

	r, w := io.Pipe()

	dataC := make(chan []byte)
	go func() {
		for data := range dataC {
			w.Write(data)
		}
	}()

	err = in.SetListener(func(data []byte, deltaMicroseconds int64) {
		dataC <- data
	})

	rd := midireader.New(r, func(m realtime.Message) {
		if m != realtime.TimingClock {
			fmt.Println(m)
		}
	})
	for {
		m, err := rd.Read()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(m)
		scene.OnMessage(m)
		scene.Draw(wr)

		switch m.(type) {
		case channel.NoteOn:

		case channel.NoteOff:
		}
	}

	/*var (
		w = mid.ConnectOut(out)
		r = mid.NewReader()
	)

	r.Msg.Each = func(pos *mid.Position, msg midi.Message) {
		w.Write(msg)
	}

	mid.ConnectIn(in, r)

	time.Sleep(1 * time.Hour)*/
}

func must(err error) {
	if err != nil {
		panic(err.Error())
	}
}

func SetString(scr tcell.Screen, x, y int, str string, style tcell.Style) {
	for i, r := range str {
		scr.SetContent(x+i, y, r, nil, style)
	}
}

func clear(w *mid.Writer) {
	for i := uint8(1); i < 8; i++ {
		for j := uint8(1); j < 8; j++ {
			w.Write(channel.Channel0.NoteOn(i*10+j, 0))
		}
	}
}

type Pad struct {
	Name        string
	Row, Column uint8
	On, Off     uint8
	IsOn        bool
}

func (p Pad) Draw(w *mid.Writer) {
	c := p.Off
	if p.IsOn {
		c = p.On
	}
	w.Write(channel.Channel0.NoteOn(p.Key(), c))
}

func (p Pad) Key() uint8 {
	return p.Row*10 + p.Column
}

func (p Pad) String() string {
	return p.Name
}

func (p *Pad) Toggle(key uint8) bool {
	if key != p.Key() {
		return false
	}
	p.IsOn = !p.IsOn
	return true
}

type Drawable interface {
	Draw(*mid.Writer)
}

type Group struct {
	Name string
	Pads []*Pad
}

func (g Group) Draw(w *mid.Writer) {
	for _, p := range g.Pads {
		p.Draw(w)
	}
}

func (g Group) OnMessage(m midi.Message) {
	switch m := m.(type) {
	case channel.NoteOn:
		var target *Pad
		for _, p := range g.Pads {
			if p.Key() == m.Key() {
				target = p
				break
			}
		}
		if target == nil {
			return
		}
		for _, p := range g.Pads {
			p.IsOn = p == target
		}
		switch target.Name {
		case "C":
			chord.Root = 60
		case "C#":
			chord.Root = 61
		case "D":
			chord.Root = 62
		case "D#":
			chord.Root = 63
		case "E":
			chord.Root = 64
		case "F":
			chord.Root = 65
		case "F#":
			chord.Root = 66
		case "G":
			chord.Root = 67
		case "G#":
			chord.Root = 68
		case "A":
			chord.Root = 69
		case "A#":
			chord.Root = 70
		case "B":
			chord.Root = 71
		case "c":
			chord.Root = 72

		case "Major":
			chord.Intervals = []uint8{4, 7}
		case "Minor":
			chord.Intervals = []uint8{3, 7}
		case "Diminished":
			chord.Intervals = []uint8{3, 6}

		case "Major Seventh":
			chord.Intervals = []uint8{4, 7, 11}
		case "Minor Seventh":
			chord.Intervals = []uint8{3, 7, 10}
		case "Dominant Seventh":
			chord.Intervals = []uint8{4, 7, 10}

		case "Sus2":
			chord.Intervals = []uint8{2, 7}
		case "Sus4":
			chord.Intervals = []uint8{5, 7}
		}
	}
}

func (g Group) SelectedPad() *Pad {
	for _, p := range g.Pads {
		if p.IsOn {
			return p
		}
	}
	return nil
}

type Scene struct {
	Name     string
	Children []interface{}
}

func (s Scene) Draw(w *mid.Writer) {
	for _, c := range s.Children {
		if c, ok := c.(Drawable); ok {
			c.Draw(w)
		}
	}
}

func (s Scene) OnMessage(m midi.Message) {
	for _, c := range s.Children {
		if c, ok := c.(MessageHandler); ok {
			c.OnMessage(m)
		}
	}
}

func (s Scene) Group(name string) *Group {
	for _, c := range s.Children {
		if g, ok := c.(*Group); ok && g.Name == name {
			return g
		}
	}
	return nil
}

type MessageHandler interface {
	OnMessage(midi.Message)
}

type Chord struct {
	Root      uint8
	Intervals []uint8
}

var chord = new(Chord)
