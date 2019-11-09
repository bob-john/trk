package main

type Cursor struct {
	row, col       int
	color          int
	flashingColor  int
	flashing       bool
	oldRow, oldCol int
	moved          bool
}

func NewCursor(color, flashingColor int) *Cursor {
	return &Cursor{0, 0, color, flashingColor, false, 0, 0, false}
}

func (c *Cursor) Render(lp *Launchpad) {
	if c.moved {
		lp.Draw(c.oldRow, c.oldCol, 0)
		c.moved = false
	}
	lp.Draw(c.row, c.col, c.color)
	if c.flashing {
		lp.SetFlashing(c.row, c.col, &c.flashingColor)
	} else {
		lp.SetFlashing(c.row, c.col, nil)
	}
}

func (c *Cursor) SetColor(color int) {
	c.color = color
}

func (c *Cursor) Move(row, col int) {
	if c.row == row && c.col == col {
		return
	}
	if !c.moved {
		c.oldRow, c.oldCol = c.row, c.col
		c.moved = true
	}
	c.row = row
	c.col = col
}

func (c *Cursor) IsAt(row, col int) bool {
	return c.row == row && c.col == col
}

func (c *Cursor) SetFlashing(flashing bool) {
	c.flashing = flashing
}
