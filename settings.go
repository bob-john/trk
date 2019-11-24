package main

import (
	"encoding/json"
	"io"
)

type Settings struct {
	Digitakt *DeviceSettings
	Digitone *DeviceSettings
}

type DeviceSettings struct {
	Input          []string
	Output         []string
	Channels       []int
	ProgChgSrc     DeviceSource
	MuteSrc        DeviceSource
	AutoChannel    int
	ProgChgInChan  int
	ProgChgOutChan int
}

type DeviceSource int

const (
	DeviceSourceDigitakt DeviceSource = iota
	DeviceSourceDigitone
	DeviceSourceBoth
)

func NewSettings() *Settings {
	return &Settings{
		Digitakt: &DeviceSettings{
			Channels:    []int{1, 2, 3, 4, 5, 6, 7, 8},
			ProgChgSrc:  DeviceSourceDigitakt,
			MuteSrc:     DeviceSourceDigitakt,
			AutoChannel: 10,
		},
		Digitone: &DeviceSettings{
			Channels:    []int{1, 2, 3, 4},
			ProgChgSrc:  DeviceSourceDigitone,
			MuteSrc:     DeviceSourceDigitone,
			AutoChannel: 10,
		},
	}
}

func ReadSettings(f io.Reader) (*Settings, error) {
	s := NewSettings()
	return s, nil
}

func (s *Settings) Write(f io.Writer) error {
	return json.NewEncoder(f).Encode(s)
}
