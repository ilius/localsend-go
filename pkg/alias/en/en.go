package en

func New() *impl {
	return &impl{}
}

type impl struct {
}

func (*impl) Adjectives() []string {
	return adjectives
}

func (*impl) Fruits() []string {
	return fruits
}

func (*impl) Combination(adjective string, fruit string) string {
	return adjective + " " + fruit
}
