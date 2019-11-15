package main

import (
	"strconv"

	"github.com/nsf/termbox-go"
)

type lenCell struct {
	stringCell
}

func newLenCell(doc *Arrangement, row, col int) Cell {
	return lenCell{stringCell{doc, row, col}}
}

func (c lenCell) Edit() CellEditor {
	return newLenCellEditor(&c)
}

func (c lenCell) Output(*Device) {}

type lenCellEditor struct {
	*lenCell
	buffer string
}

func newLenCellEditor(c *lenCell) CellEditor {
	return &lenCellEditor{c, ""}
}

func (c *lenCellEditor) Input(e termbox.Event) {
	if !isKeyDigit(e) {
		return
	}
	if c.buffer == "" && e.Ch == '0' {
		return
	}
	n, err := strconv.Atoi(c.buffer + string(e.Ch))
	if err != nil || n > 1024 {
		return
	}
	c.buffer += string(e.Ch)
	c.Set(c.buffer)
}

func (c *lenCellEditor) Commit() {
	if c.buffer == "" {
		return
	}
	n, _ := strconv.Atoi(c.buffer)
	n = clamp(n, 1, 1024)
	c.Set(strconv.Itoa(n))
}
