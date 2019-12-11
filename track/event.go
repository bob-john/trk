package track

import "fmt"

type Event struct {
	Tick    int
	Port    string
	Message []byte
}

func (e Event) Beat() string {
	return fmt.Sprintf("%3d.%02d", 1+(e.Tick/(4*24)), 1+(e.Tick*4/24)%16)
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

func (e Event) Subtype() string {
	if len(e.Message) < 1 {
		return "-"
	}
	switch e.Message[0] & 0xf0 {
	case 0xb0:
		if len(e.Message) == 3 {
			return fmt.Sprintf("CC %02d", e.Message[1])
		}
	}
	return "-"
}

func (e Event) Value() int {
	if len(e.Message) < 1 {
		return 0
	}
	switch e.Message[0] & 0xf0 {
	case 0x80, 0x90, 0xa0, 0xb0:
		if len(e.Message) == 3 {
			return int(e.Message[2] & 0x7f)
		}
	case 0xc0, 0xd0:
		if len(e.Message) == 2 {
			return int(e.Message[1] & 0x7f)
		}
	case 0xe0:
		if len(e.Message) == 3 {
			return int((e.Message[1]&0x7f)<<7 | (e.Message[2] & 0x7f))
		}
	}
	return 0
}
