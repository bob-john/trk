package main

import "fmt"

type PC int

func (pc PC) String() string {
	return fmt.Sprintf("%s%02d", string('A'+pc/16), 1+pc%16)
}

type Muted []bool

func (m Muted) String() string {
	var s string
	for _, m := range m {
		if m {
			s += "-" //"\u258e"
		} else {
			s += "+"
		}
	}
	return s
}
