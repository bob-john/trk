package main

import (
	"github.com/gomidi/midi"
	"github.com/gomidi/midi/midimessage/channel"
	"github.com/gomidi/midi/midimessage/sysex"
)

type Launchpad struct {
	*Device
	color  map[uint8]uint8
	color2 map[uint8]*uint8
	dirty  map[uint8]bool
}

func ConnectLaunchpad() (*Launchpad, error) {
	d, err := ConnectDevice("MIDIIN2 (LPMiniMK3 MIDI)", "MIDIOUT2 (LPMiniMK3 MIDI)")
	// d, err := ConnectDevice("MIDIIN2 (LPMiniMK3 MIDI)", "loopMIDI Port")
	if err != nil {
		return nil, err
	}
	return &Launchpad{d, make(map[uint8]uint8), make(map[uint8]*uint8), make(map[uint8]bool)}, nil
}

func (lp *Launchpad) Reset() {
	// lp.Write(sysex.SysEx{0, 32, 41, 2, 24, 14, 0})
	for i := uint8(1); i < 10; i++ {
		for j := uint8(1); j < 10; j++ {
			lp.Write(channel.Channel0.NoteOn(i*10+j, 0))
		}
	}
	lp.color = make(map[uint8]uint8)
	lp.color2 = make(map[uint8]*uint8)
	lp.dirty = make(map[uint8]bool)

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
			lp.Write(channel.Channel0.NoteOn(loc, lp.color[loc]))
			if color := lp.color2[loc]; color != nil {
				lp.Write(channel.Channel1.NoteOn(loc, *color))
			}
		}
		delete(lp.dirty, loc)
	}
}
func (lp *Launchpad) Clear() {
	for i := uint8(1); i < 10; i++ {
		for j := uint8(1); j < 10; j++ {
			lp.Draw(i, j, 0)
			lp.SetFlashing(i, j, nil)
		}
	}
}

func (lp *Launchpad) Draw(row, col, color uint8) {
	loc := row*10 + col
	lp.dirty[loc] = lp.dirty[loc] || (lp.color[loc] != color)
	lp.color[loc] = color
}

func (lp *Launchpad) DrawHorizontalLine(row, col, count, color uint8) {
	for i := uint8(0); i < count; i++ {
		lp.Draw(row, col+i, color)
	}
}

func (lp *Launchpad) SetFlashing(row, col uint8, color *uint8) {
	loc := row*10 + col
	lp.dirty[loc] = lp.dirty[loc] || (lp.color2[loc] != color)
	lp.color2[loc] = color
}

func (lp *Launchpad) Loc(m midi.Message) uint8 {
	switch m := m.(type) {
	case channel.NoteOn:
		return m.Key()
	case channel.NoteOff:
		return m.Key()
	case channel.ControlChange:
		return m.Controller()
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

func (lp *Launchpad) Row(m midi.Message) uint8 {
	return lp.Loc(m) / 10
}

func (lp *Launchpad) Col(m midi.Message) uint8 {
	return lp.Loc(m) % 10
}

func (lp *Launchpad) ClearNavigationButtons() {
	lp.DrawHorizontalLine(9, 1, 4, 0)
}

func (lp *Launchpad) ClearModeButtons() {
	lp.DrawHorizontalLine(9, 5, 4, 0)
}

func (lp *Launchpad) ClearGrid() {
	for i := uint8(0); i < 8; i++ {
		for j := uint8(0); j < 8; j++ {
			lp.Draw(i, j, 0)
		}
	}
}
