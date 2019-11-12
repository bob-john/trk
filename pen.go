package main

import (
	"github.com/nsf/termbox-go"
)

type Pen struct {
	arr      *Arrangement
	row, col int
	input    []termbox.Event
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
			p.input = append(p.input, e)
			p.Cell().Input(p.input)
		}
	}
	p.row = clamp(p.row, 0, p.arr.RowCount()-1)
	p.col = clamp(p.col, 0, p.arr.Row(p.row).CellCount()-1)
	if p.row != oldRow || p.col != oldCol {
		p.input = nil
	}
}
