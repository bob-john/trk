package track

import (
	"github.com/asdine/storm/q"
	"github.com/nsf/termbox-go"
)

func Parts(trk *Track) (parts []*Part, err error) {
	err = trk.db.All(&parts)
	return
}

func Pattern(trk *Track, part *Part, tick int) int {
	var pc PatternChange
	query := trk.db.Select(q.Eq("Part", part.Name), q.Lte("Tick", tick)).OrderBy("Tick").Reverse().Limit(1)
	query.First(&pc)
	return pc.Pattern
}

func IsPatternModified(trk *Track, part *Part, tick int) bool {
	var pc PatternChange
	n, err := trk.db.Select(q.Eq("Part", part.Name), q.Eq("Tick", tick)).Count(&pc)
	return err == nil && n != 0
}

func Mute(trk *Track, part *Part, tick int) [16]bool {
	var mc MuteChange
	query := trk.db.Select(q.Eq("Part", part.Name), q.Lte("Tick", tick)).OrderBy("Tick").Reverse().Limit(1)
	query.First(&mc)
	return mc.Mute
}

func IsMuteModified(trk *Track, part *Part, tick int) bool {
	var mc MuteChange
	n, err := trk.db.Select(q.Eq("Part", part.Name), q.Eq("Tick", tick)).Count(&mc)
	return err == nil && n != 0
}

func IsPartModified(trk *Track, part *Part, tick int) bool {
	return IsPatternModified(trk, part, tick) || IsMuteModified(trk, part, tick)
}

func IsModified(trk *Track, tick int) bool {
	parts, _ := Parts(trk)
	for _, part := range parts {
		if IsPatternModified(trk, part, tick) {
			return true
		}
		if IsMuteModified(trk, part, tick) {
			return true
		}
	}
	return false
}

func InputPorts(trk *Track) (ports []string) {
	parts, err := Parts(trk)
	if err != nil {
		return
	}
	for _, part := range parts {
		ports = append(ports, part.ProgChgPortIn...)
		ports = append(ports, part.MutePortIn...)
	}
	return
}

func OutputPorts(trk *Track) (ports []string) {
	parts, err := Parts(trk)
	if err != nil {
		return
	}
	for _, part := range parts {
		ports = append(ports, part.ProgChgPortOut...)
		ports = append(ports, part.MutePortOut...)
	}
	return
}

func DrawString(x, y int, s string, fg, bg termbox.Attribute) {
	for _, c := range s {
		switch c {
		case '\n':
			x = 0
			y++
		default:
			termbox.SetCell(x, y, c, fg, bg)
			x++
		}
	}
}
