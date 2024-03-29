package rtmididrv

import "trk/rtmidi"

type inPort struct {
	port   int
	name   string
	input  rtmidi.MIDIIn
	opened bool
}

func (i *inPort) Open() (err error) {
	if i.opened {
		return nil
	}
	// i.input, err = rtmidi.NewMIDIInDefault()
	i.input, err = rtmidi.NewMIDIIn(rtmidi.APIUnspecified, "", 32*1024)
	if err != nil {
		return
	}
	defer func() {
		if err != nil {
			i.input.Destroy()
		}
	}()
	// err = i.input.OpenPort(i.port, i.name)
	err = i.input.OpenPort(i.port, "")
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			i.input.Close()
		}
	}()
	err = i.input.IgnoreTypes(true, false, true)
	if err != nil {
		return
	}
	i.opened = true
	return
}

func (i *inPort) Close() (err error) {
	if !i.opened {
		return nil
	}
	err = i.input.Close()
	if err != nil {
		return
	}
	i.opened = false
	return
}

func (i *inPort) IsOpen() bool {
	return i.opened
}

func (i *inPort) Number() int {
	return i.port
}

func (i *inPort) String() string {
	return i.name
}

func (i *inPort) Underlying() interface{} {
	return i.input
}

func (i *inPort) SetListener(ln func([]byte, int64)) error {
	return i.input.SetCallback(func(in rtmidi.MIDIIn, data []byte, deltaMilliseconds float64) {
		ln(data, int64(deltaMilliseconds*1000))
	})
}

func (i *inPort) StopListening() error {
	return i.input.CancelCallback()
}
