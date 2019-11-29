package rtmididrv

import "gitlab.com/gomidi/midi/mid"

import "trk/rtmidi"

func New() (mid.Driver, error) {
	return &driver{}, nil
}

type driver struct{}

func (d *driver) Ins() ([]mid.In, error) {
	in, err := rtmidi.NewMIDIInDefault()
	if err != nil {
		return nil, err
	}
	defer in.Destroy()
	count, err := in.PortCount()
	if err != nil {
		return nil, err
	}
	var res []mid.In
	for port := 0; port < count; port++ {
		name, err := in.PortName(port)
		if err != nil {
			return nil, err
		}
		res = append(res, &inPort{port, name, nil, false})
	}
	return res, nil
}

func (d *driver) Outs() ([]mid.Out, error) {
	out, err := rtmidi.NewMIDIOutDefault()
	if err != nil {
		return nil, err
	}
	defer out.Destroy()
	count, err := out.PortCount()
	if err != nil {
		return nil, err
	}
	var res []mid.Out
	for port := 0; port < count; port++ {
		name, err := out.PortName(port)
		if err != nil {
			return nil, err
		}
		res = append(res, &outPort{port, name, nil, false})
	}
	return res, nil
}

func (d *driver) String() string {
	return "rtmidi"
}

func (d *driver) Close() error {
	return nil
}
