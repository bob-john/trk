package main

import (
	"trk/elektron"
	"trk/tracker"
	"trk/ui"
)

var (
	DT = elektron.Digitakt()
	DN = elektron.Digitone()

	C01 = elektron.C01
	C02 = elektron.C02

	C05 = elektron.C05
	C06 = elektron.C06
	C07 = elektron.C07
	C08 = elektron.C08

	play = tracker.Play
)

func main() {
	DT.SetProgChgOutCh(16)

	DN.SetChannel(1, 9)
	DN.SetChannel(2, 10)
	DN.SetChannel(3, 11)
	DN.SetChannel(4, 12)
	DN.SetProgChgOutCh(16)

	ui.Close()

	play(DT.Pattern(C01))
	play(DT.Mute(1))
	play(DT.Mute(2))
	play(DT.Mute(3))
	play(DT.Unmute(4))
	play(DT.Unmute(5))
	play(DT.Mute(6))
	play(DT.Mute(7))
	play(DT.Mute(8))

	play(DN.Pattern(C01))
	play(DN.Unmute(1))
	play(DN.Mute(2))
	play(DN.Mute(3))
	play(DN.Mute(4))
}
