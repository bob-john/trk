package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/gomidi/midi/midimessage/realtime"
	"github.com/nsf/termbox-go"
)

const (
	pageSize = 16
)

var (
	seq      *Seq
	digitakt *Device
	digitone *Device
	head     int
	state    State
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
		var (
			oldHead = head
			pattern = head / 64
			page    = (head % 64) / 16
			trigger = head % 16
		)
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
					if state == Viewing {
						trigger++
					}
				case termbox.KeyArrowLeft:
					if state == Viewing {
						trigger--
					}
				case termbox.KeyArrowUp:
					if state == Viewing {
						trigger -= 8
					}
				case termbox.KeyArrowDown:
					if state == Viewing {
						trigger += 8
					}
				case termbox.KeyPgup:
					if state == Viewing {
						pattern--
					}
				case termbox.KeyPgdn:
					if state == Viewing {
						pattern++
					}
				case termbox.KeyTab, termbox.KeySpace:
					if state == Viewing {
						page++
					}

				case termbox.KeyDelete, termbox.KeyBackspace:
					if state == Viewing || state == Playing {
						seq.Clear(head)
					}
				case termbox.KeyEnter:
					if state == Viewing {
						state = Recording
					} else if state == Recording {
						state = Viewing
					}
				}
			}

		case m := <-midiC:
			switch m.Message {
			case realtime.TimingClock:
				if state == Playing {
					tick++
				}
			case realtime.Start:
				state = Playing
				tick = 0
			case realtime.Continue:
				state = Playing
			case realtime.Stop:
				state = Viewing
			}
			if state == Recording {
				seq.Insert(m.Device.Name(), head, m.Message)
			}
		}
		if state == Playing {
			switch tick {
			case 12:
				row := seq.ConsolidatedRow(head + 2)
				row.Digitone.Pattern.Play(digitone, 15)
				row.Digitakt.Pattern.Play(digitakt, 15)

			case 18:
				row := seq.ConsolidatedRow(head + 1)
				row.Digitone.Mute.Play(digitone, row.Digitone.Channels)
				row.Digitakt.Mute.Play(digitakt, row.Digitakt.Channels)
				row.Digitone.Mute.Play(digitakt, row.Digitone.Channels) //HACK

			case 24:
				head++
				tick = 0
			}
		} else {
			SetString(0, 8, strconv.Itoa(trigger), termbox.ColorDefault, termbox.ColorDefault)
			head = clamp(pattern*64+page*16+trigger, 0, 8*16*64-1)
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

func color(on, ch bool) (fg termbox.Attribute) {
	if ch {
		fg = termbox.ColorRed
	}
	if !on {
		return
	}
	switch state {
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
	row, org := seq.ConsolidatedRow(head), seq.Row(head)
	SetString(6, 0, row.Digitone.Pattern.String(), color(false, org.Digitone.Pattern != -1), termbox.ColorDefault)
	SetString(6, 1, row.Digitone.Mute.Format(row.Digitone.Channels), color(false, len(org.Digitone.Mute) != 0), termbox.ColorDefault)
	SetString(6+row.Digitone.Channels.Len+1, 0, row.Digitakt.Pattern.String(), color(false, org.Digitakt.Pattern != -1), termbox.ColorDefault)
	SetString(6+row.Digitone.Channels.Len+1, 1, row.Digitakt.Mute.Format(row.Digitakt.Channels), color(false, len(org.Digitakt.Mute) != 0), termbox.ColorDefault)
	SetString(0, 3, fmt.Sprintf("%s%02d", string('A'+head/64/16), 1+(head/64)%16), termbox.ColorDefault, termbox.ColorDefault)
	for n := 0; n < 16; n++ {
		ch := seq.Row(head/16*16 + n).HasChanges(seq.Row(head/16*16 + n - 1))
		SetString(6+(n%8)*3, 3+n/8, fmt.Sprintf("%02d", 1+n), color(n == head%16, ch), termbox.ColorDefault)
	}
	for n := 0; n < 4; n++ {
		SetString(32+n*4, 3, fmt.Sprintf("%d:4", 1+n), color(n == (head%64)/16, false), termbox.ColorDefault)
	}
	termbox.Flush()
}

type State int

const (
	Viewing State = iota
	Recording
	Playing
)
