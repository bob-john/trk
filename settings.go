package main

import (
	"encoding/json"
	"io"
)

type Settings struct {
	Digitakt *DeviceSettings
	Digitone *DeviceSettings
}

func NewSettings() *Settings {
	return &Settings{
		Digitakt: &DeviceSettings{
			Inputs:     make(map[string]struct{}),
			Outputs:    make(map[string]struct{}),
			Channels:   make(map[int]struct{}),
			ProgChgSrc: DeviceSourceDigitakt,
			MuteSrc:    DeviceSourceDigitakt,
		},
		Digitone: &DeviceSettings{
			Inputs:     make(map[string]struct{}),
			Outputs:    make(map[string]struct{}),
			Channels:   make(map[int]struct{}),
			ProgChgSrc: DeviceSourceDigitone,
			MuteSrc:    DeviceSourceDigitone,
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

type DeviceSettings struct {
	Inputs       map[string]struct{}
	Outputs      map[string]struct{}
	Channels     map[int]struct{}
	ProgChgSrc   DeviceSource
	MuteSrc      DeviceSource
	ProgChgInCh  int
	ProgChgOutCh int
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
