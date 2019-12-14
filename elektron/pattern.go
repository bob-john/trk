package elektron

import (
	"gitlab.com/gomidi/midi"
	"gitlab.com/gomidi/midi/midimessage/channel"
)

type Pattern struct {
	bank, pattern int
}

func (p *Pattern) Chain(next ...*Pattern) *Chain {
	return &Chain{append([]*Pattern{p}, next...)}
}

func (p *Pattern) Message() midi.Message {
	return channel.Channel(15).ProgramChange(uint8(p.bank*16 + p.pattern))
}

type Chain struct {
	patterns []*Pattern
}
