package data

import "errors"

type Stack struct {
	s []Snapshot
}

func NewStack() *Stack {
	return &Stack{make([]Snapshot, 0)}
}

func (s *Stack) Push(v Snapshot) {
	s.s = append(s.s, v)
}

func (s *Stack) Pop() (Snapshot, error) {
	// FIXME: What do we do if the stack is empty, though?

	l := len(s.s)
	if l <= 0 {
		return Snapshot{}, errors.New("stack is empty")
	}

	return s.s[l-1], nil
}

func (s *Stack) Len() int {
	return len(s.s)
}

// save partialRoute,stations(so far),partialStart
type Snapshot struct {
	Stations     []Node
	PartialRoute NodeRoute
	PartialStart Node
}
