package track

import (
	"crypto/rand"

	"github.com/btcsuite/btcutil/base58"
)

func CreateIfNotExists(trk *Track, part *Part) error {
	for _, p := range trk.parts {
		if p.Name == part.Name {
			return nil
		}
	}
	err := trk.db.Save(part)
	if err != nil {
		return err
	}
	trk.parts = append(trk.parts, part)
	sortPartSlice(trk.parts)
	return nil
}

func SetPart(trk *Track, part *Part) error {
	err := trk.db.Save(part)
	if err != nil {
		return err
	}
	for i, p := range trk.parts {
		if p.Name == part.Name {
			trk.parts[i] = part
			return nil
		}
	}
	trk.parts = append(trk.parts, part)
	sortPartSlice(trk.parts)
	return nil
}

func SetPattern(trk *Track, part *Part, tick, pattern int) error {
	var ch *PatternChange
	for _, pc := range trk.pc {
		if pc.Part == part.Name && pc.Tick == tick {
			ch = pc
			break
		}
	}
	var id string
	if ch != nil {
		id = ch.ID
	} else {
		id = makeID()
	}
	err := trk.db.Save(&PatternChange{id, part.Name, tick, pattern})
	if err != nil {
		return err
	}
	if ch != nil {
		ch.Pattern = pattern
	} else {
		trk.pc = append(trk.pc, &PatternChange{id, part.Name, tick, pattern})
		sortPatternChangeSlice(trk.pc)
	}
	return nil
}

func SetMute(trk *Track, part *Part, tick int, mute [16]bool) error {
	var ch *MuteChange
	for _, mc := range trk.mc {
		if mc.Part == part.Name && mc.Tick == tick {
			ch = mc
			break
		}
	}
	var id string
	if ch != nil {
		id = ch.ID
	} else {
		id = makeID()
	}
	err := trk.db.Save(&MuteChange{id, part.Name, tick, mute})
	if err != nil {
		return err
	}
	if ch != nil {
		ch.Mute = mute
	} else {
		trk.mc = append(trk.mc, &MuteChange{id, part.Name, tick, mute})
		sortMuteChangeSlice(trk.mc)
	}
	return nil
}

func SetMuted(trk *Track, part *Part, tick int, track int, muted bool) error {
	mute := Mute(trk, part, tick)
	mute[track] = muted
	return SetMute(trk, part, tick, mute)
}

func Clear(trk *Track, tick int) (err error) {
	var i int
	for _, pc := range trk.pc {
		if pc.Tick == tick {
			err = trk.db.DeleteStruct(pc)
			if err != nil {
				return
			}
		} else {
			trk.pc[i] = pc
			i++
		}
	}
	trk.pc = trk.pc[:i]

	i = 0
	for _, mc := range trk.mc {
		if mc.Tick == tick {
			err = trk.db.DeleteStruct(mc)
			if err != nil {
				return
			}
		} else {
			trk.mc[i] = mc
			i++
		}
	}
	trk.mc = trk.mc[:i]

	return
}

func makeID() string {
	b := make([]byte, 12)
	_, err := rand.Read(b)
	if err != nil {
		panic(err)
	}
	return base58.Encode(b)
}
