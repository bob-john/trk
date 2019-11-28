package main

import (
	"bufio"
	"bytes"
	"log"

	"github.com/nsf/termbox-go"
)

type Console struct {
	Enabled       bool
	width, height int
	lines         []string
}

func NewConsole() *Console {
	c := &Console{width: 80, height: 10}
	log.SetOutput(c)
	return c
}

func (c *Console) Write(p []byte) (n int, err error) {
	scanner := bufio.NewScanner(bytes.NewReader(p))
	for scanner.Scan() {
		line := scanner.Text()
		if len(line) > c.width {
			line = line[:c.width-3] + "..."
		}
		c.lines = append(c.lines, line)
	}
	if len(c.lines) > c.height {
		c.lines = c.lines[len(c.lines)-c.height:]
	}
	return len(p), scanner.Err()
}

func (c *Console) Render() {
	if !c.Enabled {
		return
	}
	DrawBox(0, 7, c.width+1, 7+c.height+1, termbox.ColorDefault, termbox.ColorDefault)
	DrawString(1, 7, " DEBUG CONSOLE ", termbox.ColorDefault, termbox.ColorDefault)
	for n, line := range c.lines {
		DrawString(1, 8+n, line, termbox.ColorDefault, termbox.ColorDefault)
	}
}
