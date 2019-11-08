package main

import "gitlab.com/gomidi/midi/mid"

type midiDriver struct{}

func (d midiDriver) Ins() ([]mid.In, error) {
	return nil, nil
}

func (d midiDriver) Outs() ([]mid.Out, error) {
	return nil, nil
}

func (d midiDriver) String() string {
	return "Stub Linux MIDI driver"
}

func (d midiDriver) Close() error {
	return nil
}
