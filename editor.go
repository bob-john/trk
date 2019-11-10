package main

import (
	"fmt"
	"strconv"
	"strings"
)

type LineEditor struct {
	str, prev string
	i         int
}

func (e *LineEditor) Reset(str, prev string) {
	e.str = str
	e.prev = prev
	e.i = 0
}

func (e *LineEditor) ActiveCell() Cell {
	return e.Cell(e.i)
}

func (e *LineEditor) MoveToNextCell() {
	if e.i < e.CellCount()-1 {
		e.i++
	}
}

func (e *LineEditor) MoveToPreviousCell() {
	if e.i > 0 {
		e.i--
	}
}

func (e *LineEditor) CellCount() int {
	return 14
}

func (e *LineEditor) Cell(i int) Cell {
	switch i {
	case 0:
		return PatternCell{e, 4, Range{4, 3}.Substr(e.prev)}
	case 1, 2, 3, 4, 5, 6, 7, 8:
		return MuteCell{e, 8, 8 + i - 1, Range{8, 8}.Substr(e.prev)}
	case 9:
		return PatternCell{e, 17, Range{17, 3}.Substr(e.prev)}
	case 10, 11, 12, 13:
		return MuteCell{e, 21, 21 + i - 10, Range{21, 4}.Substr(e.prev)}
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
	Clear()
}

type PatternCell struct {
	editor *LineEditor
	index  int
	old    string
}

func (c PatternCell) Index() int {
	return c.index
}

func (c PatternCell) String() string {
	return c.editor.Line()[c.index : c.index+3]
}

func (c PatternCell) Inc() {
	str := c.String()
	switch str {
	case "...":
		c.setDefaultValue()
	default:
		p := DecodePattern(str)
		if p < 127 {
			c.editor.Replace(c.index, EncodePattern(p+1))
		}
	}
}

func (c PatternCell) Dec() {
	str := c.String()
	switch str {
	case "...":
		c.setDefaultValue()
	default:
		p := DecodePattern(str)
		if p > 0 {
			c.editor.Replace(c.index, EncodePattern(p-1))
		}
	}
}

func (c PatternCell) setDefaultValue() {
	if c.old != "..." {
		c.editor.Replace(c.index, c.old)
	} else {
		c.editor.Replace(c.index, "A01")
	}
}

func (c PatternCell) Clear() {
	c.editor.Replace(c.index, "...")
}

type MuteCell struct {
	editor     *LineEditor
	groupIndex int
	index      int
	oldGroup   string
}

func (c MuteCell) Index() int {
	return c.index
}

func (c MuteCell) String() string {
	return c.editor.Line()[c.index : c.index+1]
}

func (c MuteCell) Inc() {
	switch c.String() {
	case ".":
		c.setDefaultValue()
	default:
		c.editor.Replace(c.Index(), "+")
		c.editor.MoveToNextCell()
	}
}

func (c MuteCell) Dec() {
	switch c.String() {
	case ".":
		c.setDefaultValue()
	default:
		c.editor.Replace(c.Index(), "-")
		c.editor.MoveToNextCell()
	}

}

func (c MuteCell) Clear() {
	c.editor.Replace(c.groupIndex, strings.Repeat(".", len(c.oldGroup)))
}

func (c MuteCell) setDefaultValue() {
	if strings.ContainsAny(c.oldGroup, "+-") {
		c.editor.Replace(c.groupIndex, c.oldGroup)
	} else {
		c.editor.Replace(c.groupIndex, strings.Repeat("+", len(c.oldGroup)))
	}
}

func DecodePattern(str string) int {
	if strings.ContainsAny(str, ".") {
		return -1
	}
	bank := int(str[0] - 'A')
	trig, _ := strconv.Atoi(str[1:])
	return bank*16 + trig - 1
}

func EncodePattern(val int) string {
	return fmt.Sprintf("%s%02d", string('A'+val/16), 1+val%16)
}
