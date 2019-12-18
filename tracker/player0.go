package tracker

import (
	"bytes"
	"log"

	"gitlab.com/gomidi/midi"
	"gitlab.com/gomidi/midi/mid"
)

type Player0 struct {
	driver mid.Driver
	ports  map[string]mid.Out
	last   *playEventSet
	next   *playEventSet
}

func NewPlayer0(driver mid.Driver) *Player0 {
	return &Player0{driver, make(map[string]mid.Out), new(playEventSet), new(playEventSet)}
}

func (p *Player0) Close() {
	for name := range p.ports {
		p.close(name)
	}
}

func (p *Player0) Play(port string, message midi.Message) error {
	out, ok := p.ports[port]
	if !ok {
		var err error
		out, err = mid.OpenOut(p.driver, -1, port)
		if err != nil {
			return err
		}
		p.ports[port] = out
	}
	out.Send(message.Raw())
	return nil
}

func (p *Player0) Queue(ports []string, message midi.Message) {
	p.next.Insert(ports, message)
}

func (p *Player0) Flush() {
	curr := p.next.Substract(p.last)
	for _, e := range curr.Events {
		port, ok := p.ports[e.Port]
		if !ok {
			var err error
			port, err = mid.OpenOut(p.driver, -1, e.Port)
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

func (p *Player0) close(name string) {
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
	return false
	// if e.Port != o.Port {
	// 	return false
	// }
	// switch e := e.Message.(type) {
	// case channel.ProgramChange:
	// 	o, ok := o.Message.(channel.ProgramChange)
	// 	if !ok {
	// 		return false
	// 	}
	// 	return e.Channel() == o.Channel()
	// case channel.ControlChange:
	// 	o, ok := o.Message.(channel.ControlChange)
	// 	if !ok {
	// 		return false
	// 	}
	// 	return e.Channel() == o.Channel() && e.Controller() == o.Controller()
	// }
	// return e.Equals(o)
}

type playEventSet struct {
	Events []playEvent
}

func (s *playEventSet) Insert(ports []string, message midi.Message) {
	for _, port := range ports {
		e := playEvent{port, message}
		if s.contains(e) {
			continue
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
