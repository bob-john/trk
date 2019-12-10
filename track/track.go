package track

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"os"
)

var (
	ErrBadType  = errors.New("track: bad type")
	ErrBadRoute = errors.New("track: bad route")
)

type Track struct {
	Devices []*Device `json:",omitempty"`
	Routes  []*Route  `json:",omitempty"`
	Events  []*Event  `json:",omitempty"`
	file    *os.File
}

func Open(name string) (trk *Track, err error) {
	f, err := os.OpenFile(name, os.O_RDWR|os.O_CREATE|os.O_SYNC, 0666)
	if err != nil {
		return
	}
	defer func() {
		if err != nil {
			f.Close()
		}
	}()

	trk = &Track{file: f}

	st, err := f.Stat()
	if err != nil && st.Size() == 0 {
		d := json.NewDecoder(f)
		err = d.Decode(trk)
		if err != nil {
			return nil, err
		}
	}

	_, err = trk.CreateDeviceIfNotExist("DIGITAKT")
	if err != nil {
		return nil, err
	}
	_, err = trk.CreateDeviceIfNotExist("DIGITONE")
	if err != nil {
		return nil, err
	}
	err = trk.Save()
	if err != nil {
		return nil, err
	}

	return
}

func (trk *Track) Close() error {
	return trk.file.Close()
}

func (trk *Track) CreateDeviceIfNotExist(name string) (*Device, error) {
	for _, dev := range trk.Devices {
		if dev.Name == name {
			return dev, nil
		}
	}
	dev := &Device{Name: name}
	trk.Devices = append(trk.Devices, dev)
	err := trk.CreateMissingRoutes()
	if err != nil {
		return dev, err
	}
	return dev, trk.Save()
}

func (trk *Track) CreateRouteIfNotExist(input, output string) (*Route, error) {
	if input == output {
		return nil, ErrBadRoute
	}
	for _, r := range trk.Routes {
		if r.Input == input && r.Output == output {
			return r, nil
		}
	}
	r := &Route{Input: input, Output: output}
	trk.Routes = append(trk.Routes, r)
	err := trk.Save()
	if err != nil {
		return nil, err
	}
	return r, nil
}

func (trk *Track) CreateMissingRoutes() (err error) {
	for _, i := range trk.Devices {
		for _, o := range trk.Devices {
			if i.Name == o.Name {
				continue
			}
			_, err := trk.CreateRouteIfNotExist(i.Name, o.Name)
			if err != nil {
				return err
			}
		}
	}
	return
}

func (trk *Track) Insert(obj interface{}) error {
	switch obj := obj.(type) {
	case *Event:
		trk.Events = append(trk.Events, obj)

	default:
		return ErrBadType
	}
	return trk.Save()
}

func (trk *Track) Save() (err error) {
	data, err := json.Marshal(trk)
	if err != nil {
		return
	}
	_, err = trk.file.Seek(0, os.SEEK_SET)
	if err != nil {
		return
	}
	_, err = io.Copy(trk.file, bytes.NewReader(data))
	if err != nil {
		return
	}
	return trk.file.Truncate(int64(len(data)))
}

func (trk *Track) InputPorts() (ports []string) {
	for _, d := range trk.Devices {
		ports = append(ports, d.Input)
	}
	return
}

func (trk *Track) OutputPorts() (ports []string) {
	for _, d := range trk.Devices {
		ports = append(ports, d.Output)
	}
	return
}
