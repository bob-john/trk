package main

type Mute []bool

func ParseMute(str string, channelCount int) Mute {
	m := make(Mute, channelCount)
	for _, ch := range str {
		n := int(ch) - '1'
		if n < 0 || n >= len(m) {
			continue
		}
		m[n] = true
	}
	for n, v := range m {
		m[n] = !v
	}
	return m
}

func (m Mute) String() string {
	str := make([]rune, len(m))
	for n, v := range m {
		if v {
			str[n] = '-'
		} else {
			str[n] = '1' + rune(n)
		}
	}
	return string(str)
}

func (m Mute) Clear() {
	for n := range m {
		m[n] = true
	}
}
