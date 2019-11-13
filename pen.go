package main

import (
	"strconv"
	"strings"
	"unicode"

	"github.com/nsf/termbox-go"
)

type Pen struct {
	arr      *Arrangement
	row, col int
	editor   CellEditor
}

func NewPen(arr *Arrangement) *Pen {
	return &Pen{arr: arr}
}

func (p *Pen) Row() int {
	return p.row
}

func (p *Pen) Cell() Cell {
	return p.arr.Cell(p.row, p.col)
}

func (p *Pen) Range() Range {
	return p.arr.Row(p.row).Range(p.col)
}

func (p *Pen) Handle(e termbox.Event) {
	oldRow, oldCol := p.row, p.col
	switch e.Type {
	case termbox.EventKey:
		switch e.Key {
		case termbox.KeyArrowUp:
			p.row--
		case termbox.KeyPgup:
			p.row -= pageSize
		case termbox.KeyHome:
			p.row = 0
		case termbox.KeyArrowDown:
			p.row++
		case termbox.KeyPgdn:
			p.row += pageSize
		case termbox.KeyEnd:
			p.row = p.arr.RowCount() - 1

		case termbox.KeyArrowLeft:
			p.col--
		case termbox.KeyArrowRight:
			p.col++

		default:
			if p.editor == nil {
				p.editor = p.Cell().Edit()
			}
			p.editor.Input(e)
		}
	}
	if p.col < 0 && p.row > 0 {
		p.row--
		p.col = p.arr.Row(p.row).CellCount() - 1
	} else if p.col >= p.arr.Row(p.row).CellCount() && p.row < p.arr.RowCount()-1 {
		p.row++
		p.col = 0
	}
	p.row = clamp(p.row, 0, p.arr.RowCount()-1)
	p.col = clamp(p.col, 0, p.arr.Row(p.row).CellCount()-1)
	if p.row != oldRow || p.col != oldCol {
		if p.editor != nil {
			p.editor.Commit()
		}
		p.editor = p.Cell().Edit()
	}
}

type CellEditor interface {
	Input(termbox.Event)
	Commit()
}

type indexCellEditor struct {
	*indexCell
}

func newIndexCellEditor(c *indexCell) CellEditor {
	return &indexCellEditor{c}
}

func (c *indexCellEditor) Input(e termbox.Event) {}
func (c *indexCellEditor) Commit()               {}

type patternCellEditor struct {
	*patternCell
	old, buffer string
}

func newPatternCellEditor(c *patternCell) CellEditor {
	return &patternCellEditor{c, c.String(), ""}
}

func (c *patternCellEditor) Input(e termbox.Event) {
	if e.Type != termbox.EventKey {
		return
	}
	switch e.Key {
	case termbox.KeyDelete:
		c.buffer = "..."

	case termbox.KeyEsc:
		c.buffer = c.old

	default:
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
	}
	c.Set(c.row, c.col, pad(c.buffer, ' ', 3))
}

func (c *patternCellEditor) Commit() {
	p, ok := ParsePattern(c.buffer)
	if ok {
		c.Set(c.row, c.col, p.String())
	}
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
		c.Set(c.row, c.col, strings.Repeat(".", len(c.mute)))
		return
	}
	if e.Type == termbox.EventKey && e.Ch == '-' {
		c.mute.Clear()
		c.Set(c.row, c.col, c.mute.String())
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
	c.Set(c.row, c.col, c.mute.String())
}

func (c *muteCellEditor) Commit() {}

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
	n, err := strconv.Atoi(c.buffer + string(e.Ch))
	if err != nil || n > 1024 {
		return
	}
	c.buffer += string(e.Ch)
	c.Set(c.row, c.col, c.buffer)
}

func (c *lenCellEditor) Commit() {
	if c.buffer == "" {
		return
	}
	n, _ := strconv.Atoi(c.buffer)
	n = clamp(n, 1, 1024)
	c.Set(c.row, c.col, strconv.Itoa(n))
}

func isKeyLetter(e termbox.Event) bool {
	return e.Type == termbox.EventKey && unicode.IsLetter(e.Ch)
}

func isKeyDigit(e termbox.Event) bool {
	return e.Type == termbox.EventKey && unicode.IsDigit(e.Ch)
}

func isKeyDelete(e termbox.Event) bool {
	return e.Type == termbox.EventKey && e.Key == termbox.KeyDelete
}

func isKeyBackspace(e termbox.Event) bool {
	return e.Type == termbox.EventKey && e.Key == termbox.KeyBackspace
}
