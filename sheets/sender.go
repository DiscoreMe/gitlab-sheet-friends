package sheets

import (
	"fmt"
	"unicode/utf8"
)

type Sender struct {
	rows [][]interface{}
	// r and c is the columns and the rows in the table where you will start adding records
	r int
	c string
}

// SetStartRange sets the start of adding entries field
// Example: SetStartRange("B", 4)
// It will start adding entries from the rows variable from the B4 field
func (s *Sender) SetStartRange(c string, r int) {
	s.r = r
	s.c = c
}

func (s *Sender) IsValidRange() bool {
	if s.r < 1 || utf8.RuneCountInString(s.c) != 1 {
		return false
	}
	return true
}

func (s *Sender) AddRows(rows ...interface{}) {
	s.rows = append(s.rows, rows)
}

func (s *Sender) Rows() [][]interface{} {
	return s.rows
}

func (s *Sender) StartRange() string {
	if !s.IsValidRange() {
		return "A3"
	}
	return fmt.Sprintf("%s%d", s.c, s.r)
}
