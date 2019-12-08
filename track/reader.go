package track

func (trk *Track) Parts() []*Part {
	return trk.parts
}

func (trk *Track) Pattern(part *Part, tick int) (pattern int) {
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

func (trk *Track) IsPatternModified(part *Part, tick int) bool {
	for _, pc := range trk.pc {
		if pc.Tick == tick && pc.Part == part.Name {
			return true
		}
	}
	return false
}

func (trk *Track) Mute(part *Part, tick int) (mute [16]bool) {
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

func (trk *Track) IsMuteModified(part *Part, tick int) bool {
	for _, mc := range trk.mc {
		if mc.Tick == tick && mc.Part == part.Name {
			return true
		}
	}
	return false
}

func (trk *Track) Events(tick int) (events []*Event) {
	for _, e := range trk.events {
		if e.Tick == tick {
			events = append(events, e)
		} else if e.Tick > tick {
			return
		}
	}
	return
}

func (trk *Track) IsPartModified(part *Part, tick int) bool {
	return trk.IsPatternModified(part, tick) || trk.IsMuteModified(part, tick)
}

func (trk *Track) IsModified(tick int) bool {
	for _, e := range trk.events {
		if e.Tick == tick {
			return true
		} else if e.Tick > tick {
			break
		}
	}
	for _, part := range trk.parts {
		if trk.IsPatternModified(part, tick) {
			return true
		}
		if trk.IsMuteModified(part, tick) {
			return true
		}
	}
	return false
}

func (trk *Track) Filters() []*Filter {
	return trk.filters
}

func (trk *Track) InputPorts() (ports []string) {
	for _, part := range trk.parts {
		ports = append(ports, part.ProgChgPortIn...)
		ports = append(ports, part.MutePortIn...)
	}
	return
}

func (trk *Track) OutputPorts() (ports []string) {
	for _, part := range trk.parts {
		ports = append(ports, part.ProgChgPortOut...)
		ports = append(ports, part.MutePortOut...)
	}
	return
}
