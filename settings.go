package main

import (
	"strings"

	"github.com/nsf/termbox-go"
)

type SettingsItem interface {
	String(int) string
	Handle(termbox.Event) *Settings
}

type Settings struct {
	Title string
	Items []SettingsItem
}

type Checkbox struct {
	Label string
	On    bool
}

func (c *Checkbox) String(width int) string {
	box := "[ ]"
	if c.On {
		box = "[x]"
	}
	return LayoutString(c.Label, box, width)
}

func (c *Checkbox) Handle(e termbox.Event) *Settings {
	if e.Type != termbox.EventKey {
		return nil
	}
	switch e.Key {
	case termbox.KeyEnter:
		c.On = !c.On
	}
	return nil
}

func LayoutString(lhs, rhs string, width int) string {
	if len(lhs)+1+len(rhs) > width {
		lhs = strings.TrimSpace(lhs[:width-len(rhs)-5]) + "..."
	}
	return lhs + strings.Repeat(" ", width-len(rhs)-1-len(lhs)) + rhs
}
