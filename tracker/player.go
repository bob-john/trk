package tracker

import (
	"trk/rtmididrv"

	"gitlab.com/gomidi/midi"
	"gitlab.com/gomidi/midi/mid"
)

var (
	midiDriver mid.Driver
	outs       = make(map[string]mid.Out)
)

func init() {
	var err error
	midiDriver, err = rtmididrv.New()
	must(err)
}

func Play(port string, message midi.Message) {
	out, ok := outs[port]
	if !ok {
		var err error
		out, err = mid.OpenOut(midiDriver, -1, port)
		must(err)
		outs[port] = out
	}
	out.Send(message.Raw())
}
