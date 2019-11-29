package rtmididrv

import "trk/rtmidi"

type outPort struct {
	port   int
	name   string
	output rtmidi.MIDIOut
	opened bool
}

func (o *outPort) Open() (err error) {
	if o.opened {
		return nil
	}
	o.output, err = rtmidi.NewMIDIOutDefault()
	if err != nil {
		return
	}
	defer func() {
		if err != nil {
			o.output.Destroy()
		}
	}()
	err = o.output.OpenPort(o.port, o.name)
	if err != nil {
		return err
	}
	o.opened = true
	return
}

func (o *outPort) Close() error {
	if !o.opened {
		return nil
	}
	defer o.output.Destroy()
	return o.output.Close()
}

func (o *outPort) IsOpen() bool {
	return o.opened
}

func (o *outPort) Number() int {
	return o.port
}

func (o *outPort) String() string {
	return o.name
}

func (o *outPort) Underlying() interface{} {
	return o.output
}

func (o *outPort) Send(data []byte) error {
	return o.output.SendMessage(data)
}
