package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"sort"
	"strconv"
	"strings"
)

type Seq struct {
	lines map[int]string
}

func (s *Seq) Line(step int) string {
	line, ok := s.lines[step]
	if ok {
		return line
	}
	return fmt.Sprintf("%03X ... ........ ... ....", step)
}

func (s *Seq) ConsolidatedLine(step int) string {
	return fmt.Sprintf("%03X A01 ++++++++ A01 ++++", step)
}

func (s *Seq) Insert(line string) {
	step := s.parseStep(line)
	if strings.ContainsAny(line[4:], "ABCDEFGH0123456789") {
		if s.lines == nil {
			s.lines = make(map[int]string)
		}
		s.lines[step] = line
	} else {
		delete(s.lines, step)
	}
}

func (s *Seq) Write(w io.Writer) (n int, err error) {
	var steps []int
	for _, line := range s.lines {
		steps = append(steps, s.parseStep(line))
	}
	sort.Ints(steps)
	for _, step := range steps {
		var d int
		d, err = io.WriteString(w, s.lines[step]+"\n")
		n += d
		if err != nil {
			return
		}
	}
	return
}

func (s *Seq) WriteFile(path string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = s.Write(f)
	if err != nil {
		return err
	}
	return f.Close()
}

func (s *Seq) Read(r io.Reader) error {
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		s.Insert(scanner.Text())
	}
	if err := scanner.Err(); err != io.EOF {
		return err
	}
	return nil
}

func (s *Seq) ReadFile(path string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()
	return seq.Read(f)
}

func (s *Seq) parseStep(line string) int {
	step, _ := strconv.ParseInt(line[:3], 16, 32)
	return int(step)
}
