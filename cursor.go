package main

type Cursor struct {
	lp       *Launchpad
	color    uint8
	row, col uint8
	flashing bool
}

func NewCursor(lp *Launchpad, color uint8) *Cursor {
	lp.Set(0, 0, color)
	return &Cursor{lp, color, 0, 0, false}
}

func (c *Cursor) MoveTo(row, col uint8) {
	c.lp.Set(c.row, c.col, 0)
	c.row = row
	c.col = col
	if c.flashing {
		c.lp.Set(c.row, c.col, c.color)
		c.lp.Update()
		c.lp.StartFlashing(c.row, c.col, 0)
	} else {
		c.lp.Set(c.row, c.col, c.color)
	}
}

func (c *Cursor) IsAt(row, col uint8) bool {
	return c.row == row && c.col == col
}

func (c *Cursor) ToggleFlashing() {
	if c.flashing {
		c.StopFlashing()
	} else {
		c.StartFlashing()
	}
}

func (c *Cursor) IsFlashing() bool {
	return c.flashing
}

func (c *Cursor) StartFlashing() {
	c.flashing = true
	c.lp.StartFlashing(c.row, c.col, 0)
}

func (c *Cursor) StopFlashing() {
	c.lp.StopFlashing(c.row, c.col)
	c.flashing = false
}
