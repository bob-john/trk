package main

import (
	"archive/zip"
	"crypto/rand"
	"fmt"
	"os"
	"time"

	"github.com/asdine/storm"
	"github.com/btcsuite/btcutil/base58"
	"go.etcd.io/bbolt"
)

func OpenTrack(name string) (*storm.DB, error) {
	trk, err := storm.Open(name, storm.BoltOptions(0600, &bbolt.Options{Timeout: 1 * time.Second}))
	if err != nil {
		return nil, err
	}
	n, err := trk.Count(&Part1{})
	if err != nil {
		return nil, err
	}
	if n == 0 {
		var (
			dt = NewPart1("DIGITAKT", "DT", 16)
			dn = NewPart1("DIGITONE", "DN", 8)
		)
		err = trk.Save(dt)
		if err != nil {
			return nil, err
		}
		err = trk.Save(dn)
		if err != nil {
			return nil, err
		}
		err = trk.Save(dt.NewPatternChange(0, 0))
		if err != nil {
			return nil, err
		}
		err = trk.Save(dn.NewPatternChange(0, 0))
		if err != nil {
			return nil, err
		}
		err = trk.Save(dt.NewMuteChange(0, [16]bool{}))
		if err != nil {
			return nil, err
		}
		err = trk.Save(dn.NewMuteChange(0, [16]bool{}))
		if err != nil {
			return nil, err
		}
	}
	return trk, nil
}

type Part1 struct {
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

func NewPart1(name, shortName string, trackCount int) *Part1 {
	return &Part1{name, shortName, make([]int, trackCount), nil, nil, nil, nil, 10, 10}
}

func (p *Part1) Pattern(trk *storm.DB, tick int) (int, bool) {
	var chg []*PatternChange
	trk.Find("Part", p.Name, &chg, storm.Reverse(), storm.Limit(2))
	switch len(chg) {
	case 1:
		return chg[0].Pattern, chg[0].Tick == tick
	case 2:
		return chg[1].Pattern, chg[1].Tick == tick && chg[1].Pattern != chg[0].Pattern
	}
	return 0, false
}

func (p *Part1) Mute(trk *storm.DB, tick int) ([16]bool, bool) {
	var chg []*MuteChange
	trk.Find("Part", p.Name, &chg, storm.Reverse(), storm.Limit(2))
	switch len(chg) {
	case 1:
		return chg[0].Mute, chg[0].Tick == tick
	case 2:
		return chg[1].Mute, chg[1].Tick == tick && chg[1].Mute != chg[0].Mute
	}
	return [16]bool{}, false
}

func (p *Part1) NewPatternChange(tick, pattern int) *PatternChange {
	return &PatternChange{ID(), p.Name, tick, pattern}
}

func (p *Part1) NewMuteChange(tick int, mute [16]bool) *MuteChange {
	return &MuteChange{ID(), p.Name, tick, mute}
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

func FormatMute(mute [16]bool, part *Part1) (str string) {
	for n, ch := range part.TrackCh {
		if ch == 0 {
			continue
		}
		if mute[n] {
			str += "-"
		} else if n < 8 {
			str += string('1' + n)
		} else {
			str += string('A' + n - 8)
		}
		//TODO Handle Digitone (1234 1234)
		//TODO Improve Digitakt (12345678 ABCDEFGH)
	}
	return
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

func ID() string {
	b := make([]byte, 12)
	_, err := rand.Read(b)
	must(err)
	return base58.Encode(b)
}
