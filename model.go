package main

const (
	pageSize = 64
	maxStep  = 64 * pageSize
)

type Model struct {
	step   int
	tracks []*Track
}

func NewModel() *Model {
	m := new(Model)
	m.tracks = append(m.tracks, NewTrack(8))
	m.tracks = append(m.tracks, NewTrack(4))
	return m
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

func (m *Model) Tracks() []*Track {
	return m.tracks
}

func (m *Model) SetCursor(val int) {
	m.step = m.Page()*m.PageSize() + val%m.PageSize()
}
