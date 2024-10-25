package alias

import (
	"crypto/rand"
	"math/big"

	"github.com/ilius/localsend-go/pkg/alias/en"
)

// TODO: add language arg
func GenerateRandomAlias() string {
	return en.Adjectives[randomInt(len(en.Adjectives))] + " " + en.Fruits[randomInt(len(en.Fruits))]
}

// randomInt returns a uniform random value in [0, max). It panics if max <= 0.
func randomInt(max int) int {
	ibig, err := rand.Int(rand.Reader, big.NewInt(int64(max)))
	if err != nil {
		panic(err) // not sure how to trigger this in test
	}
	return int(ibig.Int64())
}
