package main

import (
	"github.com/gomidi/midi"
	"github.com/gomidi/midi/midimessage/channel"
	"github.com/gomidi/midi/midimessage/sysex"
)

type Launchpad struct {
	*Device
	color  map[int]int
	color2 map[int]*int
	dirty  map[int]bool
}

func ConnectLaunchpad() (*Launchpad, error) {
	d, err := ConnectDevice("MIDIIN2 (LPMiniMK3 MIDI)", "MIDIOUT2 (LPMiniMK3 MIDI)")
	// d, err := ConnectDevice("MIDIIN2 (LPMiniMK3 MIDI)", "loopMIDI Port")
	if err != nil {
		return nil, err
	}
	return &Launchpad{d, make(map[int]int), make(map[int]*int), make(map[int]bool)}, nil
}

func (lp *Launchpad) Reset() {
	// lp.Write(sysex.SysEx{0, 32, 41, 2, 24, 14, 0})
	for i := 1; i < 10; i++ {
		for j := 1; j < 10; j++ {
			lp.Write(channel.Channel0.NoteOn(uint8(i*10+j), 0))
		}
	}
	lp.color = make(map[int]int)
	lp.color2 = make(map[int]*int)
	lp.dirty = make(map[int]bool)

	//	lp.Write(sysex.SysEx{0, 32, 41, 2, 24, 11, 7, 60, 60, 30})
	lp.Write(sysex.SysEx{126, 127, 6, 1})
}

func (lp *Launchpad) Flush() {
	// start := time.Now()
	// defer func() {
	// 	fmt.Println("Flush", time.Since(start))
	// }()
	for loc, dirty := range lp.dirty {
		if dirty {
			lp.Write(channel.Channel0.NoteOn(uint8(loc), uint8(lp.color[loc])))
			if color := lp.color2[loc]; color != nil {
				lp.Write(channel.Channel1.NoteOn(uint8(loc), uint8(*color)))
			}
		}
		delete(lp.dirty, loc)
	}
}
func (lp *Launchpad) Clear() {
	for i := 1; i < 10; i++ {
		for j := 1; j < 10; j++ {
			lp.Draw(i, j, 0)
			lp.SetFlashing(i, j, nil)
		}
	}
}

func (lp *Launchpad) Draw(row, col, color int) {
	loc := row*10 + col
	lp.dirty[loc] = lp.dirty[loc] || (lp.color[loc] != color)
	lp.color[loc] = color
}

func (lp *Launchpad) DrawHorizontalLine(row, col, count, color int) {
	for i := 0; i < count; i++ {
		lp.Draw(row, col+i, color)
	}
}

func (lp *Launchpad) SetFlashing(row, col int, color *int) {
	loc := row*10 + col
	lp.dirty[loc] = lp.dirty[loc] || (lp.color2[loc] != color)
	lp.color2[loc] = color
}

func (lp *Launchpad) Loc(m midi.Message) int {
	switch m := m.(type) {
	case channel.NoteOn:
		return int(m.Key())
	case channel.NoteOff:
		return int(m.Key())
	case channel.ControlChange:
		return int(m.Controller())
	}
	return 0
}

func (lp *Launchpad) IsOn(m midi.Message) bool {
	switch m := m.(type) {
	case channel.NoteOn:
		return true
	case channel.ControlChange:
		return m.Value() != 0
	}
	return false
}

func (lp *Launchpad) IsPad(m midi.Message) bool {
	switch m.(type) {
	case channel.NoteOn, channel.NoteOff:
		return true
	}
	return false
}

func (lp *Launchpad) Row(m midi.Message) int {
	return lp.Loc(m) / 10
}

func (lp *Launchpad) Col(m midi.Message) int {
	return lp.Loc(m) % 10
}

func (lp *Launchpad) ClearNavigationButtons() {
	lp.DrawHorizontalLine(9, 1, 4, 0)
}

func (lp *Launchpad) ClearModeButtons() {
	lp.DrawHorizontalLine(9, 5, 4, 0)
}

func (lp *Launchpad) ClearGrid() {
	for i := 1; i <= 8; i++ {
		for j := 1; j <= 8; j++ {
			lp.Draw(i, j, 0)
		}
	}
}
