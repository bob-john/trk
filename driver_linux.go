package main

import (
	"gitlab.com/gomidi/midi/mid"
	driver "gitlab.com/gomidi/rtmididrv"
)

func NewDriver() (mid.Driver, error) {
	return driver.New()
}
