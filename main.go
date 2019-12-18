package main

import (
	"fmt"
	"trk/elektron"
	"trk/ui"

	"gitlab.com/gomidi/midi"
	"gitlab.com/gomidi/midi/midimessage/channel"
)

var (
	DT = elektron.Digitakt()
	DN = elektron.Digitone()
	LP = ui.Input("LPMiniMK3")

	C01 = elektron.C01
	C02 = elektron.C02

	C05 = elektron.C05
	C06 = elektron.C06
	C07 = elektron.C07
	C08 = elektron.C08
)

func main() {
	defer ui.Close()

	//
	// SETUP
	//

	DT.SetProgChgOutCh(16)

	DN.SetChannel(1, 9)
	DN.SetChannel(2, 10)
	DN.SetChannel(3, 11)
	DN.SetChannel(4, 12)
	DN.SetProgChgOutCh(16)

	//
	// PLAY
	//

	LP.Listen(func(m midi.Message) bool {
		fmt.Println(m)
		switch m := m.(type) {
		case channel.ControlChange:
			return m.Controller() != 19
		default:
			return true
		}
	})

	// DT.Schedule("C01", "---45---")
	// DN.Schedule("C01", "1---")

	// DT.Schedule("C01", "-2345---")
	// DN.Schedule("C02", "12--")

	// DT.Schedule("C01", "-2345---")
	// DN.Schedule("C01", "12--")

	// DT.Schedule("C01", "-2345---")
	// DN.Schedule("C02", "12--")
}
