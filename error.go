package xml

import (
	"fmt"
)

type XMLError struct {
	Parsing string
	Err     error
	Pos     Pos
}

func newErr(parsing string, err error, pos Pos) *XMLError {
	return &XMLError{
		Parsing: parsing,
		Err:     err,
		Pos:     pos,
	}
}

func (e *XMLError) Error() string {
	str := fmt.Sprintf("error while parsing %s at line %d column %d", e.Parsing, e.Pos.Line, e.Pos.Col)
	if e.Err != nil {
		str += fmt.Sprintf(": %s", e.Err.Error())
	}
	return str
}
