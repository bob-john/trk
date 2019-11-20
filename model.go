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

func (m *Model) Pattern() int {
	return m.Head / 64
}

func (m *Model) SetPattern(val int) {
	m.setHead(val, m.Page(), m.Trig())
}

func (m *Model) Page() int {
	return (m.Head % 64) / 16
}

func (m *Model) SetPage(val int) {
	m.setHead(m.Pattern(), val, m.Trig())
}

func (m *Model) Trig() int {
	return m.Head % 16
}

func (m *Model) SetTrig(val int) {
	m.setHead(m.Pattern(), m.Page(), val)
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

func (m *Model) setHead(pattern, page, trig int) {
	if m.State.Is(Viewing, Recording) {
		m.Head = clamp(pattern*64+page*16+trig, 0, 8*16*64-1)
		m.State = Viewing
	}
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
