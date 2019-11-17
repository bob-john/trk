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
	seq       *Seq
	digitakt  *Device
	digitone  *Device
	head      int
	playing   bool
	recording bool
)

func main() {
	var err error

	if len(os.Args) != 2 {
		fmt.Println("usage: trk <path> [<device> ...]")
		fmt.Println("trk: invalid command line")
		os.Exit(1)
	}

	seq, err = ReadSeq(os.Args[1])
	if os.IsNotExist(err) {
		seq = NewSeq()
	} else if err != nil {
		must(err)
	}

	err = termbox.Init()
	must(err)
	defer termbox.Close()

	digitakt, _ = OpenDevice("Digitakt", "Elektron Digitakt", "Elektron Digitakt")
	digitone, _ = OpenDevice("Digitone", "Elektron Digitone", "Elektron Digitone")

	seq.ConsolidatedRow(0).Play(digitone, digitakt)

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
		oldHead := head
		row, col := head%16, head/16
		select {
		case e := <-eventC:
			switch e.Type {
			case termbox.EventKey:
				switch e.Key {
				case termbox.KeyCtrlS:
					err := seq.Write(os.Args[1])
					if err != nil {
						log.Fatal(err)
					}
				case termbox.KeyCtrlX:
					err := seq.Write(os.Args[1])
					if err != nil {
						log.Fatal(err)
					}
					done = true

				case termbox.KeyArrowRight:
					if !playing {
						col++
					}
				case termbox.KeyArrowLeft:
					if !playing {
						col--
					}
				case termbox.KeyArrowUp:
					if !playing {
						row--
					}
				case termbox.KeyArrowDown:
					if !playing {
						row++
					}

				case termbox.KeySpace:
					recording = !recording
				}
			}

		case m := <-midiC:
			switch m.Message {
			case realtime.TimingClock:
				if playing {
					tick++
				}
			case realtime.Start:
				playing = true
				tick = 0
				head = 0
			case realtime.Continue:
				playing = true
			case realtime.Stop:
				playing = false
			}
			seq.Insert(m.Device.Name(), head, m.Message)
		}
		if playing {
			if tick == 0 {
				curr, next := seq.ConsolidatedRow(head), seq.ConsolidatedRow(head+2)
				next.Digitone.Pattern.Play(digitone, 15)
				curr.Digitone.Mute.Play(digitone, curr.Digitone.Channels)
				next.Digitakt.Pattern.Play(digitakt, 15)
				curr.Digitakt.Mute.Play(digitakt, curr.Digitakt.Channels)
			}
			if tick == 6 {
				head++
				tick = 0
			}
		} else {
			row = clamp(row, 0, 16-1)
			col = clamp(col, 0, 96-1)
			head = 16*col + row
			if head != oldHead {
				seq.ConsolidatedRow(head).Play(digitone, digitakt)
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

func render() {
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
	p := head / 16
	for n := 0; n < 16; n++ {
		fg := termbox.ColorBlue
		if n == (head/16)%16 {
			if recording {
				fg = termbox.ColorRed
			} else if playing {
				fg = termbox.ColorGreen
			}
			fg = fg | termbox.AttrReverse
		}
		SetString(2+n*3, 0, fmt.Sprintf("%2d", 1+16*(p/16)+n), fg, termbox.ColorDefault)
	}
	SetString(4, 2, "DN", termbox.ColorBlue, termbox.ColorDefault)
	SetString(13, 2, "DT", termbox.ColorBlue, termbox.ColorDefault)
	for n := 0; n < 16; n++ {
		step := 16*p + n
		fg, bg := termbox.ColorBlue, termbox.ColorDefault
		if step == head {
			if recording {
				fg = termbox.ColorRed
			} else if playing {
				fg = termbox.ColorGreen
			}
			fg = fg | termbox.AttrReverse
		}
		SetString(0, 3+n, seq.Text(step), fg, bg)
	}
	termbox.Flush()
}
