package main

import (
	"gitlab.com/gomidi/midi/mid"

	// #cgo LDFLAGS: -lwinmm
	// #include <windows.h>
	"C"
)
import "fmt"

type midiDriver struct{}

func (d midiDriver) Ins() ([]mid.In, error) {
	var ins []mid.In
	n := int(C.midiInGetNumDevs())
	for i := 0; i < n; i++ {
		var caps C.MIDIINCAPS
		err := C.midiInGetDevCaps(C.ulonglong(i), &caps, C.sizeof_MIDIINCAPS)
		if err != C.MMSYSERR_NOERROR {
			return nil, fmt.Errorf("mm: %d", err)
		}
		ins = append(ins, midiIn{i, C.GoString(&caps.szPname[0])})
	}
	return ins, nil
}

func (d midiDriver) Outs() ([]mid.Out, error) {
	var outs []mid.Out
	n := int(C.midiOutGetNumDevs())
	for i := 0; i < n; i++ {
		var caps C.MIDIOUTCAPS
		err := C.midiOutGetDevCaps(C.ulonglong(i), &caps, C.sizeof_MIDIOUTCAPS)
		if err != C.MMSYSERR_NOERROR {
			return nil, fmt.Errorf("mm: %d", err)
		}
		outs = append(outs, midiOut{i, C.GoString(&caps.szPname[0])})
	}
	return outs, nil
}

func (d midiDriver) String() string {
	return "WinMM MIDI driver"
}

func (d midiDriver) Close() error {
	return nil
}

type midiIn struct {
	deviceID   int
	deviceName string
}

func (d midiIn) Open() error {
	return nil
}

func (d midiIn) Close() error {
	return nil
}

func (d midiIn) IsOpen() bool {
	return false
}

func (d midiIn) Number() int {
	return d.deviceID
}

func (d midiIn) String() string {
	return d.deviceName
}

func (d midiIn) Underlying() interface{} {
	return nil
}

func (d midiIn) SetListener(func(data []byte, deltaMicroseconds int64)) error {
	return nil
}

func (d midiIn) StopListening() error {
	return nil
}

type midiOut struct {
	deviceID   int
	deviceName string
}

func (d midiOut) Open() error {
	return nil
}

func (d midiOut) Close() error {
	return nil
}

func (d midiOut) IsOpen() bool {
	return false
}

func (d midiOut) Number() int {
	return d.deviceID
}

func (d midiOut) String() string {
	return d.deviceName
}

func (d midiOut) Underlying() interface{} {
	return nil
}

func (d midiOut) Send([]byte) error {
	return nil
}
