package xml

import "fmt"

type Scanner struct {
	source []rune
	cursor int
}

func NewScanner(str string) *Scanner {
	return &Scanner{
		source: []rune(str),
		cursor: 0,
	}
}

func (s *Scanner) errorf(f string, args ...interface{}) error {
	h := fmt.Sprintf("cursor: %d ", s.cursor)
	return fmt.Errorf(h+f, args...)
}

func (s *Scanner) Get() rune {
	if len(s.source) <= s.cursor {
		return 0
	}
	return rune(s.source[s.cursor])
}

func (s *Scanner) Test(r rune) bool {
	return s.Get() == r
}

func (s *Scanner) Must(r rune) error {
	if !s.Test(r) {
		return s.errorf("expected %q", r)
	}
	s.Step()
	return nil
}

func (s *Scanner) Tests(str string) bool {
	i := s.cursor
	e := i + len([]rune(str))
	if len(s.source) < e {
		return false
	}
	return string(s.source[i:e]) == str
}

func (s *Scanner) Musts(str string) error {
	if !s.Tests(str) {
		return s.errorf("expected %q", str)
	}
	s.StepN(len(str))
	return nil
}

func (s *Scanner) Step() {
	s.cursor++
}

func (s *Scanner) StepN(n int) {
	s.cursor += n
}
