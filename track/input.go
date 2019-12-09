package track

type Input struct {
	Port   string `storm:"id"`
	Filter []Filter
}

func (trk *Track) SetInput(port string, enabled bool) (err error) {
	if enabled && trk.Input(port) == nil {
		err = trk.db.Save(&Input{Port: port})
		if err != nil {
			return
		}
	} else if !enabled {
		err = trk.db.DeleteStruct(&Input{Port: port})
		if err != nil {
			return
		}
	}
	return trk.db.All(&trk.inputs)
}

func (trk *Track) Input(port string) *Input {
	for _, i := range trk.inputs {
		if i.Port == port {
			return i
		}
	}
	return nil
}

func (trk *Track) InputPorts() (ports []string) {
	for _, i := range trk.inputs {
		ports = append(ports, i.Port)
	}
	return
}
