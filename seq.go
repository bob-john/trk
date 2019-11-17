package main

import (
	"archive/zip"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/gomidi/midi"
	"github.com/gomidi/midi/midimessage/channel"
)

type Seq struct {
	row map[int]*Row
}

func ReadSeq(path string) (*Seq, error) {
	seq := NewSeq()
	r, err := zip.OpenReader(path)
	if err != nil {
		return nil, err
	}
	defer r.Close()
	for _, f := range r.File {
		switch f.Name {
		case "seq.csv":
			f, err := f.Open()
			if err != nil {
				return nil, err
			}
			defer f.Close()
			r := csv.NewReader(f)
			r.ReuseRecord = true
			for {
				record, err := r.Read()
				if err == io.EOF {
					break
				}
				if err != nil {
					return nil, err
				}
				row, err := DecodeRow(record)
				if err != nil {
					return nil, err
				}
				seq.row[row.Index] = row
			}
			f.Close()
		}
	}
	return seq, nil
}

func NewSeq() *Seq {
	return &Seq{make(map[int]*Row)}
}

func (s *Seq) Write(path string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	w := zip.NewWriter(f)
	defer w.Close()
	o, err := w.Create("seq.csv")
	if err != nil {
		return err
	}
	u := csv.NewWriter(o)
	for _, row := range s.row {
		err = u.Write(row.Record())
		if err != nil {
			return err
		}
	}
	u.Flush()
	err = u.Error()
	if err != nil {
		return err
	}
	err = w.Close()
	if err != nil {
		return err
	}
	return f.Close()
}

func (s *Seq) Insert(device string, row int, message midi.Message) {
	switch m := message.(type) {
	case channel.ProgramChange:
		r := s.Row(row)
		if strings.Contains(device, "Digitone") {
			r.Digitone.Pattern = Pattern(m.Program())
			s.row[row] = r
		}
		if strings.Contains(device, "Digitakt") {
			r.Digitakt.Pattern = Pattern(m.Program())
			s.row[row] = r
		}

	case channel.ControlChange:
		if m.Controller() != 94 {
			return
		}
		r := s.Row(row)
		ch := int(m.Channel())
		if r.Digitone.Channels.Contains(ch) {
			r.Digitone.Mute[r.Digitone.Channels.IndexOf(ch)] = m.Value() != 0
			s.row[row] = r
		}
		if r.Digitakt.Channels.Contains(ch) {
			r.Digitakt.Mute[r.Digitakt.Channels.IndexOf(ch)] = m.Value() != 0
			s.row[row] = r
		}
	}
}

func (s *Seq) Row(index int) *Row {
	r, ok := s.row[index]
	if ok {
		return r
	}
	r = NewRow(index)
	if index == 0 {
		r.SetDefaults()
	}
	return r
}

func (s *Seq) ConsolidatedRow(index int) *Row {
	row := s.Row(index)
	for index > 0 && !row.Consolidated() {
		index--
		row.Merge(s.Row(index))
	}
	return row
}

func (s *Seq) Text(row int) string {
	if row%16 == 0 {
		return s.ConsolidatedRow(row).String()
	} else {
		return s.Row(row).String()
	}
}

type Row struct {
	Index    int
	Digitone *Part
	Digitakt *Part
}

func DecodeRow(record []string) (*Row, error) {
	if len(record) != 5 {
		return nil, errors.New("row: invalid record")
	}
	res := new(Row)
	res.Index, _ = strconv.Atoi(record[0])
	res.Digitone = DecodePart(record[1:3], MakeRange(8, 11))
	res.Digitakt = DecodePart(record[3:5], MakeRange(0, 7))
	return res, nil
}

func NewRow(index int) *Row {
	return &Row{
		Index:    index,
		Digitone: NewPart(MakeRange(8, 11)),
		Digitakt: NewPart(MakeRange(0, 7)),
	}
}

func (r *Row) Record() []string {
	return []string{
		strconv.Itoa(r.Index),
		r.Digitone.Pattern.String(),
		r.Digitone.Mute.String(),
		r.Digitakt.Pattern.String(),
		r.Digitakt.Mute.String(),
	}
}

func (r *Row) String() string {
	return fmt.Sprintf("%3d %s %s", 1+r.Index%16, r.Digitone, r.Digitakt)
}

func (r *Row) SetDefaults() {
	r.Digitone.SetDefaults()
	r.Digitakt.SetDefaults()
}

func (r *Row) Consolidated() bool {
	return r.Digitone.Consolidated() && r.Digitakt.Consolidated()
}

func (r *Row) Merge(o *Row) {
	r.Digitone.Merge(o.Digitone)
	r.Digitakt.Merge(o.Digitakt)
}

func (r *Row) Play(digitone, digitakt *Device) {
	r.Digitone.Play(digitone)
	r.Digitakt.Play(digitakt)
}

type Part struct {
	Pattern  Pattern
	Mute     Mute
	Channels Range
}

func NewPart(channels Range) *Part {
	return &Part{
		Pattern:  -1,
		Mute:     make(Mute),
		Channels: channels,
	}
}

func DecodePart(rec []string, channels Range) *Part {
	if len(rec) != 2 {
		return nil
	}
	return &Part{
		Pattern:  DecodePattern(rec[0]),
		Mute:     DecodeMute(rec[1]),
		Channels: channels,
	}
}

func (p *Part) SetDefaults() {
	p.Pattern = 0
	for n := 0; n < p.Channels.Len; n++ {
		p.Mute[n] = false
	}
}

func (p *Part) String() string {
	return fmt.Sprintf("%s %s", p.Pattern, p.Mute.Format(p.Channels))
}

func (p *Part) Consolidated() bool {
	return p.Pattern != -1 && len(p.Mute) != 0
}

func (p *Part) Merge(o *Part) {
	if p.Pattern == -1 {
		p.Pattern = o.Pattern
	}
	if len(p.Mute) == 0 {
		for n, m := range o.Mute {
			p.Mute[n] = m
		}
	}
}

func (p *Part) Play(out *Device) {
	p.Pattern.Play(out, 15)
	p.Mute.Play(out, p.Channels)
}

type Pattern int

func DecodePattern(field string) Pattern {
	if len(field) != 3 {
		return -1
	}
	bank := int(field[0] - 'A')
	pattern, err := strconv.Atoi(field[1:])
	if bank >= 0 && bank < 8 && pattern >= 0 && pattern < 16 && err == nil {
		return Pattern(16*bank + pattern)
	}
	return -1
}

func (p Pattern) String() string {
	if p == -1 {
		return "..."
	}
	return fmt.Sprintf("%s%02d", string('A'+int(p)/16), 1+int(p)%16)
}

func (p Pattern) Play(out *Device, ch int) {
	out.Write(channel.Channel(ch).ProgramChange(uint8(p)))
}

type Mute map[int]bool

func DecodeMute(field string) (res Mute) {
	res = make(map[int]bool)
	for _, c := range field {
		if c-'1' >= 0 && c-'1' < 8 {
			res[int(c-'1')] = false
		} else if c-'A' >= 0 && c-'A' < 8 {
			res[int(8+c-'A')] = false
		}
	}
	for n := 0; n < 16; n++ {
		if _, ok := res[n]; !ok {
			res[n] = true
		}
	}
	return
}

func (m Mute) Format(channels Range) string {
	if len(m) == 0 {
		return strings.Repeat(".", channels.Len)
	}
	var str string
	for n := 0; n < channels.Len; n++ {
		if m[n] {
			str += "-"
		} else if n < 8 {
			str += string('1' + n)
		}
	}
	return str
}

func (m Mute) String() string {
	return m.Format(MakeRange(0, 15))
}

func (m Mute) Play(out *Device, channels Range) {
	for k, v := range m {
		ch := channels.Index + k
		var muted uint8
		if v {
			muted = 1
		}
		out.Write(channel.Channel(ch).ControlChange(94, muted))
	}
}
