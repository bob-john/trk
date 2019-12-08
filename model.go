package main

import "trk/track"

type Model struct {
	Track *track.Track
	Head  int
	State State
}

func NewModel() *Model {
	return new(Model)
}

func (m *Model) Page() int {
	return m.Head / 16
}

func (m *Model) SetPage(val int) {
	m.setHead(Clamp(val, 0, m.LastPage()), m.X(), m.Y())
}

func (m *Model) X() int {
	return (m.Head % 16) % 8
}

func (m *Model) SetX(val int) {
	m.setHead(m.Page(), val, m.Y())
}

func (m *Model) Y() int {
	return (m.Head % 16) / 8
}

func (m *Model) SetY(val int) {
	if m.Page() == 0 && m.Y() == 0 && val < 0 {
		return
	}
	if m.Page() == m.LastPage() && m.Y() == 1 && val > 1 {
		return
	}
	m.setHead(m.Page(), m.X(), val)
}

func (m *Model) LastPage() int {
	return 512 - 1
}

func (m *Model) LastStep() int {
	return (m.LastPage()+1)*16 - 1
}

func (m *Model) ToggleRecording() {
	switch m.State {
	case Viewing:
		m.State = Recording
	case Recording:
		m.State = Viewing
	}
}

func (m *Model) HeadForTrig(val int) int {
	return m.makeHead(m.Page(), val%8, val/8)
}

func (m *Model) SetTrig(val int) {
	m.setHead(m.Page(), val%8, val/8)
}

func (m *Model) SetHead(step int) {
	m.setHead(0, step, 0)
}

func (m *Model) setHead(pattern, x, y int) {
	if m.State.Is(Viewing, Recording) {
		m.Head = m.makeHead(pattern, x, y)
		m.State = Viewing
	}
}

func (m *Model) makeHead(pattern, x, y int) int {
	return Clamp(pattern*16+y*8+x, 0, 512*16-1)
}

type State int

const (
	Viewing State = iota
	Recording
	Playing
)

func (s State) Is(values ...State) bool {
	for _, val := range values {
		if s == val {
			return true
		}
	}
	return false
}
