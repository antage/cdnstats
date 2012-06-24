package main

import (
	"sync"
)

type Id uint

type Sequence struct {
	lock   sync.Mutex
	lastId Id
}

func (s *Sequence) PeekId() (id Id) {
	s.lock.Lock()
	defer s.lock.Unlock()

	return s.lastId
}

func (s *Sequence) NextId() (id Id) {
	s.lock.Lock()
	defer s.lock.Unlock()

	id = s.lastId
	s.lastId++

	return id
}
