package track

import (
	"crypto/rand"

	"github.com/btcsuite/btcutil/base58"
)

func makeID() string {
	b := make([]byte, 12)
	_, err := rand.Read(b)
	if err != nil {
		panic(err)
	}
	return base58.Encode(b)
}
