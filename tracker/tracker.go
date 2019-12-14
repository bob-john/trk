package tracker

import (
	"gitlab.com/gomidi/midi"
	"trk/rtmididrv"
)

type Tracker struct {
	player *Player
}

func New() (*Tracker, error) {
	var err error
	driver, err := rtmididrv.New()
	if err != nil {
		return nil, err
	}
	return &Tracker{player: NewPlayer(driver)}, nil
}

func (t *Tracker) Out(port string) *Out {
	return &Out{t, port}
}

type Out struct {
	tracker *Tracker
	port    string
}

func (o *Out) Play(e Event) error {
	return o.tracker.player.Play([]string{o.port}, e.Message())
}

type Event interface {
	Message() midi.Message
}
