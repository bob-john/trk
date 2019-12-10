package track

type Event struct {
	Tick    int
	Port    string
	Message []byte
}

func (e Event) Channel() int {
	if len(e.Message) < 0 {
		return 0
	}
	if e.Message[0] >= 0x80 && e.Message[0] < 0xff {
		return 1 + int(e.Message[0]&0x0f)
	}
	return 0
}
func (e Event) Type() string {
	if len(e.Message) < 0 {
		return ""
	}
	switch e.Message[0] & 0xf0 {
	case 0x80:
		return "Note Off"
	case 0x90:
		//TODO Check velocity
		return "Note On"
	case 0xa0:
		return "Polyphonic Aftertouch"
	case 0xb0:
		return "Control Change"
	case 0xc0:
		return "Program Change"
	case 0xd0:
		return "Channel Aftertouch"
	case 0xe0:
		return "Pitch Bend Change"
	}
	return ""
}

// func (e Event) Type() string {
// 	if len(e.Message) < 0 {
// 		return ""
// 	}
// 	switch e.Message[0] & 0xf0 {
// 	case 0x80, 0x90:
// 		return "Note Off"
// 	case 0x90:
// 		return "Note On"
// 	case 0xa0:
// 		return "Polyphonic Aftertouch"
// 	case 0xb0:
// 		return "Control Change"
// 	case 0xc0:
// 		return "Program Change"
// 	case 0xd0:
// 		return "Channel Aftertouch"
// 	case 0xe0:
// 		return "Pitch Bend Change"
// 	}
// 	return ""
// }
