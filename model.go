package main

const (
	pageSize = 64
	maxStep  = 64 * pageSize
)

type Model struct {
	step     int
	Digitakt *Synth
	Digitone *Synth
}

func NewModel() *Model {
	return &Model{0, NewSynth(8), NewSynth(4)}
}

func (m *Model) Step() int {
	return m.step
}

func (m *Model) PageSize() int {
	return pageSize
}

func (m *Model) Page() int {
	return m.step / pageSize
}

func (m *Model) Cursor() int {
	return m.step % pageSize
}

func (m *Model) StepForCursor(cursor int) int {
	return m.Page()*m.PageSize() + cursor%m.PageSize()
}

func (m *Model) CanDecPage() bool {
	return m.step >= pageSize
}

func (m *Model) DecPage() {
	if m.CanDecPage() {
		m.step -= pageSize
	}
}

func (m *Model) CanIncPage() bool {
	return m.step < maxStep-pageSize
}

func (m *Model) IncPage() {
	if m.CanIncPage() {
		m.step += pageSize
	}
}

func (m *Model) SetCursor(val int) {
	m.step = m.Page()*m.PageSize() + val%m.PageSize()
}

func (m *Model) HasMessage(step int) bool {
	_, ch := m.Digitakt.Pattern(step)
	if ch {
		return true
	}
	_, ch = m.Digitone.Pattern(step)
	if ch {
		return true
	}
	return false
}
