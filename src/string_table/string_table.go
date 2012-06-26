package string_table

import (
	"sync"
	"sequence"
)

type Id uint32

type StringTable struct {
	lock       sync.RWMutex
	sequence   sequence.Uint32Sequence
	stringById map[Id]string
	idByString map[string]Id
}

func New() (t *StringTable) {
	t = new(StringTable)
	t.stringById = make(map[Id]string)
	t.idByString = make(map[string]Id)
	return t
}

func (t *StringTable) Store(s string) (id Id) {
	t.lock.Lock()
	defer t.lock.Unlock()

	id, ok := t.idByString[s]
	if ok {
		return
	} else {
		id := Id(t.sequence.Next())
		t.stringById[id] = s
		t.idByString[s] = id
		return id
	}
	panic("unreachable")
}

func (t *StringTable) Lookup(id Id) (s string, ok bool) {
	t.lock.RLock()
	defer t.lock.RUnlock()

	s, ok = t.stringById[id]
	return
}

func (t *StringTable) Len() uint32 {
	t.lock.RLock()
	defer t.lock.RUnlock()

	return uint32(t.sequence.Peek())
}
