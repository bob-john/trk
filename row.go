package main

import (
	"fmt"
	"strings"
)

func makeEmptyRow() []string {
	return []string{"...", "........", "...", "....", "64"}
}

type Row struct {
	doc *Arrangement
	row int
}

func (r Row) CellCount() int {
	return 5
}

func (r Row) Cell(col int) Cell {
	switch col {
	case 0:
		return newPatternCell(r.doc, r.row, 0)
	case 1:
		return newMuteCell(r.doc, r.row, 1, 8)
	case 2:
		return newPatternCell(r.doc, r.row, 2)
	case 3:
		return newMuteCell(r.doc, r.row, 3, 4)
	case 4:
		return newLenCell(r.doc, r.row, 4)
	}
	return nil
}

func (r Row) Range(i int) Range {
	index := 4
	for j := 0; j < i; j++ {
		index += len(r.Cell(j).String()) + 1
	}
	return Range{index, len(r.Cell(i).String())}
}

func (r Row) Index() Cell {
	return r.Cell(0)
}

func (r Row) Digitakt() Part {
	return Part{r, 0}
}

func (r Row) Digitone() Part {
	return Part{r, 2}
}

func (r Row) Len() Cell {
	return r.Cell(4)
}

func (r Row) String() string {
	cells := []string{fmt.Sprintf("%3d", 1+r.row)}
	for i := 0; i < r.CellCount(); i++ {
		cells = append(cells, r.Cell(i).String())
	}
	return strings.Join(cells, " ")
}

func (r Row) Output(digitakt, digitone *Device) {
	r.OutputPattern(digitakt, digitone)
	r.OutputMute(digitakt, digitone)
}

func (r Row) OutputPattern(digitakt, digitone *Device) {
	r.doc.ConsolidatedCell(r.row, 0).Output(digitakt)
	r.doc.ConsolidatedCell(r.row, 2).Output(digitone)
}

func (r Row) OutputMute(digitakt, digitone *Device) {
	r.doc.ConsolidatedCell(r.row, 1).Output(digitakt)
	r.doc.ConsolidatedCell(r.row, 3).Output(digitone)
}

type Part struct {
	row Row
	col int
}

func (p Part) Pattern() Cell {
	return p.row.Cell(p.col)
}

func (p Part) Mute() Cell {
	return p.row.Cell(p.col + 1)
}

func isCellEmpty(c Cell) bool {
	return c.String() == strings.Repeat(".", len(c.String()))
}
