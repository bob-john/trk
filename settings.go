package main

import (
	"strings"

	"github.com/nsf/termbox-go"
)

type Settings struct {
	Title string
	Items []settingsItem
}

func NewSettings(title string) *Settings {
	return &Settings{title, nil}
}

func (s *Settings) AddMenu(title string, build func(page *Settings)) {
	page := &Settings{title, nil}
	s.Items = append(s.Items, &menuItem{title, page})
	build(page)
}

func (s *Settings) AddCheckbox(title string, on bool) {
	s.Items = append(s.Items, &checkboxItem{title, on})
}

func (s *Settings) PreferredSize() Size {
	sz := MakeSize(len(s.Title)+4, len(s.Items)+1)
	for _, item := range s.Items {
		minWidth := item.MinWidth() + 2
		if sz.Width < minWidth {
			sz.Width = minWidth
		}
	}
	return sz
}

type settingsItem interface {
	Handle(*Dialog, termbox.Event)
	String(int) string
	MinWidth() int
}

type menuItem struct {
	label string
	page  *Settings
}

func (m *menuItem) Handle(dialog *Dialog, e termbox.Event) {
	if IsKey(e, termbox.KeyArrowRight, termbox.KeyEnter) {
		dialog.Enter(m.page)
	}
}

func (m *menuItem) String(width int) string {
	return LayoutString(m.label, "", width)
}

func (m *menuItem) MinWidth() int {
	return len(m.label)
}

type checkboxItem struct {
	label string
	on    bool
}

func (c *checkboxItem) Handle(dialog *Dialog, e termbox.Event) {
	if IsKey(e, termbox.KeyArrowRight, termbox.KeyEnter) {
		c.on = !c.on
	}
}

func (c *checkboxItem) String(width int) string {
	box := "[ ]"
	if c.on {
		box = "[x]"
	}
	return LayoutString(c.label, box, width)
}

func (c *checkboxItem) MinWidth() int {
	return len(c.label) + 4
}

func LayoutString(lhs, rhs string, width int) string {
	if len(lhs)+1+len(rhs) > width {
		lhs = strings.TrimSpace(lhs[:width-len(rhs)-5]) + "..."
	}
	return lhs + strings.Repeat(" ", width-len(rhs)-1-len(lhs)) + rhs
}
