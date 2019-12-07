package main

import (
	"bufio"
	"bytes"
	"log"
	"os"
	"sort"
	"sync"
	"time"

	"github.com/nsf/termbox-go"
)

type Console struct {
	Enabled       bool
	width, height int
	lines         []string
	mx            sync.Mutex
}

func NewConsole() *Console {
	os.Remove("trk.log")
	c := &Console{width: 80, height: 10}
	log.SetOutput(c)
	return c
}

func (c *Console) Write(p []byte) (n int, err error) {
	c.writeToFile(p)
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

func (c *Console) writeToFile(b []byte) (err error) {
	f, err := os.OpenFile("trk.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		return
	}
	defer f.Close()
	_, err = f.Write(b)
	if err != nil {
		return
	}
	err = f.Sync()
	if err != nil {
		return
	}
	return f.Close()
}

type TimeTracker struct {
	start    time.Time
	measures map[string][]time.Duration
}

func NewTimeTracker() *TimeTracker {
	return &TimeTracker{time.Now(), make(map[string][]time.Duration)}
}

func (tt *TimeTracker) Trace(label string) {
	tt.measures[label] = append(tt.measures[label], time.Since(tt.start))
	tt.start = time.Now()
}

func (tt *TimeTracker) Log() {
	var labels []string
	var global time.Duration
	var total = make(map[string]time.Duration)
	var count = make(map[string]time.Duration)
	for label, measures := range tt.measures {
		for _, measure := range measures {
			global += measure
			total[label] += measure
			count[label]++
		}
		labels = append(labels, label)
	}
	sort.Strings(labels)
	for _, label := range labels {
		if count[label] == 0 {
			continue
		}
		log.Printf("[%s] %s x %d", label, total[label]/count[label], count[label])
	}
	log.Println("===", global)
}
