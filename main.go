package main

import (
	"fmt"
	"io"
	"log"

	"gitlab.com/gomidi/midi/mid"
	"gitlab.com/gomidi/midi/midimessage/realtime"
	"gitlab.com/gomidi/midi/midireader"
)

func main() {
	fmt.Println("TRK")

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

	in, err := mid.OpenIn(midiDriver{}, -1, "Elektron Digitone")
	must(err)
	defer in.Close()

	r, w := io.Pipe()

	dataC := make(chan []byte)
	go func() {
		for data := range dataC {
			w.Write(data)
		}
	}()

	err = in.SetListener(func(data []byte, deltaMicroseconds int64) {
		dataC <- data
	})

	rd := midireader.New(r, func(m realtime.Message) {
		if m != realtime.TimingClock {
			fmt.Println(m)
		}
	})
	for {
		m, err := rd.Read()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(m)
	}

	/*out, err := mid.OpenOut(midiDriver{}, -1, "Microsoft GS Wavetable Synth")
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

	time.Sleep(1 * time.Hour)*/
}

func must(err error) {
	if err != nil {
		panic(err.Error())
	}
}
