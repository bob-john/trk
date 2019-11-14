package main

import (
	"unicode"

	"github.com/nsf/termbox-go"
)

type Pen struct {
	doc      *Arrangement
	row, col int
	editor   CellEditor
	undo     string
}

func NewPen(doc *Arrangement) *Pen {
	return &Pen{doc: doc}
}

func (p *Pen) Row() int {
	return p.row
}

func (p *Pen) Cell() Cell {
	return p.doc.Cell(p.row, p.col)
}

func (p *Pen) Range() Range {
	return p.doc.Row(p.row).Range(p.col)
}

func (p *Pen) Editing() bool {
	return p.editor != nil
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
			p.row = p.doc.RowCount() - 1

		case termbox.KeyArrowLeft:
			if p.editor != nil {
				p.col--
			}
		case termbox.KeyArrowRight:
			if p.editor != nil {
				p.col++
			}

		case termbox.KeyEsc:
			if p.editor != nil {
				p.Cell().Set(p.undo)
				p.editor = p.Cell().Edit()
			}
		case termbox.KeyEnter:
			if p.editor != nil {
				p.editor.Commit()
				p.editor = nil
				p.undo = ""
			} else {
				p.undo = p.Cell().String()
				p.editor = p.Cell().Edit()
			}

		case termbox.KeyInsert:
			p.doc.InsertAfter(p.row)
		case termbox.KeyDelete:
			if p.editor != nil {
				p.editor.Input(e)
			} else {
				p.doc.Delete(p.Row())
			}

		default:
			if p.editor != nil {
				p.editor.Input(e)
			}
		}
	}
	if p.col < 0 && p.row > 0 {
		p.row--
		p.col = p.doc.Row(p.row).CellCount() - 1
	} else if p.col >= p.doc.Row(p.row).CellCount() && p.row < p.doc.RowCount()-1 {
		p.row++
		p.col = 0
	}
	p.row = clamp(p.row, 0, p.doc.RowCount()-1)
	p.col = clamp(p.col, 0, p.doc.Row(p.row).CellCount()-1)
	if p.editor != nil && (p.row != oldRow || p.col != oldCol) {
		p.editor.Commit()
		p.undo = p.Cell().String()
		p.editor = p.Cell().Edit()
	}
}

type CellEditor interface {
	Input(termbox.Event)
	Commit()
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
