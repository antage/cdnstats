package sequence

import (
	"sync"
)

type Id16 uint16
type Id32 uint32
type Id64 uint64

type Uint16Sequence struct {
	lock   sync.RWMutex
	lastId Id16
}

type Uint32Sequence struct {
	lock   sync.RWMutex
	lastId Id32
}

type Uint64Sequence struct {
	lock   sync.RWMutex
	lastId Id64
}

func (s *Uint16Sequence) Peek() (id Id16) {
	s.lock.RLock()
	defer s.lock.RUnlock()

	return s.lastId
}

func (s *Uint16Sequence) Next() (id Id16) {
	s.lock.Lock()
	defer s.lock.Unlock()

	id = s.lastId
	s.lastId++

	return id
}

func (s *Uint32Sequence) Peek() (id Id32) {
	s.lock.RLock()
	defer s.lock.RUnlock()

	return s.lastId
}

func (s *Uint32Sequence) Next() (id Id32) {
	s.lock.Lock()
	defer s.lock.Unlock()

	id = s.lastId
	s.lastId++

	return id
}

func (s *Uint64Sequence) Peek() (id Id64) {
	s.lock.RLock()
	defer s.lock.RUnlock()

	return s.lastId
}

func (s *Uint64Sequence) Next() (id Id64) {
	s.lock.Lock()
	defer s.lock.Unlock()

	id = s.lastId
	s.lastId++

	return id
}
