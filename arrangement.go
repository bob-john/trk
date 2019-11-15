package main

import (
	"encoding/csv"
	"os"
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
	return &Arrangement{[][]string{makeEmptyRow()}}
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

func (a *Arrangement) ConsolidatedCell(row, col int) Cell {
	c := a.Cell(row, col)
	for isCellEmpty(c) && row > 0 {
		row--
		c = a.Cell(row, col)
	}
	return c
}

func (a *Arrangement) Set(row, col int, value string) {
	a.rows[row][col] = value
}

func (a *Arrangement) Get(row, col int) string {
	return a.rows[row][col]
}

func (a *Arrangement) Delete(row int) {
	a.rows = append(a.rows[:row], a.rows[row+1:]...)
	if len(a.rows) == 0 {
		a.rows = [][]string{makeEmptyRow()}
	}
}

func (a *Arrangement) InsertAfter(row int) {
	a.rows = append(a.rows[:row+1], append([][]string{makeEmptyRow()}, a.rows[row+1:]...)...)
}

func (a *Arrangement) WriteFile(path string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	w := csv.NewWriter(f)
	err = w.WriteAll(a.rows)
	if err != nil {
		return err
	}
	return f.Close()
}
