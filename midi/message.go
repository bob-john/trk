package midi

// type Message interface{}

// func Split(data []byte, atEOF bool) (advance int, token []byte, err error) {
// 	advance = Len(data)
// 	if advance > 0 {

// 	}

// }

// func Len(data []byte) int {
// 	if len(data) == 0 {
// 		return 0
// 	}
// 	switch data[0] & 0xF0 {
// 	case 0x80, 0x90, 0xa0, 0xb0, 0xe0:
// 		return 3
// 	case 0xc0, 0xd0:
// 		return 2
// 	case 0xf0:
// 		switch data[0] {
// 		case 0xf0:
// 			//TODO SysEx
// 		case 0xf1:
// 			//TODO
// 		case 0xf2:
// 			return 3
// 		case 0xf3:
// 			return 2
// 		case 0xf4, 0xf5, 0xf9, 0xfd:
// 			for n, b := range data[1:] {
// 				if b&0x80 != 0 {
// 					return n
// 				}
// 			}
// 			return 0
// 		case 0xf6, 0xf7, 0xf8, 0xfa, 0xfb, 0xfc, 0xfe, 0xff:
// 			return 1
// 		}
// 	}
// 	return 0
// }

// func Parse(b []byte) Message {
// 	if len(b) != Len(b) {
// 		return nil
// 	}
// 	switch b[0] & 0xf0 {
// 	case 0x80:
// 		return NoteOff(b)
// 	case 0x90:
// 		return NoteOn(b)
// 	case 0xA0:
// 		return PolyphonicAftertouch(b)
// 	case 0xB0:
// 		return ControlChange(b)
// 	case 0xC0:
// 		return ProgramChange(b)
// 	case 0xD0:
// 		return ChannelAftertouch(b)
// 	case 0xE0:
// 		return PitchBendChange(b)
// 	case 0xF0:
// 		switch b[0] {
// 		case 0xF0:
// 			return SysEx(b)
// 		case 0xF1:
// 			return TimeCodeQtrFrame(b)
// 		case 0xF2:
// 			return SongPositionPointer(b)
// 		case 0xF3:
// 			return SongSelect(b)
// 		case 0xF4, 0xF5:
// 			return Undefined(b)
// 		case 0xF6:
// 			return TuneRequest(b)
// 		case 0xF7:
// 			return EndOfSysEx(b)
// 		case 0xF8:
// 			return Timingclock(b)
// 		case 0xF9:
// 			return Undefined(b)
// 		case 0xFA:
// 			return Start(b)
// 		case 0xFB:
// 			return Continue(b)
// 		case 0xFC:
// 			return Stop(b)
// 		case 0xFD:
// 			return Undefined(b)
// 		case 0xFE:
// 			return ActiveSensing(b)
// 		case 0xFF:
// 			return SystemReset(b)
// 		}
// 	}
// 	return nil
// }

// type NoteOff []byte
// type NoteOn []byte
// type PolyphonicAftertouch []byte
// type ControlChange []byte
// type ProgramChange []byte
// type ChannelAftertouch []byte
// type PitchBendChange []byte
// type SysEx []byte
// type TimeCodeQtrFrame []byte
// type SongPositionPointer []byte
// type SongSelect []byte
// type Undefined []byte
// type TuneRequest []byte
// type EndOfSysEx []byte
// type Timingclock []byte
// type Start []byte
// type Continue []byte
// type Stop []byte
// type ActiveSensing []byte
// type SystemReset []byte
