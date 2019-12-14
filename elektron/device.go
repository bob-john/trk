package elektron

import (
	"trk/tracker"

	"github.com/gomidi/midi/midimessage/channel"
	"gitlab.com/gomidi/midi"
)

type Device struct {
	port         string
	channels     map[int]int
	progChgOutCh int
	trackCount   int
	muteInit     bool
}

func NewDevice(port string, trackCount int) *Device {
	return &Device{port, make(map[int]int), 10, trackCount, true}
}

func (d *Device) SetChannel(track, channel int) {
	d.channels[track] = channel
}

func (d *Device) SetProgChgOutCh(channel int) {
	d.progChgOutCh = channel
}

func (d *Device) SetPattern(ptn Pattern) {
	tracker.Play(progChg{d.port, d.progChgOutCh, ptn.Program()})
}

func (d *Device) Mute(tracks ...int) {
	if d.muteInit {
		d.muteInit = false
		d.Unmute()
	}
	if len(tracks) == 0 {
		for i := 1; i <= d.trackCount; i++ {
			tracks = append(tracks, i)
		}
	}
	for _, track := range tracks {
		tracker.Play(mute{d.port, d.channel(track)})
	}
}

func (d *Device) Unmute(tracks ...int) {
	if d.muteInit {
		d.muteInit = false
		d.Mute()
	}
	if len(tracks) == 0 {
		for i := 1; i <= d.trackCount; i++ {
			tracks = append(tracks, i)
		}
	}
	for _, track := range tracks {
		tracker.Play(unmute{d.port, d.channel(track)})
	}
}

func (d *Device) channel(track int) int {
	ch, ok := d.channels[track]
	if ok {
		return ch
	}
	return track
}

type mute struct {
	port    string
	channel int
}

func (m mute) Port() string {
	return m.port
}

func (m mute) Message() midi.Message {
	return channel.Channel(m.channel-1).ControlChange(94, 1)
}

type unmute struct {
	port    string
	channel int
}

func (u unmute) Port() string {
	return u.port
}

func (u unmute) Message() midi.Message {
	return channel.Channel(u.channel-1).ControlChange(94, 0)
}

type progChg struct {
	port    string
	channel int
	program uint8
}

func (p progChg) Port() string {
	return p.port
}

func (p progChg) Message() midi.Message {
	return channel.Channel(p.channel - 1).ProgramChange(p.program)
}
