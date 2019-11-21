package main

import "os"

type Model struct {
	Seq   *Seq
	Head  int
	State State
}

func (m *Model) LoadSeq(path string) error {
	var err error
	m.Seq, err = ReadSeq(path)
	if os.IsNotExist(err) {
		m.Seq = NewSeq()
	} else if err != nil {
		return err
	}
	return nil
}

func (m *Model) HeadForTrig(val int) int {
	return m.makeHead(m.Pattern(), val)
}

func (m *Model) Pattern() int {
	return m.Head / 16
}

func (m *Model) Trig() int {
	return m.Head % 16
}

func (m *Model) SetPattern(val int) {
	m.setHead(clamp(val, 0, m.LastPattern()), m.Trig())
}

func (m *Model) SetTrig(val int) {
	m.setHead(m.Pattern(), clamp(val, 0, 16-1))
}

func (m *Model) LastPattern() int {
	return 512 - 1
}

func (m *Model) ClearStep() {
	if m.State.Is(Viewing, Playing) {
		m.Seq.Clear(m.Head)
	}
}

func (m *Model) ToggleRecording() {
	switch m.State {
	case Viewing:
		m.State = Recording
	case Recording:
		m.State = Viewing
	}
}

func (m *Model) setHead(pattern, trig int) {
	if m.State.Is(Viewing, Recording) {
		m.Head = m.makeHead(pattern, trig)
		m.State = Viewing
	}
}

func (m *Model) makeHead(pattern, trig int) int {
	return clamp(pattern*16+trig, 0, 512*16-1)
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
