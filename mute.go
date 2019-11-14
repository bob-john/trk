package main

import (
	"strings"

	"github.com/nsf/termbox-go"
)

type Mute []bool

func ParseMute(str string, channelCount int) Mute {
	m := make(Mute, channelCount)
	for _, ch := range str {
		n := int(ch) - '1'
		if n < 0 || n >= len(m) {
			continue
		}
		m[n] = true
	}
	for n, v := range m {
		m[n] = !v
	}
	return m
}

func (m Mute) String() string {
	str := make([]rune, len(m))
	for n, v := range m {
		if v {
			str[n] = '-'
		} else {
			str[n] = '1' + rune(n)
		}
	}
	return string(str)
}

func (m Mute) Clear() {
	for n := range m {
		m[n] = true
	}
}

type muteCell struct {
	stringCell
	channelCount int
}

func newMuteCell(doc *Arrangement, row, col, channelCount int) Cell {
	return muteCell{stringCell{doc, row, col}, channelCount}
}

func (c muteCell) Edit() CellEditor {
	return newMuteCellEditor(&c)
}

type muteCellEditor struct {
	*muteCell
	mute Mute
}

func newMuteCellEditor(c *muteCell) CellEditor {
	return &muteCellEditor{c, ParseMute(c.String(), c.channelCount)}
}

func (c *muteCellEditor) Input(e termbox.Event) {
	if isKeyDelete(e) {
		c.mute.Clear()
		c.Set(strings.Repeat(".", len(c.mute)))
		return
	}
	if e.Type == termbox.EventKey && e.Ch == '-' {
		c.mute.Clear()
		c.Set(c.mute.String())
		return
	}
	if !isKeyDigit(e) {
		return
	}
	n := int(e.Ch) - '1'
	if n < 0 || n >= len(c.mute) {
		return
	}
	c.mute[n] = !c.mute[n]
	c.Set(c.mute.String())
}

func (c *muteCellEditor) Commit() {}
