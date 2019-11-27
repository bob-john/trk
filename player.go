package main

import (
	"trk/track"

	"github.com/gomidi/midi"
	"github.com/gomidi/midi/midimessage/channel"
	"gitlab.com/gomidi/midi/mid"
)

type Player struct {
	ports map[string]mid.Out
}

func NewPlayer() *Player {
	return &Player{make(map[string]mid.Out)}
}

func (p *Player) Play(trk *track.Track, tick int) {
	parts, _ := track.Parts(trk)
	for _, part := range parts {
		p.playPattern(trk, part, tick)
		p.playMute(trk, part, tick)
	}
}

func (p *Player) PlayPattern(trk *track.Track, tick int) {
	parts, _ := track.Parts(trk)
	for _, part := range parts {
		p.playPattern(trk, part, tick)
	}
}

func (p *Player) PlayMute(trk *track.Track, tick int) {
	parts, _ := track.Parts(trk)
	for _, part := range parts {
		p.playMute(trk, part, tick)
	}
}

func (p *Player) playPattern(trk *track.Track, part *track.Part, tick int) {
	pattern := track.Pattern(trk, part, tick)
	p.write(part.ProgChgPortOut, channel.Channel(part.ProgChgOutCh).ProgramChange(uint8(pattern)))
}

func (p *Player) playMute(trk *track.Track, part *track.Part, tick int) {
	mute := track.Mute(trk, part, tick)
	for n, ch := range part.TrackCh {
		if ch == -1 {
			continue
		}
		if mute[n] {
			p.write(part.MutePortOut, channel.Channel(ch).ControlChange(94, 1))
		} else {
			p.write(part.MutePortOut, channel.Channel(ch).ControlChange(94, 0))
		}
	}
}

func (p *Player) write(ports []string, message midi.Message) {
	required := make(map[string]bool)
	for _, name := range ports {
		required[name] = true
		out, ok := p.ports[name]
		if !ok {
			var err error
			out, err = mid.OpenOut(driver, -1, name)
			if err != nil {
				continue
			}
			p.ports[name] = out
		}
		out.Send(message.Raw())
	}
	for name, port := range p.ports {
		if !required[name] {
			if port.Close() == nil {
				delete(p.ports, name)
			}
		}
	}
}
