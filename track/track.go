package track

import (
	"time"

	"github.com/asdine/storm"
	"go.etcd.io/bbolt"
)

type Track struct {
	db     *storm.DB
	inputs []*Input
	events []*Event
}

func Open(name string) (trk *Track, err error) {
	db, err := storm.Open(name, storm.BoltOptions(0600, &bbolt.Options{Timeout: 1 * time.Second}))
	if err != nil {
		return nil, err
	}
	trk = &Track{db: db}
	err = trk.db.All(&trk.inputs)
	if err != nil {
		return
	}
	err = trk.db.All(&trk.events)
	if err != nil {
		return
	}
	return
}

func (trk *Track) Close() error {
	return trk.db.Close()
}
