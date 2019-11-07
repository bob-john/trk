package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/gomidi/midi/midimessage/realtime"
)

func main() {
	var drv midiDriver

	if len(os.Args) < 2 {
		fmt.Println("trk: missing a command. See 'trk help'.")
		os.Exit(1)
	}

	switch strings.ToLower(os.Args[1]) {
	case "help":
		fmt.Println("usage: trk <command> [<args>]")
		fmt.Println()
		fmt.Println("These are the available commands:")
		fmt.Println("  ", "devices", "\t", "List the MIDI devices")
		fmt.Println("  ", "record", "\t", "Record MIDI devices")

	case "devices":
		ins, err := drv.Ins()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("MIDI Inputs:")
		for _, port := range ins {
			fmt.Printf("%d) %s\n", port.Number(), port.String())
		}
		outs, err := drv.Outs()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println()
		fmt.Println("MIDI Outputs:")
		for _, port := range outs {
			fmt.Printf("%d) %s\n", port.Number(), port.String())
		}

	case "record":
		if len(os.Args) == 2 {
			fmt.Println("trk: no device to record")
			os.Exit(1)
		}
		quit := make(chan struct{})
		for _, name := range os.Args[2:] {
			in, err := OpenInput(name)
			if err != nil {
				fmt.Printf("trk: fail to open %q: %v\n", name, err)
				os.Exit(2)
			}
			defer in.Close()
			go func(name string) {
				for m := range in.In() {
					if m != realtime.TimingClock {
						fmt.Println(name, m)
					}
					if m == realtime.Stop {
						close(quit)
					}
				}
			}(name)
		}
		<-quit

	default:
		fmt.Printf("trk: '%s' is not a command. See 'trk help'.\n")
		os.Exit(1)
	}
}
