package main

import "fmt"

type Track struct {
	pattern    map[int]int
	muted      map[mutedKey]bool
	voiceCount int
}

func NewTrack(voiceCount int) *Track {
	return &Track{make(map[int]int), make(map[mutedKey]bool), voiceCount}
}

func (t *Track) VoiceCount() int {
	return t.voiceCount
}

func (t *Track) SetPattern(step int, pattern int) {
	t.pattern[step] = pattern
}

func (t *Track) Pattern(step int) (pattern int, change bool) {
	pattern, change = t.pattern[step]
	if change {
		return
	}
	bestKey, pattern := 0, t.pattern[0]
	for key, val := range t.pattern {
		if key < step && key > bestKey {
			bestKey = key
			pattern = val
		}
	}
	if bestKey == 0 && step == 0 {
		change = true
	}
	return
}

func (t *Track) SetMuted(step int, voice int, muted bool) {
	t.muted[mutedKey{step, voice}] = muted
}

func (t *Track) Muted(step, voice int) (muted bool, change bool) {
	muted, change = t.muted[mutedKey{step, voice}]
	if change {
		return
	}
	bestKey, muted := mutedKey{0, voice}, t.muted[mutedKey{0, voice}]
	for key, val := range t.muted {
		if key.voice == voice && key.step < step && key.step > bestKey.step {
			bestKey = key
			muted = val
		}
	}
	if bestKey.step == 0 && step == 0 {
		change = true
	}
	return
}

type mutedKey struct {
	step, voice int
}

type Pattern int

func (p Pattern) String() string {
	return fmt.Sprintf("%s%02d", string('A'+int(p)/16), 1+int(p)%16)
}
