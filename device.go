package main

import (
	"io"

	"github.com/gomidi/midi"
	"github.com/gomidi/midi/midimessage/realtime"
	"github.com/gomidi/midi/midireader"
	"gitlab.com/gomidi/midi/mid"
)

type Device struct {
	*Input
	*Output
}

func ConnectDevice(inputName, outputName string) (*Device, error) {
	in, err := OpenInput(inputName)
	if err != nil {
		return nil, err
	}
	out, err := OpenOutput(outputName)
	if err != nil {
		return nil, err
	}
	return &Device{in, out}, nil
}

func (d *Device) Close() {
	d.Input.Close()
	d.out.Close()
}

type Input struct {
	in  mid.In
	inC chan midi.Message
}

func OpenInput(name string) (*Input, error) {
	in, err := mid.OpenIn(midiDriver{}, -1, name)
	if err != nil {
		return nil, err
	}
	inC := make(chan midi.Message)
	r, w := io.Pipe()
	datac := make(chan []byte)
	go func() {
		for data := range datac {
			w.Write(data)
		}
	}()
	err = in.SetListener(func(data []byte, deltaMicroseconds int64) {
		datac <- data
	})
	if err != nil {
		in.Close()
		return nil, err
	}
	go func() {
		rd := midireader.New(r, func(m realtime.Message) {
			inC <- m
		})
		for {
			m, err := rd.Read()
			if err != nil {
				return
			}
			inC <- m
		}
	}()
	return &Input{in, inC}, nil
}

func (i *Input) Close() error {
	close(i.inC)
	return i.in.Close()
}

func (i *Input) In() <-chan midi.Message {
	return i.inC
}

type Output struct {
	out mid.Out
	w   *mid.Writer
}

func OpenOutput(name string) (*Output, error) {
	out, err := mid.OpenOut(midiDriver{}, -1, name)
	if err != nil {
		return nil, err
	}
	w := mid.ConnectOut(out)
	return &Output{out, w}, nil
}

func (o *Output) Close() error {
	return o.out.Close()
}

func (o *Output) Write(m midi.Message) error {
	return o.w.Write(m)
}