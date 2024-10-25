package alias

import (
	"crypto/rand"
	"log/slog"
	"math/big"

	"github.com/ilius/localsend-go/pkg/alias/en"
	"github.com/ilius/localsend-go/pkg/alias/fa"
)

type LangInterface interface {
	Adjectives() []string
	Fruits() []string
	Combination(adjective string, fruit string) string
}

var (
	_en LangInterface = en.New()
	_fa LangInterface = fa.New()
)

func genAlias(lang LangInterface) string {
	adjectives := lang.Adjectives()
	fruits := lang.Fruits()
	return lang.Combination(adjectives[randomInt(len(adjectives))], fruits[randomInt(len(fruits))])
}

func GenerateRandomAlias(lang string) string {
	switch lang {
	case "", "en":
		return genAlias(_en)
	case "fa":
		return genAlias(_fa)
	default:
		slog.Error("unsupported language name", "lang", lang)
	}
	return genAlias(_en)
}

// randomInt returns a uniform random value in [0, max). It panics if max <= 0.
func randomInt(max int) int {
	ibig, err := rand.Int(rand.Reader, big.NewInt(int64(max)))
	if err != nil {
		panic(err) // not sure how to trigger this in test
	}
	return int(ibig.Int64())
}