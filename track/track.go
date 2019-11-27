package track

import (
	"time"

	"github.com/asdine/storm"
	"go.etcd.io/bbolt"
)

type Track struct {
	db *storm.DB
}

func Open(name string) (*Track, error) {
	db, err := storm.Open(name, storm.BoltOptions(0600, &bbolt.Options{Timeout: 1 * time.Second}))
	if err != nil {
		return nil, err
	}
	trk := &Track{db}
	for _, part := range []*Part{newPart("DIGITAKT", "DT", 16), newPart("DIGITONE", "DN", 8)} {
		err = CreateIfNotExists(trk, part)
		if err != nil {
			return nil, err
		}
		err = SetPattern(trk, part, 0, 0)
		if err != nil {
			return nil, err
		}
		err = SetMute(trk, part, 0, [16]bool{})
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
	return &Part{name, shortName, make([]int, trackCount), nil, nil, nil, nil, 10, 10}
}

type PatternChange struct {
	ID      string
	Part    string `storm:"index"`
	Tick    int    `storm:"index"`
	Pattern int
}

type MuteChange struct {
	ID   string
	Part string `storm:"index"`
	Tick int    `storm:"index"`
	Mute [16]bool
}
