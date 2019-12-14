package main

import (
	"trk/elektron"
	"trk/tracker"
	"trk/ui"
)

var (
	C01 = elektron.C01
	C02 = elektron.C02
	C05 = elektron.C05
	C06 = elektron.C06
	C07 = elektron.C07
	C08 = elektron.C08
)

func main() {
	ui, err := ui.New()
	must(err)
	trk, err := tracker.New()
	must(err)

	dt, err := ui.Out(trk, "Digitakt")
	must(err)
	dn, err := ui.Out(trk, "Digitone")
	must(err)

	ui.Close()

	must(dt.Play(C01))
	must(dn.Play(C01))

	C01.Chain(C02)
}
