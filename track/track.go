package track

import (
	"sort"
	"time"

	"github.com/asdine/storm"
	"go.etcd.io/bbolt"
)

type Track struct {
	db    *storm.DB
	parts []*Part
	pc    []*PatternChange
	mc    []*MuteChange
}

func Open(name string) (*Track, error) {
	db, err := storm.Open(name, storm.BoltOptions(0600, &bbolt.Options{Timeout: 1 * time.Second}))
	if err != nil {
		return nil, err
	}
	trk := &Track{db, nil, nil, nil}
	err = db.All(&trk.parts)
	if err != nil {
		return nil, err
	}
	sortPartSlice(trk.parts)
	var pcs []*PatternChange
	err = db.All(&pcs)
	if err != nil {
		return nil, err
	}
	for _, pc := range pcs {
		trk.pc = append(trk.pc, pc)
	}
	sortPatternChangeSlice(trk.pc)
	var mcs []*MuteChange
	err = db.All(&mcs)
	if err != nil {
		return nil, err
	}
	for _, mc := range mcs {
		trk.mc = append(trk.mc, mc)
	}
	sortMuteChangeSlice(trk.mc)
	for _, part := range []*Part{newPart("DIGITAKT", "DT", 16), newPart("DIGITONE", "DN", 8)} {
		err = trk.CreateIfNotExists(part)
		if err != nil {
			return nil, err
		}
	}
	return trk, nil
}

func (trk *Track) Close() error {
	return trk.db.Close()
}

type Part struct {
	Name           string `storm:"id"`
	ShortName      string
	TrackCh        []int
	ProgChgPortIn  []string
	ProgChgPortOut []string
	MutePortIn     []string
	MutePortOut    []string
	ProgChgInCh    int
	ProgChgOutCh   int
}

func newPart(name, shortName string, trackCount int) *Part {
	return &Part{name, shortName, make([]int, trackCount), nil, nil, nil, nil, 9, 9}
}

func (p *Part) TrackOf(ch int) int {
	for n, c := range p.TrackCh {
		if c == ch {
			return n
		}
	}
	return -1
}

func sortPartSlice(sl []*Part) {
	sort.SliceStable(sl, func(i, j int) bool {
		return sl[i].Name < sl[j].Name
	})
}

type PatternChange struct {
	ID      string
	Part    string `storm:"index"`
	Tick    int    `storm:"index"`
	Pattern int
}

func sortPatternChangeSlice(sl []*PatternChange) {
	sort.SliceStable(sl, func(i, j int) bool {
		if sl[i].Tick == sl[j].Tick {
			return sl[i].Part < sl[j].Part
		}
		return sl[i].Tick < sl[j].Tick
	})
}

type MuteChange struct {
	ID   string
	Part string `storm:"index"`
	Tick int    `storm:"index"`
	Mute [16]bool
}

func sortMuteChangeSlice(sl []*MuteChange) {
	sort.SliceStable(sl, func(i, j int) bool {
		if sl[i].Tick == sl[j].Tick {
			return sl[i].Part < sl[j].Part
		}
		return sl[i].Tick < sl[j].Tick
	})
}
