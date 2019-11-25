package main

import (
	"io"

	"github.com/gomidi/midi"
	"github.com/gomidi/midi/midimessage/realtime"
	"github.com/gomidi/midi/midireader"
	"gitlab.com/gomidi/midi/mid"
)

var driver, _ = NewDriver()

type Device struct {
	name string
	*Input
	*Output
}

func OpenDevice(name, inputName, outputName string) (*Device, error) {
	in, err := OpenInput(inputName)
	if err != nil {
		return nil, err
	}
	out, err := OpenOutput(outputName)
	if err != nil {
		return nil, err
	}
	return &Device{name, in, out}, nil
}

func (d *Device) Name() string {
	return d.name
}

func (d *Device) Close() {
	d.Input.Close()
	d.out.Close()
}

type Input struct {
	name string
	in  mid.In
	inC chan midi.Message
}

func OpenInput(name string) (*Input, error) {
	in, err := mid.OpenIn(driver, -1, name)
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
	return &Input{name, in, inC}, nil
}

func (i *Input) Close() error {
	err := i.in.Close()
	if err != nil {
		return err
	}
	close(i.inC)
	return nil
}

func (i *Input) In() <-chan midi.Message {
	return i.inC
}

func (i *Input) Name() string {
	return i.name
}

type Output struct {
	out mid.Out
	w   *mid.Writer
}

func OpenOutput(name string) (*Output, error) {
	out, err := mid.OpenOut(driver, -1, name)
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

func (o *Output) String() string {
	return o.out.String()
}

type Message struct {
	midi.Message
	Port string
}

// type Ports struct {
// 	input  map[string]*Input
// 	output map[string]*Output
// }

// func NewPorts() *Ports {
// 	return &Ports{make(map[string]*Input), make(map[string]*Output)}
// }

// func (p *Ports) SetInputs(names ...string) {
// 	required := make(map[string]bool)
// 	for _, name := range names {
// 		port, err := OpenInput(name)
// 		if err == nil {
// 			p.input[name] = port
// 			required[name] = true
// 		}
// 	}
// 	for name, port := range p.input {
// 		if !required[name] {
// 			port.Close()
// 			delete(p.input, name)
// 		}
// 	}
// }

// func (p *Ports) SetOutputs(names ...string) {
// 	required := make(map[string]bool)
// 	for _, name := range names {
// 		port, err := OpenOutput(name)
// 		if err == nil {
// 			p.output[name] = port
// 			required[name] = true
// 		}
// 	}
// 	for name, port := range p.output {
// 		if !required[name] {
// 			port.Close()
// 			delete(p.output, name)
// 		}
// 	}
// }

// func (p *Ports) Input(name string) *Input {
// 	port, _ := p.input[name]
// 	return port
// }

// func (p *Ports) Output(name string) *Output {
// 	port, _ := p.output[name]
// 	return port
// }
