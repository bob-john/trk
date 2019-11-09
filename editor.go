package main

import (
	"fmt"
	"strconv"
)

type LineEditor struct {
	str string
}

func (e *LineEditor) Reset(str string) {
	e.str = str
}

func (e *LineEditor) CellCount() int {
	return 16
}

func (e *LineEditor) Cell(i int) Cell {
	switch i {
	case 0:
		return BankCell{e, 4}
	case 1:
		return PatternCell{e, 5}
	case 2, 3, 4, 5, 6, 7, 8, 9:
		return MuteCell{e, 8 + i - 2}
	case 10:
		return BankCell{e, 17}
	case 11:
		return PatternCell{e, 18}
	case 12, 13, 14, 15:
		return MuteCell{e, 21 + i - 12}
	}
	return nil
}

func (e *LineEditor) Line() string {
	return e.str
}

func (e *LineEditor) Replace(i int, repl string) {
	e.str = e.str[:i] + repl + e.str[i+len(repl):]
}

type Range struct {
	Index, Len int
}

func (r Range) Substr(str string) string {
	return str[r.Index : r.Index+r.Len]
}

type Cell interface {
	Index() int
	String() string
	Inc()
	Dec()
}

type BankCell struct {
	editor *LineEditor
	index  int
}

func (c BankCell) Index() int {
	return c.index
}

func (c BankCell) String() string {
	return c.editor.Line()[c.index : c.index+1]
}

func (c BankCell) Inc() {
	str := c.String()
	next := str
	switch str {
	case ".":
		next = "A"
	case "A", "B", "C", "D", "E", "F", "G":
		next = string(str[0] + 1)
	}
	c.editor.Replace(c.Index(), next)
}

func (c BankCell) Dec() {
	str := c.String()
	next := str
	switch str {
	case "B", "C", "D", "E", "F", "G", "H":
		next = string(str[0] - 1)
	}
	c.editor.Replace(c.Index(), next)
}

type PatternCell struct {
	editor *LineEditor
	index  int
}

func (c PatternCell) Index() int {
	return c.index
}

func (c PatternCell) String() string {
	return c.editor.Line()[c.index : c.index+2]
}

func (c PatternCell) Inc() {
	var next string
	str := c.String()
	switch str {
	case ".":
		next = "01"
	default:
		val, _ := strconv.Atoi(str)
		if val < 16 {
			val++
		}
		next = fmt.Sprintf("%02d", val)
	}
	c.editor.Replace(c.Index(), next)
}

func (c PatternCell) Dec() {
	var next string
	str := c.String()
	switch str {
	case ".":
		return
	default:
		val, _ := strconv.Atoi(str)
		if val > 1 {
			val--
		}
		next = fmt.Sprintf("%02d", val)
	}
	c.editor.Replace(c.Index(), next)
}

type MuteCell struct {
	editor *LineEditor
	index  int
}

func (c MuteCell) Index() int {
	return c.index
}

func (c MuteCell) String() string {
	return c.editor.Line()[c.index : c.index+1]
}

func (c MuteCell) Inc() {
	var next string
	switch c.String() {
	case ".", "-":
		next = "+"
	default:
		return
	}
	c.editor.Replace(c.Index(), next)
}

func (c MuteCell) Dec() {
	var next string
	switch c.String() {
	case ".", "+":
		next = "-"
	default:
		return
	}
	c.editor.Replace(c.Index(), next)
}
