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
	return m.makeHead(m.Pattern(), m.Page(), val)
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
	if m.Pattern() == 0 && val < 0 {
		return
	}
	if m.Pattern() == m.LastPattern() && val > 3 {
		return
	}
	m.setHead(m.Pattern(), val, m.Trig())
}

func (m *Model) LastPattern() int {
	return 8*16 - 1
}

func (m *Model) LastPage() int {
	return 4 - 1
}

func (m *Model) Trig() int {
	return m.Head % 16
}

func (m *Model) SetTrig(val int) {
	if val < 0 || val > 15 {
		return
	}
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
		m.Head = m.makeHead(pattern, page, trig)
		m.State = Viewing
	}
}

func (m *Model) makeHead(pattern, page, trig int) int {
	return clamp(pattern*64+page*16+trig, 0, 8*16*64-1)
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
