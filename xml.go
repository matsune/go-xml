package xml

func newParser(str string) *parser {
	return &parser{
		scanner: &scanner{
			source: []rune(str),
			cursor: 0,
		},
	}
}

func Parse(str string) (*XML, error) {
	return newParser(str).parse()
}
