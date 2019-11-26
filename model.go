package main

type Model struct {
	Seq0  *Seq0
	Track *Track
	Head  int
	State State
}

func NewModel() *Model {
	return &Model{Seq0: NewSeq0()}
}

func (m *Model) LoadTrack(path string) error {
	m.Track = NewTrack()
	// var err error
	// m.Track, err = ReadTrack(path)
	// if os.IsNotExist(err) {
	// 	m.Track = NewTrack()
	// } else if err != nil {
	// 	return err
	// }
	return nil
}

func (m *Model) Pattern() int {
	return m.Head / 16
}

func (m *Model) SetPattern(val int) {
	m.setHead(clamp(val, 0, m.LastPattern()), m.X(), m.Y())
}

func (m *Model) X() int {
	return (m.Head % 16) % 8
}

func (m *Model) SetX(val int) {
	m.setHead(m.Pattern(), val, m.Y())
}

func (m *Model) Y() int {
	return (m.Head % 16) / 8
}

func (m *Model) SetY(val int) {
	if m.Pattern() == 0 && m.Y() == 0 && val < 0 {
		return
	}
	if m.Pattern() == m.LastPattern() && m.Y() == 1 && val > 1 {
		return
	}
	m.setHead(m.Pattern(), m.X(), val)
}

func (m *Model) LastPattern() int {
	return 512 - 1
}

func (m *Model) ClearStep() {
	if m.State.Is(Viewing, Playing) {
		m.Track.Seq.Clear(m.Head)
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

func (m *Model) HeadForTrig(val int) int {
	return m.makeHead(m.Pattern(), val%8, val/8)
}

func (m *Model) SetTrig(val int) {
	m.setHead(m.Pattern(), val%8, val/8)
}

func (m *Model) setHead(pattern, x, y int) {
	if m.State.Is(Viewing, Recording) {
		m.Head = m.makeHead(pattern, x, y)
		m.State = Viewing
	}
}

func (m *Model) makeHead(pattern, x, y int) int {
	return clamp(pattern*16+y*8+x, 0, 512*16-1)
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
