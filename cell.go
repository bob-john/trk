package main

type Cell interface {
	Edit() CellEditor
	Set(string)
	String() string
	Output(*Device)
}

type stringCell struct {
	doc      *Arrangement
	row, col int
}

func (c stringCell) Set(val string) {
	c.doc.Set(c.row, c.col, val)
}

func (c stringCell) String() string {
	return c.doc.Get(c.row, c.col)
}
