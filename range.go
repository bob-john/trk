package main

type Range struct {
	Index, Len int
}

func MakeRange(min, max int) Range {
	if min < max {
		return Range{min, max - min + 1}
	}
	return Range{max, min - max + 1}
}

func (r Range) Substr(str string) string {
	return str[r.Index : r.Index+r.Len]
}

func (r Range) Replace(str, repl string) string {
	return str[:r.Index] + repl + str[r.Index+r.Len:]
}

func (r Range) Contains(n int) bool {
	i := n - r.Index
	return i >= 0 && i < r.Len
}

func (r Range) IndexOf(n int) int {
	return n - r.Index
}
