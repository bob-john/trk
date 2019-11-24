package main

import (
	"archive/zip"
	"os"
)

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
