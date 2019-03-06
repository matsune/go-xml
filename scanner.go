package xml

import "fmt"

type scanner struct {
	source []rune
	cursor int
}

func (s *scanner) errorf(f string, args ...interface{}) error {
	h := fmt.Sprintf("cursor: %d ", s.cursor-1)
	return fmt.Errorf(h+f, args...)
}

func (s *scanner) isEnd() bool {
	return len(s.source) <= s.cursor
}

func (s *scanner) Get() rune {
	if s.isEnd() {
		return 0
	}
	return rune(s.source[s.cursor])
}

func (s *scanner) Test(r rune) bool {
	return s.Get() == r
}

func (s *scanner) Must(r rune) error {
	if !s.Test(r) {
		return s.errorf("expected %q", r)
	}
	s.Step()
	return nil
}

func (s *scanner) Tests(str string) bool {
	i := s.cursor
	e := i + len([]rune(str))
	if len(s.source) < e {
		return false
	}
	return string(s.source[i:e]) == str
}

func (s *scanner) Musts(str string) error {
	if !s.Tests(str) {
		return s.errorf("expected %q", str)
	}
	s.StepN(len(str))
	return nil
}

func (s *scanner) Step() {
	s.cursor++
}

func (s *scanner) StepN(n int) {
	s.cursor += n
}
