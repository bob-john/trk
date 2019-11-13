package main

import (
	"strconv"
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
	Reset()
	Commit()
}

type indexCellEditor struct {
	*indexCell
}

func newIndexCellEditor(c *indexCell) CellEditor {
	return &indexCellEditor{c}
}

func (c *indexCellEditor) Input(e termbox.Event) {}
func (c *indexCellEditor) Reset()                {}
func (c *indexCellEditor) Commit()               {}

type patternCellEditor struct {
	*patternCell
}

func newPatternCellEditor(c *patternCell) CellEditor {
	return &patternCellEditor{c}
}

func (c *patternCellEditor) Input(e termbox.Event) {

}

//TODO Clear()

func (c *patternCellEditor) Reset()  {}
func (c *patternCellEditor) Commit() {}

type muteCellEditor struct {
	*muteCell
	unmuted []bool
}

func newMuteCellEditor(c *muteCell) CellEditor {
	u := make([]bool, c.len)
	for _, ch := range c.String() {
		n := int(ch - '1')
		if n >= 0 && n < len(u) {
			u[n] = true
		}
	}
	return &muteCellEditor{c, u}
}

func (c *muteCellEditor) Input(e termbox.Event) {
	if !isKeyDigit(e) {
		return
	}
	n := int(e.Ch - '1')
	if n < 0 || n >= len(c.unmuted) {
		return
	}
	c.unmuted[n] = !c.unmuted[n]
	var val string
	for n, u := range c.unmuted {
		if u {
			val += strconv.Itoa(1 + n)
		} else {
			val += "-"
		}
	}
	c.Set(c.row, c.col, val)
}

// func (c *muteCellEditor) Clear()  {}
func (c *muteCellEditor) Reset()  {}
func (c *muteCellEditor) Commit() {}

type lenCellEditor struct {
	*lenCell
}

func newLenCellEditor(c *lenCell) CellEditor {
	return &lenCellEditor{c}
}

func (c *lenCellEditor) Input(e termbox.Event) {}
func (c *lenCellEditor) Reset()                {}
func (c *lenCellEditor) Commit()               {}

func isKeyDelete(e termbox.Event) bool {
	return e.Type == termbox.EventKey && e.Key == termbox.KeyDelete
}

func isKeyDigit(e termbox.Event) bool {
	return e.Type == termbox.EventKey && unicode.IsDigit(e.Ch)
}

func isKeyLetter(e termbox.Event) bool {
	return e.Type == termbox.EventKey && unicode.IsLetter(e.Ch)
}

func isKeyEnter(e termbox.Event) bool {
	return e.Type == termbox.EventKey && e.Key == termbox.KeyEnter
}
