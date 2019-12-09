package track

type EventStringer struct {
	Event *Event
}

func (s EventStringer) Channel() string {
	return ""
}

func (s EventStringer) Type() string {
	return ""
}

func (s EventStringer) Subtype() string {
	return ""
}

func (s EventStringer) Value() string {
	return ""
}
