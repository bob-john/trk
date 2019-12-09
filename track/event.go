package track

type Event struct {
	ID    string
	Port  string
	Tick  int
	Bytes []byte
}

// func (trk *Track) Events(tick int) (events []*Event) {
// 	for _, e := range trk.events {
// 		if e.Tick == tick {
// 			events = append(events, e)
// 		} else if e.Tick > tick {
// 			return
// 		}
// 	}
// 	return
// }

func (trk *Track) Events() (events []*Event) {
	return trk.events
}

func (e *Event) Save(trk *Track) (err error) {
	if e.ID == "" {
		e.ID = makeID()
	}
	err = trk.db.Save(e)
	if err != nil {
		return
	}
	return trk.db.All(&trk.events)
}

// func sortEventSlice(sl []*Event) {
// 	sort.SliceStable(sl, func(i, j int) bool {
// 		if sl[i].Tick == sl[j].Tick {
// 			return sl[i].Port < sl[j].Port
// 		}
// 		return sl[i].Tick < sl[j].Tick
// 	})
// }
