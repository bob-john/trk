package track

import (
	"fmt"
)

func FormatTrackName(part string, n int) string {
	switch part {
	case "DIGITAKT":
		if n < 8 {
			return fmt.Sprintf("TRACK %d", 1+n)
		}
		return fmt.Sprintf("TRACK %s", string('A'+n-8))

	case "DIGITONE":
		if n < 4 {
			return fmt.Sprintf("TRACK %d", 1+n)
		}
		return fmt.Sprintf("MIDI %d", 1+n-4)

	default:
		return fmt.Sprintf("TRACK %d", 1+n)
	}
}

func FormatPattern(p int) string {
	return fmt.Sprintf("%s%02d", string('A'+p/16), 1+p%16)
}

func FormatMute(mute [16]bool, part *Part) (str string) {
	for n, ch := range part.TrackCh {
		if ch < 0 || mute[n] {
			str += "-"
		} else if n < 8 {
			str += string('1' + n)
		} else {
			str += string('A' + n - 8)
		}
		if n == len(part.TrackCh)/2-1 {
			str += " "
		}
	}
	return
}
