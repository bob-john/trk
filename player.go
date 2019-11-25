package main

import (
	"github.com/gomidi/midi"
	"github.com/gomidi/midi/midimessage/channel"
)

type Player struct {
	output map[string]*Output
}

func NewPlayer() *Player {
	return &Player{make(map[string]*Output)}
}

func (p *Player) Play(track *Track, row int) {
	p.PlayPattern(track, row)
	p.PlayMute(track, row)
}

func (p *Player) PlayPattern(track *Track, step int) {
	row := track.Seq.ConsolidatedRow(step)
	for name, part := range row.Parts {
		dev, ok := track.Settings.Devices[name]
		if !ok {
			continue
		}
		ch := dev.ProgChgOutCh - 1
		if ch < 0 {
			continue
		}
		p.Write(dev.Outputs, channel.Channel(ch).ProgramChange(uint8(part.Pattern)))
	}
}

func (p *Player) PlayMute(track *Track, step int) {
	row := track.Seq.ConsolidatedRow(step)
	for name, part := range row.Parts {
		device, ok := track.Settings.Devices[name]
		if !ok {
			continue
		}
		//FIXME Drop source. Use device names.
		if device.MuteSrc == DeviceSourceBoth {
			for _, device := range track.Settings.Devices {
				p.playMute(part.Mute, device)
			}
		} else {
			p.playMute(part.Mute, device)
		}
	}
}

func (p *Player) playMute(mute Mute, device *DeviceSettings) {
	for n, ch := range device.Channels {
		ch := ch - 1
		if ch < 0 {
			continue
		}
		if mute[n] {
			p.Write(device.Outputs, channel.Channel(ch).ControlChange(94, 1))
		} else {
			p.Write(device.Outputs, channel.Channel(ch).ControlChange(94, 0))
		}
	}
}

func (p *Player) Write(names map[string]struct{}, message midi.Message) {
	var err error
	required := make(map[string]bool)
	for name := range names {
		required[name] = true
		port, ok := p.output[name]
		if !ok {
			port, err = OpenOutput(name)
			if err != nil {
				continue
			}
			p.output[name] = port
		}
		port.Write(message)
	}
	for name, port := range p.output {
		if !required[name] {
			port.Close()
			delete(p.output, name)
		}
	}
}
