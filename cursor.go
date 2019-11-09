package main

var cursorFlashingColor uint8 = 0

type Cursor struct {
	row, col       uint8
	color          uint8
	flashing       bool
	oldRow, oldCol uint8
	moved          bool
}

func NewCursor(color uint8) *Cursor {
	return &Cursor{0, 0, color, false, 0, 0, false}
}

func (c *Cursor) Render(lp *Launchpad) {
	if c.moved {
		lp.Draw(c.oldRow, c.oldCol, 0)
		c.moved = false
	}
	lp.Draw(c.row, c.col, c.color)
	if c.flashing {
		lp.SetFlashing(c.row, c.col, &cursorFlashingColor)
	} else {
		lp.SetFlashing(c.row, c.col, nil)
	}
}

func (c *Cursor) Set(row, col uint8) {
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

func (c *Cursor) IsAt(row, col uint8) bool {
	return c.row == row && c.col == col
}

func (c *Cursor) SetFlashing(flashing bool) {
	c.flashing = flashing
}

func (c *Cursor) IsFlashing() bool {
	return c.flashing
}
