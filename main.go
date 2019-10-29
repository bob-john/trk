package main

import (
	"fmt"
	"time"

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
	err = mid.ConnectIn(in, mid.NewReader())
	must(err)
	time.Sleep(1 * time.Hour)
}

func must(err error) {
	if err != nil {
		panic(err.Error())
	}
}
