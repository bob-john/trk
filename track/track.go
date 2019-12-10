package track

import "os"

import "encoding/json"

type Track struct {
	Devices []*Device
	Routes  []*Route
	Events  []*Event
	file    *os.File
}

func Open(name string) (trk *Track, err error) {
	f, err := os.Open(name)
	if os.IsNotExist(err) {
		f, err = os.Create(name)
		if err != nil {
			return
		}
		_, err = f.WriteString("{}")
		if err != nil {
			return
		}
		trk = new(Track)
	} else if err == nil {
		defer func() {
			if err != nil {
				f.Close()
			}
		}()
		d := json.NewDecoder(f)
		trk = new(Track)
		err = d.Decode(trk)
		if err != nil {
			return
		}
	} else {
		return
	}
	trk.file = f
	return
}

func (trk *Track) Close() error {
	return trk.file.Close()
}

func (trk *Track) Insert(tick int, port string, message []byte) error {
	trk.Events = append(trk.Events, &Event{Tick: tick, Port: port, Message: message})
	return trk.Save()
}

func (trk *Track) Save() (err error) {
	err = trk.file.Truncate(0)
	if err != nil {
		return
	}
	e := json.NewEncoder(trk.file)
	err = e.Encode(trk)
	if err != nil {
		return
	}
	return trk.file.Sync()
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
