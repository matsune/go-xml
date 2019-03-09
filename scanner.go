package xml

type scanner struct {
	source []rune
	cursor uint
}

func (s *scanner) pos() Pos {
	c := int(s.cursor)
	if len(s.source) < c {
		c = len(s.source)
	}

	sub := s.source[0:c]
	p := Pos{
		Line: 1,
		Col:  1,
	}
	for _, v := range sub {
		if v == '\n' {
			p.Line++
			p.Col = 1
		} else {
			p.Col++
		}
	}
	return p
}

func (s *scanner) isEnd() bool {
	return len(s.source) <= int(s.cursor)
}

func (s *scanner) Get() rune {
	if s.isEnd() {
		return 0
	}
	return rune(s.source[s.cursor])
}

func (s *scanner) Step() {
	s.cursor++
}

func (s *scanner) StepN(n int) {
	i := int(s.cursor)
	i += n
	s.cursor = uint(i)
}
