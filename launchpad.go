package main

import (
	"github.com/gomidi/midi"
	"github.com/gomidi/midi/midimessage/channel"
)

type Launchpad struct {
	*Device
	color map[uint8]uint8
	dirty map[uint8]bool
}

func ConnectLaunchpad() (*Launchpad, error) {
	d, err := ConnectDevice("MIDIIN2 (LPMiniMK3 MIDI)", "MIDIOUT2 (LPMiniMK3 MIDI)")
	if err != nil {
		return nil, err
	}
	return &Launchpad{d, make(map[uint8]uint8), make(map[uint8]bool)}, nil
}

func (lp *Launchpad) Reset() {
	for i := uint8(1); i < 10; i++ {
		for j := uint8(1); j < 10; j++ {
			lp.Write(channel.Channel0.NoteOn(i*10+j, 0))
		}
	}
	lp.color = make(map[uint8]uint8)
	lp.dirty = make(map[uint8]bool)
}

func (lp *Launchpad) Clear() {
	for i := uint8(1); i < 10; i++ {
		for j := uint8(1); j < 10; j++ {
			lp.Set(i, j, 0)
		}
	}
}

func (lp *Launchpad) Set(row, col, color uint8) {
	loc := row*10 + col
	lp.dirty[loc] = lp.dirty[loc] || (lp.color[loc] != color)
	lp.color[loc] = color
}

func (lp *Launchpad) Get(row, col uint8) uint8 {
	loc := row*10 + col
	return lp.color[loc]
}

func (lp *Launchpad) StartFlashing(row, col, color uint8) {
	loc := row*10 + col
	lp.Write(channel.Channel1.NoteOn(loc, color))
}

func (lp *Launchpad) StopFlashing(row, col uint8) {
	loc := row*10 + col
	lp.Write(channel.Channel0.NoteOn(loc, lp.color[loc]))
	delete(lp.dirty, loc)
}

func (lp *Launchpad) Update() {
	// start := time.Now()
	// defer func() {
	// 	fmt.Println("Update", time.Since(start))
	// }()
	for loc, dirty := range lp.dirty {
		if dirty {
			lp.Write(channel.Channel0.NoteOn(loc, lp.color[loc]))
		}
		delete(lp.dirty, loc)
	}
}

func (lp *Launchpad) Location(m midi.Message) uint8 {
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
