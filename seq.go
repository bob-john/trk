package main

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"reflect"
	"strconv"
	"strings"

	"github.com/gomidi/midi"
	"github.com/gomidi/midi/midimessage/channel"
)

type Seq struct {
	row map[int]*Row
}

func NewSeq() *Seq {
	return &Seq{make(map[int]*Row)}
}

func ReadSeq(f io.Reader) (*Seq, error) {
	seq := NewSeq()
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
	return seq, nil
}

func (s *Seq) Write(f io.Writer) error {
	w := csv.NewWriter(f)
	for _, row := range s.row {
		err := w.Write(row.Record())
		if err != nil {
			return err
		}
	}
	w.Flush()
	return w.Error()
}

func (s *Seq) Insert(device string, row int, message midi.Message) {
	r := s.Row(row).Copy()
	switch m := message.(type) {
	case channel.ProgramChange:
		if strings.Contains(device, "Digitone") {
			r.Parts[Digitone].Pattern = Pattern(m.Program())
		}
		if strings.Contains(device, "Digitakt") {
			r.Parts[Digitakt].Pattern = Pattern(m.Program())
		}

	case channel.ControlChange:
		if m.Controller() != 94 {
			return
		}
		ch := int(m.Channel())
		if r.Parts[Digitone].Channels.Contains(ch) {
			r.Parts[Digitone].Mute = s.ConsolidatedRow(row).Parts[Digitone].Mute.Copy()
			r.Parts[Digitone].Mute[r.Parts[Digitone].Channels.IndexOf(ch)] = m.Value() != 0
		}
		if r.Parts[Digitakt].Channels.Contains(ch) {
			r.Parts[Digitakt].Mute = s.ConsolidatedRow(row).Parts[Digitakt].Mute.Copy()
			r.Parts[Digitakt].Mute[r.Parts[Digitakt].Channels.IndexOf(ch)] = m.Value() != 0
		}
	}
	if r.HasChanges(s.ConsolidatedRow(row)) {
		s.row[row] = r
	}
}

func (s *Seq) Row(index int) *Row {
	r, ok := s.row[index]
	if ok {
		return r
	}
	r = NewRow(index)
	if index <= 0 {
		r.SetDefaults()
	}
	return r
}

func (s *Seq) ConsolidatedRow(index int) *Row {
	row := s.Row(index).Copy()
	for index > 0 && !row.Consolidated() {
		index--
		row.Merge(s.Row(index))
	}
	return row
}

func (s *Seq) Clear(row int) {
	s.row[row] = NewRow(row)
}

func (s *Seq) Text(row int) string {
	if row%16 == 0 {
		return s.ConsolidatedRow(row).String()
	} else {
		return s.Row(row).String()
	}
}

type Row struct {
	Index int
	Parts map[string]*Part
}

func DecodeRow(record []string) (*Row, error) {
	if len(record) != 5 {
		return nil, errors.New("row: invalid record")
	}
	index, _ := strconv.Atoi(record[0])
	res := NewRow(index)
	res.Parts[Digitone] = DecodePart(record[1:3], MakeRange(8, 11))
	res.Parts[Digitakt] = DecodePart(record[3:5], MakeRange(0, 7))
	return res, nil
}

func NewRow(index int) *Row {
	return &Row{
		Index: index,
		Parts: map[string]*Part{
			Digitone: NewPart(MakeRange(8, 11)),
			Digitakt: NewPart(MakeRange(0, 7)),
		},
	}
}

func (r *Row) Record() []string {
	return []string{
		strconv.Itoa(r.Index),
		r.Parts[Digitone].Pattern.String(),
		r.Parts[Digitone].Mute.String(),
		r.Parts[Digitakt].Pattern.String(),
		r.Parts[Digitakt].Mute.String(),
	}
}

func (r *Row) String() string {
	return fmt.Sprintf("%3d %s %s", 1+r.Index%16, r.Parts[Digitone], r.Parts[Digitakt])
}

func (r *Row) SetDefaults() {
	r.Parts[Digitone].SetDefaults()
	r.Parts[Digitakt].SetDefaults()
}

func (r *Row) Consolidated() bool {
	return r.Parts[Digitone].Consolidated() && r.Parts[Digitakt].Consolidated()
}

func (r *Row) Merge(o *Row) {
	r.Parts[Digitone].Merge(o.Parts[Digitone])
	r.Parts[Digitakt].Merge(o.Parts[Digitakt])
}

func (r *Row) HasChanges(prev *Row) bool {
	if r.Parts[Digitone].HasChanges(prev.Parts[Digitone]) {
		return true
	}
	if r.Parts[Digitakt].HasChanges(prev.Parts[Digitakt]) {
		return true
	}
	return false
}

func (r *Row) Copy() *Row {
	dup := &Row{Index: r.Index, Parts: make(map[string]*Part)}
	for name, part := range r.Parts {
		dup.Parts[name] = part.Copy()
	}
	return dup
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

func (p *Part) HasChanges(prev *Part) bool {
	if p.Pattern != -1 && p.Pattern != prev.Pattern {
		return true
	}
	if len(p.Mute) != 0 && !reflect.DeepEqual(p.Mute, prev.Mute) {
		return true
	}
	return false
}

func (p *Part) Copy() *Part {
	return &Part{
		Pattern:  p.Pattern,
		Mute:     p.Mute.Copy(),
		Channels: p.Channels,
	}
}

type Pattern int

func DecodePattern(field string) Pattern {
	if len(field) != 3 {
		return -1
	}
	bank := int(field[0] - 'A')
	trig, err := strconv.Atoi(field[1:])
	if bank >= 0 && bank < 8 && trig >= 1 && trig <= 16 && err == nil {
		return Pattern(16*bank + trig - 1)
	}
	return -1
}

func (p Pattern) String() string {
	if p == -1 {
		return "..."
	}
	return fmt.Sprintf("%s%02d", string('A'+int(p)/16), 1+int(p)%16)
}

type Mute map[int]bool

func DecodeMute(field string) (res Mute) {
	res = make(map[int]bool)
	for _, c := range field {
		n := int(c - '1')
		if n >= 0 && n < 8 {
			res[n] = false
		}
	}
	if len(res) == 0 {
		return
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

func (m Mute) Copy() Mute {
	cpy := make(Mute)
	for k, v := range m {
		cpy[k] = v
	}
	return cpy
}
