package main

import (
	"fmt"
	"time"

	"gitlab.com/gomidi/midi"
	"gitlab.com/gomidi/midi/mid"
)

func main() {
	fmt.Println("trk")

	var d midiDriver

	ins, err := d.Ins()
	must(err)

	outs, err := d.Outs()
	must(err)

	for _, port := range ins {
		fmt.Printf("[%v] %s\n", port.Number(), port.String())
	}
	for _, port := range outs {
		fmt.Printf("[%v] %s\n", port.Number(), port.String())
	}

	in, err := mid.OpenIn(midiDriver{}, -1, "Keystation Mini 32")
	must(err)
	defer in.Close()

	out, err := mid.OpenOut(midiDriver{}, -1, "Microsoft GS Wavetable Synth")
	must(err)
	defer out.Close()

	var (
		w = mid.ConnectOut(out)
		r = mid.NewReader()
	)

	r.Msg.Each = func(pos *mid.Position, msg midi.Message) {
		w.Write(msg)
	}

	mid.ConnectIn(in, r)

	time.Sleep(1 * time.Hour)
}

func must(err error) {
	if err != nil {
		panic(err.Error())
	}
}
