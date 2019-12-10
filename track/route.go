package track

import "fmt"

type Route struct {
	Input  string // Input device name (!= port)
	Output string // Output device name (!= port)
	ProgCh bool   // Forward program changes
	Notes  bool   // Forward notes and related
	CC     bool   // Forward CC
}

func (r *Route) String() string {
	return fmt.Sprintf("%s -> %s", r.Input, r.Output)
}

func (r *Route) Accept(message []byte) bool {
	if len(message) == 0 {
		return false
	}
	switch message[0] & 0xf0 {
	case 0x80, 0x90, 0xa0, 0xd0, 0xe0:
		return r.Notes
	case 0xb0:
		return r.CC
	case 0xc0:
		return r.ProgCh
	case 0xf0:
		return false
	}
	return false
}
