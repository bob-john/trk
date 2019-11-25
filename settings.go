package main

import (
	"encoding/json"
	"io"
)

const (
	Digitakt = "Digitakt"
	Digitone = "Digitone"
)

type Settings struct {
	Devices map[string]*DeviceSettings
}

func NewSettings() *Settings {
	return &Settings{
		Devices: map[string]*DeviceSettings{
			Digitakt: NewDeviceSettings(DeviceSourceDigitakt, 8),
			Digitone: NewDeviceSettings(DeviceSourceDigitone, 4),
		},
	}
}

func ReadSettings(f io.Reader) (*Settings, error) {
	s := NewSettings()
	err := json.NewDecoder(f).Decode(s)
	return s, err
}

func (s *Settings) Write(f io.Writer) error {
	return json.NewEncoder(f).Encode(s)
}

func (s *Settings) InputPortNames() (names []string) {
	for _, device := range s.Devices {
		for name := range device.Inputs {
			names = append(names, name)
		}
	}
	return
}

type DeviceSettings struct {
	Inputs       map[string]struct{}
	Outputs      map[string]struct{}
	Channels     []int
	ProgChgSrc   DeviceSource
	MuteSrc      DeviceSource
	ProgChgInCh  int
	ProgChgOutCh int
}

func NewDeviceSettings(source DeviceSource, trackCount int) *DeviceSettings {
	return &DeviceSettings{
		Inputs:       make(map[string]struct{}),
		Outputs:      make(map[string]struct{}),
		ProgChgSrc:   source,
		ProgChgInCh:  10,
		ProgChgOutCh: 10,
		MuteSrc:      source,
		Channels:     make([]int, trackCount),
	}
}

type DeviceSource int

const (
	DeviceSourceDigitakt DeviceSource = iota
	DeviceSourceDigitone
	DeviceSourceBoth
)

func (s DeviceSource) String() string {
	switch s {
	case DeviceSourceDigitakt:
		return "Digitakt"
	case DeviceSourceDigitone:
		return "Digitone"
	case DeviceSourceBoth:
		return "Both"
	}
	return ""
}
