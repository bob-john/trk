package main

import (
	"gitlab.com/gomidi/midi/mid"

	// #cgo LDFLAGS: -lwinmm
	// #include <windows.h>
	"C"
)
import (
	"fmt"
	"log"
	"syscall"
)

var (
	midiInListeners = make(map[int]func(data []byte, deltaMicroseconds int64))
)

// https://docs.microsoft.com/fr-fr/windows/win32/multimedia/midi-functions
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
		ins = append(ins, &midiIn{deviceID: i, deviceName: C.GoString(&caps.szPname[0])})
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
		outs = append(outs, &midiOut{deviceID: i, deviceName: C.GoString(&caps.szPname[0])})
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
	handle     C.HMIDIIN
}

func (d *midiIn) Open() error {
	err := C.midiInOpen(&d.handle, C.UINT(d.deviceID), C.DWORD_PTR(syscall.NewCallback(midiInProc)), C.DWORD_PTR(d.deviceID), C.CALLBACK_FUNCTION)
	if err != C.MMSYSERR_NOERROR {
		return fmt.Errorf("mm: %d", err)
	}
	return nil
}

func (d *midiIn) Close() error {
	err := C.midiInClose(d.handle)
	if err != C.MMSYSERR_NOERROR {
		return fmt.Errorf("mm: %d", err)
	}
	d.handle = nil
	return nil
}

func (d *midiIn) IsOpen() bool {
	return d.handle != nil
}

func (d *midiIn) Number() int {
	return d.deviceID
}

func (d *midiIn) String() string {
	return d.deviceName
}

func (d *midiIn) Underlying() interface{} {
	return nil
}

func (d *midiIn) SetListener(ls func(data []byte, deltaMicroseconds int64)) error {
	midiInListeners[d.deviceID] = ls
	err := C.midiInStart(d.handle)
	if err != C.MMSYSERR_NOERROR {
		return fmt.Errorf("mm: %d", err)
	}
	return nil
}

func (d *midiIn) StopListening() error {
	err := C.midiInStop(d.handle)
	if err != C.MMSYSERR_NOERROR {
		return fmt.Errorf("mm: %d", err)
	}
	return nil
}

func midiInProc(hMidiIn C.HMIDIIN, wMsg C.UINT, dwInstance C.DWORD_PTR, dwParam1 C.DWORD_PTR, dwParam2 C.DWORD_PTR) uintptr {
	switch wMsg {
	case C.MIM_OPEN:
		// log.Println(dwInstance, "MIM_OPEN")

	case C.MIM_CLOSE:
		// log.Println(dwInstance, "MIM_CLOSE")

	case C.MIM_DATA:
		// log.Println(dwInstance, "MIM_DATA", dwParam1, dwParam2)
		ls, ok := midiInListeners[int(dwInstance)]
		if !ok {
			return 0
		}
		b := []byte{byte(dwParam1), byte(dwParam1 >> 8), byte(dwParam1 >> 16)}
		switch b[0] & 0xF0 {
		case 0xF0:
			switch b[0] {
			case 0xF1, 0xF3:
				b = b[:2]
			case 0xF2:
				b = b[:3]
			default:
				b = b[:1]
			}
		case 0xC0, 0xD0:
			b = b[:2]
		}
		if b[0] == 0xF8 {
			return 0
		}
		ls(b, int64(dwParam2)*1000)
	}
	return 0
}

type midiOut struct {
	deviceID   int
	deviceName string
	handle     C.HMIDIOUT
}

func (d *midiOut) Open() error {
	err := C.midiOutOpen(&d.handle, C.UINT(d.deviceID), C.DWORD_PTR(0), C.DWORD_PTR(0), C.CALLBACK_NULL)
	if err != C.MMSYSERR_NOERROR {
		return fmt.Errorf("mm: %d", err)
	}
	return nil
}

func (d *midiOut) Close() error {
	err := C.midiOutClose(d.handle)
	if err != C.MMSYSERR_NOERROR {
		return fmt.Errorf("mm: %d", err)
	}
	d.handle = nil
	return nil
}

func (d *midiOut) IsOpen() bool {
	return d.handle != nil
}

func (d *midiOut) Number() int {
	return d.deviceID
}

func (d *midiOut) String() string {
	return d.deviceName
}

func (d *midiOut) Underlying() interface{} {
	return nil
}

func (d *midiOut) Send(b []byte) error {
	if len(b) == 3 {
		err := C.midiOutShortMsg(d.handle, C.DWORD(int32(b[0])|(int32(b[1])<<8)|(int32(b[2])<<16)))
		if err != C.MMSYSERR_NOERROR {
			return fmt.Errorf("mm: %d", err)
		}
	} else {
		log.Printf("cannot send %v, len=%d", b, len(b))
	}
	return nil
}
