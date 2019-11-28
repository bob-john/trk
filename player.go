package main

import (
	"bytes"
	"log"
	"trk/track"

	"github.com/gomidi/midi"
	"github.com/gomidi/midi/midimessage/channel"
	"gitlab.com/gomidi/midi/mid"
)

type Player struct {
	ports map[string]mid.Out
	last  *playEventSet
	next  *playEventSet
}

func NewPlayer() *Player {
	return &Player{make(map[string]mid.Out), new(playEventSet), new(playEventSet)}
}

func (p *Player) Play(trk *track.Track, tick int) {
	parts, _ := track.Parts(trk)
	for _, part := range parts {
		p.writePattern(trk, part, tick)
		p.writeMute(trk, part, tick)
	}
	p.flush()
}

func (p *Player) PlayPattern(trk *track.Track, tick int) {
	parts, _ := track.Parts(trk)
	for _, part := range parts {
		p.writePattern(trk, part, tick)
	}
	p.flush()
}

func (p *Player) PlayMute(trk *track.Track, tick int) {
	parts, _ := track.Parts(trk)
	for _, part := range parts {
		p.writeMute(trk, part, tick)
	}
	p.flush()
}

func (p *Player) writePattern(trk *track.Track, part *track.Part, tick int) {
	pattern := track.Pattern(trk, part, tick)
	p.write(part.ProgChgPortOut, p.pattern(part, pattern))
}

func (p *Player) pattern(part *track.Part, pattern int) midi.Message {
	return channel.Channel(part.ProgChgOutCh).ProgramChange(uint8(pattern))
}

func (p *Player) writeMute(trk *track.Track, part *track.Part, tick int) {
	mute := track.Mute(trk, part, tick)
	for n, ch := range part.TrackCh {
		if ch == -1 {
			continue
		}
		p.write(part.MutePortOut, p.mute(part, ch, mute[n]))
	}
}

func (p *Player) mute(part *track.Part, ch int, muted bool) midi.Message {
	if muted {
		return channel.Channel(ch).ControlChange(94, 1)
	}
	return channel.Channel(ch).ControlChange(94, 0)
}

func (p *Player) write(ports []string, message midi.Message) {
	p.next.Insert(ports, message)
}

func (p *Player) flush() {
	curr := p.next.Substract(p.last)
	log.Printf("player flush %d events", len(curr.Events))
	for _, e := range curr.Events {
		port, ok := p.ports[e.Port]
		if !ok {
			var err error
			port, err = mid.OpenOut(driver, -1, e.Port)
			if err != nil {
				continue
			}
			p.ports[e.Port] = port
		}
		port.Send(e.Message)
	}

	required := make(map[string]bool)
	for _, e := range p.next.Events {
		required[e.Port] = true
	}
	for name, port := range p.ports {
		if !required[name] {
			if port.Close() == nil {
				delete(p.ports, name)
			}
		}
	}

	p.last, p.next = p.next, new(playEventSet)
}

type playEvent struct {
	Port    string
	Message []byte
}

func (e playEvent) Equals(o playEvent) bool {
	return e.Port == o.Port && bytes.Equal(e.Message, o.Message)
}

type playEventSet struct {
	Events []playEvent
}

func (s *playEventSet) Insert(ports []string, message midi.Message) {
	for _, port := range ports {
		e := playEvent{port, message.Raw()}
		if s.contains(e) {
			return
		}
		s.Events = append(s.Events, e)
	}
}

func (s *playEventSet) Substract(o *playEventSet) *playEventSet {
	r := new(playEventSet)
	for _, e := range s.Events {
		if o.contains(e) {
			continue
		}
		r.Events = append(r.Events, e)
	}
	return r
}

func (s *playEventSet) contains(e playEvent) bool {
	for _, o := range s.Events {
		if o.Equals(e) {
			return true
		}
	}
	return false
}
