package main

type Synth struct {
	pattern    map[int]Pattern
	muted      map[mutedKey]bool
	voiceCount int
}

func NewSynth(voiceCount int) *Synth {
	return &Synth{make(map[int]Pattern), make(map[mutedKey]bool), voiceCount}
}

func (t *Synth) VoiceCount() int {
	return t.voiceCount
}

func (t *Synth) SetPattern(step int, pattern Pattern) {
	t.pattern[step] = pattern
}

func (t *Synth) Pattern(step int) (pattern Pattern, change bool) {
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

func (t *Synth) SetMuted(step int, voice int, muted bool) {
	t.muted[mutedKey{step, voice}] = muted
}

func (t *Synth) Muted(step, voice int) (muted bool, change bool) {
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
