package main

import (
	"sync"
)

type NameMapper struct {
	lock     sync.RWMutex
	seq      *Sequence
	nameById map[Id]string
	idByName map[string]Id
}

func NewNameMapper() (m *NameMapper) {
	m = new(NameMapper)
	m.seq = new(Sequence)
	m.nameById = make(map[Id]string, 1024)
	m.idByName = make(map[string]Id, 1024)
	return m
}

func (m *NameMapper) NameToId(n string) (id Id) {
	m.lock.Lock()
	defer m.lock.Unlock()

	id, ok := m.idByName[n]
	if ok {
		return id
	} else {
		newId := m.seq.NextId()
		m.nameById[newId] = n
		m.idByName[n] = newId
		return newId
	}
	panic("unreachable")
}

func (m *NameMapper) IdToName(id Id) (n string, ok bool) {
	m.lock.RLock()
	defer m.lock.RUnlock()

	n, ok = m.nameById[id]
	return n, ok
}
