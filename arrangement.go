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
