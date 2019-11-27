package track

import (
	"crypto/rand"

	"github.com/asdine/storm"
	"github.com/asdine/storm/q"
	"github.com/btcsuite/btcutil/base58"
)

func CreateIfNotExists(trk *Track, part *Part) error {
	var val Part
	err := trk.db.One("Name", part.Name, &val)
	if err == storm.ErrNotFound {
		return trk.db.Save(part)
	}
	return err
}

func SetPart(trk *Track, part *Part) error {
	return trk.db.Save(part)
}

func SetPattern(trk *Track, part *Part, tick, pattern int) error {
	var ch PatternChange
	err := trk.db.Select(q.Eq("Part", part.Name), q.Eq("Tick", tick)).First(&ch)
	if err == nil {
		ch.Pattern = pattern
	} else if err == storm.ErrNotFound {
		ch = PatternChange{makeID(), part.Name, tick, pattern}
	} else {
		return err
	}
	return trk.db.Save(&ch)
}

func SetMute(trk *Track, part *Part, tick int, mute [16]bool) error {
	var ch MuteChange
	err := trk.db.Select(q.Eq("Part", part.Name), q.Eq("Tick", tick)).First(&ch)
	if err == nil {
		ch.Mute = mute
	} else if err == storm.ErrNotFound {
		ch = MuteChange{makeID(), part.Name, tick, mute}
	} else {
		return err
	}
	return trk.db.Save(&ch)
}

func SetMuted(trk *Track, part *Part, tick int, track int, muted bool) error {
	mute, _ := Mute(trk, part, tick)
	mute[track] = muted
	return SetMute(trk, part, tick, mute)
}

func makeID() string {
	b := make([]byte, 12)
	_, err := rand.Read(b)
	if err != nil {
		panic(err)
	}
	return base58.Encode(b)
}
