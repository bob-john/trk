package track

type Route struct {
	ID     string
	Input  string
	Output string
	Filter Filter
}

func (o *Route) Save(trk *Track) error {
	if o.ID == "" {
		o.ID = makeID()
	}
	return trk.db.Save(o)
}

type Filter int

const (
	Note Filter = 1 << iota
	PolyphonicAftertouch
	ControlChange
	ProgramChange
	ChannelAftertouch
	PitchBendChange
)

func Filters() []Filter {
	return []Filter{Note, PolyphonicAftertouch, ControlChange, ProgramChange, ChannelAftertouch, PitchBendChange}
}

func (f Filter) String() string {
	switch f {
	case Note:
		return "Note"
	case PolyphonicAftertouch:
		return "Polyphonic Aftertouch"
	case ControlChange:
		return "Control Change"
	case ProgramChange:
		return "Program Change"
	case ChannelAftertouch:
		return "Channel Aftertouch"
	case PitchBendChange:
		return "Pitch Bend Change"
	}
	return ""
}
