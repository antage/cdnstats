package main

import (
	"net/http"
	"string_table"
	"sync"
)

type ringByString struct {
	lock sync.RWMutex
	m    map[string]*StatRing
}

var ring = NewStatRing()
var ringByBucket = ringByString{m: make(map[string]*StatRing)}
var ringByServer = ringByString{m: make(map[string]*StatRing)}

var rx = make(chan *http.Request, 1024)

var pathTable *string_table.StringTable
var refererTable = string_table.New()

func (r *ringByString) Lookup(s string) (*StatRing, bool) {
	r.lock.RLock()
	defer r.lock.RUnlock()

	rng, ok := r.m[s]
	return rng, ok
}

func (r *ringByString) LookupOrCreate(s string) *StatRing {
	r.lock.Lock()
	defer r.lock.Unlock()

	rng, ok := r.m[s]
	if ok {
		return rng
	} else {
		rng := NewStatRing()
		r.m[s] = rng
		return rng
	}
	panic("unreachable")
}

func (r *ringByString) Keys() []string {
	r.lock.RLock()
	defer r.lock.RUnlock()

	keys := make([]string, len(r.m))
	i := 0
	for k, _ := range r.m {
		keys[i] = k
		i++
	}
	return keys
}
