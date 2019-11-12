package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"strings"
	"unicode"

	"github.com/nsf/termbox-go"
)

type Arrangement struct {
	rows [][]string
}

func LoadArrangement(path string) (*Arrangement, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	rows, err := csv.NewReader(f).ReadAll()
	if err != nil {
		return nil, err
	}
	return &Arrangement{rows}, f.Close()
}

func NewArrangement() *Arrangement {
	return &Arrangement{}
}

func (a *Arrangement) RowCount() int {
	return len(a.rows)
}

func (a *Arrangement) Row(i int) Row {
	return Row{a, i}
}

func (a *Arrangement) Cell(row, col int) Cell {
	return a.Row(row).Cell(col)
}

func (a *Arrangement) Set(row, col int, value string) {
	a.rows[row][col] = value
}

func (a *Arrangement) Get(row, col int) string {
	return a.rows[row][col]
}

type Row struct {
	*Arrangement
	row int
}

func (r Row) CellCount() int {
	return 6
}

func (r Row) Cell(i int) Cell {
	switch i {
	case 0:
		return r.Index()
	case 1:
		return r.Digitakt().Pattern()
	case 2:
		return r.Digitakt().Mute()
	case 3:
		return r.Digitone().Pattern()
	case 4:
		return r.Digitone().Mute()
	case 5:
		return r.Len()
	}
	return nil
}

func (r Row) Range(i int) Range {
	var index int
	for j := 0; j < i; j++ {
		index += len(r.Cell(j).String()) + 1
	}
	return Range{index, len(r.Cell(i).String())}
}

func (r Row) Index() Cell {
	return indexCell{r.Arrangement, r.row, 0}
}

func (r Row) Digitakt() Part {
	return Part{r.Arrangement, r.row, 0, 8}
}

func (r Row) Digitone() Part {
	return Part{r.Arrangement, r.row, 2, 4}
}

func (r Row) Len() Cell {
	return lenCell{r.Arrangement, r.row, 4}
}

func (r Row) String() string {
	var cells []string
	for i := 0; i < r.CellCount(); i++ {
		cells = append(cells, r.Cell(i).String())
	}
	return strings.Join(cells, " ")
}

type Part struct {
	*Arrangement
	row, col     int
	channelCount int
}

func (p Part) Pattern() Cell {
	return patternCell{p.Arrangement, p.row, p.col}
}

func (p Part) Mute() Cell {
	return muteCell{p.Arrangement, p.row, p.col + 1, p.channelCount}
}

type Cell interface {
	Input([]termbox.Event)
	String() string
}

type indexCell struct {
	*Arrangement
	row, col int
}

func (c indexCell) Input(events []termbox.Event) {}

func (c indexCell) String() string {
	return fmt.Sprintf("%3d", 1+c.row)
}

type patternCell struct {
	*Arrangement
	row, col int
}

func (c patternCell) Input(events []termbox.Event) {}

func (c patternCell) String() string {
	return c.Get(c.row, c.col)
}

type muteCell struct {
	*Arrangement
	row, col int
	len      int
}

func (c muteCell) Input(events []termbox.Event) {
	muted := make([]bool, c.len)
	for _, e := range events {
		if !isKeyDigit(e) {
			continue
		}
		n := int(e.Ch - '1')
		if n >= 0 && n < c.len {
			muted[n] = !muted[n]
		}
	}
	var val string
	for n, m := range muted {
		if m {
			val += "-"
		} else {
			val += strconv.Itoa(1 + n)
		}
	}
	c.Set(c.row, c.col, val)
}

func (c muteCell) String() string {
	return c.Get(c.row, c.col)
}

type lenCell struct {
	*Arrangement
	row, col int
}

func (c lenCell) Input(events []termbox.Event) {}

func (c lenCell) String() string {
	return c.Get(c.row, c.col)
}

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
