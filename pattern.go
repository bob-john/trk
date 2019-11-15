package main

import (
	"fmt"
	"strconv"
	"strings"
	"unicode"

	"github.com/gomidi/midi/midimessage/channel"
	"github.com/nsf/termbox-go"
)

type Pattern int

func ParsePattern(str string) (Pattern, bool) {
	if len(str) == 0 {
		return 0, false
	}
	str = strings.ToUpper(str)
	bank := int(str[0] - 'A')
	if bank < 0 || bank >= 8 {
		return MakePattern(0, 0), false
	}
	trig, err := strconv.Atoi(strings.TrimPrefix(str[1:], "0"))
	if err != nil || trig < 1 || trig > 16 {
		return MakePattern(bank, 0), false
	}
	return MakePattern(bank, trig-1), true
}

func MakePattern(bank, trig int) Pattern {
	return Pattern(bank*16 + trig)
}

func (p Pattern) String() string {
	return fmt.Sprintf("%s%02d", string('A'+int(p)/16), 1+int(p)%16)
}

// func (p Pattern) Bank() int {
// 	return int(p) / 16
// }

// func (p Pattern) Trig() int {
// 	return int(p) % 16
// }

// func (p Pattern) SetBank(bank int) Pattern {
// 	return MakePattern(bank, p.Trig())
// }

// func (p Pattern) SetTrig(trig int) Pattern {
// 	return MakePattern(p.Bank(), trig)
// }

type patternCell struct {
	stringCell
}

func newPatternCell(doc *Arrangement, row, col int) Cell {
	return patternCell{stringCell{doc, row, col}}
}

func (c patternCell) Edit() CellEditor {
	return newPatternCellEditor(&c)
}

func (c patternCell) Output(out *Device) {
	p, ok := ParsePattern(c.String())
	if !ok {
		return
	}
	out.Write(channel.Channel9.ProgramChange(uint8(p)))
}

type patternCellEditor struct {
	*patternCell
	buffer string
}

func newPatternCellEditor(c *patternCell) CellEditor {
	return &patternCellEditor{c, ""}
}

func (c *patternCellEditor) Input(e termbox.Event) {
	if e.Type != termbox.EventKey {
		return
	}
	if isKeyBackspace(e) {
		c.buffer = ""
		c.Set("...")
		return

	}
	var ok bool
	if len(c.buffer) == 0 {
		ch := unicode.ToUpper(e.Ch)
		ok = ch >= 'A' && ch <= 'H'
	} else {
		_, ok = ParsePattern(c.buffer + string(e.Ch))
	}
	if !ok {
		return
	}
	c.buffer += string(e.Ch)
	c.Set(pad(c.buffer, ' ', 3))
}

func (c *patternCellEditor) Commit() {
	if c.buffer == "" || c.buffer == "..." {
		return
	}
	p, _ := ParsePattern(c.buffer)
	c.Set(p.String())
}
