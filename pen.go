package main

import (
	"unicode"

	"github.com/nsf/termbox-go"
)

type Pen struct {
	arr      *Arrangement
	row, col int
	input    string
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
	switch e.Type {
	case termbox.EventKey:
		switch e.Key {
		case termbox.KeyArrowUp:
			p.row--
			p.input = ""
		case termbox.KeyPgup:
			p.row -= pageSize
			p.input = ""
		case termbox.KeyHome:
			p.row = 0
			p.input = ""
		case termbox.KeyArrowDown:
			p.row++
			p.input = ""
		case termbox.KeyPgdn:
			p.row += pageSize
			p.input = ""
		case termbox.KeyEnd:
			p.row = p.arr.RowCount() - 1
			p.input = ""

		case termbox.KeyArrowLeft:
			p.col--
			p.input = ""
		case termbox.KeyArrowRight:
			p.col++
			p.input = ""

		default:
			if unicode.IsPrint(e.Ch) {
				p.input += string(e.Ch)
				p.Cell().Set(p.input, true)
			}

			// default:
			// 	if arr.Cell(cur).Input(e) {
			// 		p.row++
			// 	}
		}
	}
	p.row = clamp(p.row, 0, p.arr.RowCount()-1)
	p.col = clamp(p.col, 0, p.arr.Row(p.row).CellCount()-1)
}
