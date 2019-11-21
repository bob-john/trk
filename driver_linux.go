package main

import (
	"gitlab.com/gomidi/midi/mid"
	"gitlab.com/gomidi/rtmididrv"
)

func NewDriver() (mid.Driver, error) {
	return rtmididrv.New()
}
