package track

import "github.com/asdine/storm"

func Parts(trk *Track) (parts []*Part, err error) {
	err = trk.db.All(&parts)
	return
}

func Pattern(trk *Track, part *Part, tick int) (int, bool) {
	var chg []*PatternChange
	trk.db.Find("Part", part.Name, &chg, storm.Reverse(), storm.Limit(2))
	switch len(chg) {
	case 1:
		return chg[0].Pattern, chg[0].Tick == tick
	case 2:
		return chg[1].Pattern, chg[1].Tick == tick && chg[1].Pattern != chg[0].Pattern
	}
	return 0, false
}

func Mute(trk *Track, part *Part, tick int) ([16]bool, bool) {
	var chg []*MuteChange
	trk.db.Find("Part", part.Name, &chg, storm.Reverse(), storm.Limit(2))
	switch len(chg) {
	case 1:
		return chg[0].Mute, chg[0].Tick == tick
	case 2:
		return chg[1].Mute, chg[1].Tick == tick && chg[1].Mute != chg[0].Mute
	}
	return [16]bool{}, false
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
