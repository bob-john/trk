package track

func Parts(trk *Track) (parts []*Part, err error) {
	return trk.parts, nil
}

func Pattern(trk *Track, part *Part, tick int) (pattern int) {
	for _, pc := range trk.pc {
		if pc.Tick <= tick {
			if pc.Part == part.Name {
				pattern = pc.Pattern
			}
		} else {
			break
		}
	}
	return
}

func IsPatternModified(trk *Track, part *Part, tick int) bool {
	for _, pc := range trk.pc {
		if pc.Tick == tick && pc.Part == part.Name {
			return true
		}
	}
	return false
}

func Mute(trk *Track, part *Part, tick int) (mute [16]bool) {
	for _, mc := range trk.mc {
		if mc.Tick <= tick {
			if mc.Part == part.Name {
				mute = mc.Mute
			}
		} else {
			break
		}
	}
	return
}

func IsMuteModified(trk *Track, part *Part, tick int) bool {
	for _, mc := range trk.mc {
		if mc.Tick == tick && mc.Part == part.Name {
			return true
		}
	}
	return false
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
