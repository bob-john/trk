package tracker

import (
	"trk/rtmididrv"

	"gitlab.com/gomidi/midi"
)

var (
	player *Player
)

func init() {
	driver, err := rtmididrv.New()
	must(err)
	player = NewPlayer(driver)
}

func Play(e Event) {
	must(player.Play(e.Port(), e.Message()))
}

type Event interface {
	Port() string
	Message() midi.Message
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
