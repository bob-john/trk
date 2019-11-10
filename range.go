package main

type Range struct {
	Index, Len int
}

func (r Range) Substr(str string) string {
	return str[r.Index : r.Index+r.Len]
}

func (r Range) Replace(str, repl string) string {
	return str[:r.Index] + repl + str[r.Index+r.Len:]
}
