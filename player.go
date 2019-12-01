package main

import (
	"bytes"
	"log"
	"trk/track"

	"gitlab.com/gomidi/midi"
	"gitlab.com/gomidi/midi/mid"
	"gitlab.com/gomidi/midi/midimessage/channel"
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
	for _, part := range trk.Parts() {
		p.writePattern(trk, part, tick)
		p.writeMute(trk, part, tick)
	}
	p.flush()
}

func (p *Player) PlayPattern(trk *track.Track, tick int) {
	for _, part := range trk.Parts() {
		p.writePattern(trk, part, tick)
	}
	p.flush()
}

func (p *Player) PlayMute(trk *track.Track, tick int) {
	for _, part := range trk.Parts() {
		p.writeMute(trk, part, tick)
	}
	p.flush()
}

func (p *Player) Close() {
	for name := range p.ports {
		p.close(name)
	}
}

func (p *Player) writePattern(trk *track.Track, part *track.Part, tick int) {
	pattern := trk.Pattern(part, tick)
	p.write(part.ProgChgPortOut, p.pattern(part, pattern))
}

func (p *Player) pattern(part *track.Part, pattern int) midi.Message {
	return channel.Channel(part.ProgChgOutCh).ProgramChange(uint8(pattern))
}

func (p *Player) writeMute(trk *track.Track, part *track.Part, tick int) {
	mute := trk.Mute(part, tick)
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
	for _, e := range curr.Events {
		port, ok := p.ports[e.Port]
		if !ok {
			var err error
			port, err = mid.OpenOut(midiDriver, -1, e.Port)
			log.Printf("player: open %s: %v", port, err)
			if err != nil {
				continue
			}
			p.ports[e.Port] = port
		}
		port.Send(e.Message.Raw())
	}

	required := make(map[string]bool)
	for _, e := range p.next.Events {
		required[e.Port] = true
	}
	for name := range p.ports {
		if !required[name] {
			p.close(name)
		}
	}

	p.last.Merge(p.next)
	p.next.Clear()
}

func (p *Player) close(name string) {
	port, ok := p.ports[name]
	if !ok {
		return
	}
	err := port.Close()
	log.Printf("player: close %s: %v", port, err)
	if err != nil {
		return
	}
	delete(p.ports, name)
}

type playEvent struct {
	Port    string
	Message midi.Message
}

func (e playEvent) Equals(o playEvent) bool {
	return e.Port == o.Port && bytes.Equal(e.Message.Raw(), o.Message.Raw())
}

func (e playEvent) Replace(o playEvent) bool {
	if e.Port != o.Port {
		return false
	}
	switch e := e.Message.(type) {
	case channel.ProgramChange:
		o, ok := o.Message.(channel.ProgramChange)
		if !ok {
			return false
		}
		return e.Channel() == o.Channel()
	case channel.ControlChange:
		o, ok := o.Message.(channel.ControlChange)
		if !ok {
			return false
		}
		return e.Channel() == o.Channel() && e.Controller() == o.Controller()
	}
	return e.Equals(o)
}

type playEventSet struct {
	Events []playEvent
}

func (s *playEventSet) Insert(ports []string, message midi.Message) {
	for _, port := range ports {
		e := playEvent{port, message}
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

func (s *playEventSet) Merge(o *playEventSet) {
	for _, o := range o.Events {
		var replaced bool
		for i, e := range s.Events {
			if o.Replace(e) {
				s.Events[i] = o
				replaced = true
				break
			}
		}
		if !replaced {
			s.Events = append(s.Events, o)
		}
	}
}

func (s *playEventSet) Clear() {
	s.Events = nil
}

func (s *playEventSet) contains(e playEvent) bool {
	for _, o := range s.Events {
		if o.Equals(e) {
			return true
		}
	}
	return false
}
