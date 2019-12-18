package tracker

import (
	"io"

	"gitlab.com/gomidi/midi"
	"gitlab.com/gomidi/midi/mid"
	"gitlab.com/gomidi/midi/midimessage/realtime"
	"gitlab.com/gomidi/midi/midireader"
)

var (
	ins = make(map[string]*In)
)

type In struct {
	in mid.In
}

func OpenIn(port string) *In {
	in, err := mid.OpenIn(midiDriver, -1, port)
	must(err)
	return &In{in}
}

func (i *In) Listen(handler func(midi.Message) bool) {
	var (
		r, w = io.Pipe()
		msg  = make(chan midi.Message)
	)
	must(i.in.SetListener(func(b []byte, deltaMicroseconds int64) {
		w.Write(b)
	}))
	go func() {
		reader := midireader.New(r, func(m realtime.Message) {
			msg <- m
		})
		for {
			m, err := reader.Read()
			if err == io.EOF {
				must(r.Close())
				return
			}
			must(err)
			msg <- m
		}
	}()
	for m := range msg {
		if !handler(m) {
			must(w.Close())
			must(i.in.StopListening())
			return
		}
	}
}
