package main

import (
	"archive/zip"
	"os"
	"time"

	"github.com/asdine/storm"
	"go.etcd.io/bbolt"
)

func OpenTrack(name string) (*storm.DB, error) {
	trk, err := storm.Open(name, storm.BoltOptions(0600, &bbolt.Options{Timeout: 1 * time.Second}))
	if err != nil {
		return nil, err
	}
	err = trk.Save(NewPart1("DIGITAKT", "DT", 8))
	if err != nil {
		return nil, err
	}
	err = trk.Save(NewPart1("DIGITONE", "DB", 4))
	if err != nil {
		return nil, err
	}
	return trk, nil
}

type Part1 struct {
	Name          string `storm:"id"`
	ShortName     string
	Track         []int
	ProgChgPortIn []string
	MutePortIn    []string
	PortOut       []string
	ProgChgInCh   int
	ProgChgOutCh  int
}

func NewPart1(name, shortName string, trackCount int) *Part1 {
	return &Part1{name, shortName, make([]int, trackCount), nil, nil, nil, 10, 10}
}

type PatternChange struct {
	Tick    int    `storm:"id"`
	Part    string `storm:"index"`
	Pattern int
}

type MuteChange struct {
	Tick int    `storm:"id"`
	Part string `storm:"index"`
	Mute [16]bool
}

type Track struct {
	Seq      *Seq
	Settings *Settings
}

func NewTrack() *Track {
	return &Track{NewSeq(), NewSettings()}
}

func ReadTrack(path string) (*Track, error) {
	track := new(Track)
	r, err := zip.OpenReader(path)
	if err != nil {
		return nil, err
	}
	defer r.Close()
	for _, f := range r.File {
		w, err := f.Open()
		if err != nil {
			return nil, err
		}
		defer w.Close()
		switch f.Name {
		case "seq.csv":
			track.Seq, err = ReadSeq(w)
			if err != nil {
				return nil, err
			}
		case "settings.json":
			track.Settings, err = ReadSettings(w)
			if err != nil {
				return nil, err
			}
		}
		w.Close()
	}
	return track, nil
}

func (t *Track) Write(path string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	w := zip.NewWriter(f)
	defer w.Close()
	for _, name := range []string{"seq.csv", "settings.json"} {
		f, err := w.Create(name)
		if err != nil {
			return err
		}
		switch name {
		case "seq.csv":
			err = t.Seq.Write(f)
			if err != nil {
				return err
			}
		case "settings.json":
			err = t.Settings.Write(f)
			if err != nil {
				return err
			}
		}
	}
	err = w.Close()
	if err != nil {
		return err
	}
	return f.Close()
}
