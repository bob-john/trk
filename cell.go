package main

// type Cell interface {
// 	Range() Range
// 	Inc()
// 	Dec()
// 	LargeInc()
// 	LargeDec()
// 	Clear()
// }

// type PatternCell struct {
// 	StringCell
// }

// func NewPatternCell(line *string, index, len int) Cell {
// 	return PatternCell{MakeStringCell(line, index, len)}
// }

// func (c PatternCell) Inc() {
// 	c.update(1)
// }

// func (c PatternCell) Dec() {
// 	c.update(-1)
// }

// func (c PatternCell) LargeInc() {
// 	c.update(16)
// }

// func (c PatternCell) LargeDec() {
// 	c.update(-16)
// }

// func (c PatternCell) Clear() {
// 	c.Set("...")
// }

// func (c PatternCell) set(val int) {
// 	c.Set(c.encode(clamp(val, 0, 127)))
// }

// func (c PatternCell) update(delta int) {
// 	c.set(c.value() + delta)
// }

// func (c PatternCell) value() int {
// 	str := c.String()
// 	if strings.ContainsAny(str, ".") {
// 		return -1
// 	}
// 	bank := int(str[0] - 'A')
// 	trig, _ := strconv.Atoi(str[1:])
// 	return bank*16 + trig - 1
// }

// func (c PatternCell) encode(val int) string {
// 	if val < 0 {
// 		return "..."
// 	}
// 	return fmt.Sprintf("%s%02d", string('A'+val/16), 1+val%16)
// }

// type MuteCell struct {
// 	StringCell
// }

// func NewMuteCell(line *string, index, len int) Cell {
// 	return MuteCell{MakeStringCell(line, index, len)}
// }

// func (c MuteCell) Inc() {
// 	c.Set("+")
// }

// func (c MuteCell) Dec() {
// 	c.Set("-")
// }

// func (c MuteCell) LargeInc() {
// 	c.Set("+")
// }

// func (c MuteCell) LargeDec() {
// 	c.Set("-")
// }

// func (c MuteCell) Clear() {
// 	c.Set(".")
// }

// type LenCell struct {
// 	StringCell
// }

// func NewLenCell(line *string, index, len int) Cell {
// 	return LenCell{MakeStringCell(line, index, len)}
// }

// func (c LenCell) Inc() {
// 	c.update(1)
// }

// func (c LenCell) Dec() {
// 	c.update(-1)
// }

// func (c LenCell) LargeInc() {
// 	c.set(snap(c.value()+16, 16))
// }

// func (c LenCell) LargeDec() {
// 	c.set(snap(c.value()-16, 16))
// }

// func (c LenCell) Clear() {
// 	c.Set("  16")
// }

// func (c LenCell) update(delta int) {
// 	c.set(c.value() + delta)
// }

// func (c LenCell) set(val int) {
// 	c.Set(c.encode(clamp(val, 1, 1024)))
// }

// func (c LenCell) value() int {
// 	val, _ := strconv.Atoi(strings.TrimSpace(c.String()))
// 	return val
// }

// func (c LenCell) encode(val int) string {
// 	return fmt.Sprintf("%4d", val)
// }

// type StringCell struct {
// 	line *string
// 	old  string
// 	rng  Range
// }

// func MakeStringCell(line *string, index, len int) StringCell {
// 	return StringCell{line, *line, Range{index, len}}
// }

// func (c StringCell) Range() Range {
// 	return c.rng
// }

// func (c StringCell) String() string {
// 	return c.rng.Substr(*c.line)
// }

// func (c StringCell) Set(val string) {
// 	*c.line = c.rng.Replace(*c.line, val)
// }
