package main

import (
	"time"
	"sync"
)

type Stat struct {
	lock           sync.RWMutex
	requests       uint64
	bytes          uint64
	bytesByPath    map[Id]uint64
	bytesByReferer map[Id]uint64
}

type StatRing struct {
	ring     [24]*Stat
	lastHour int
}

func NewStat() *Stat {
	s := new(Stat)
	s.bytesByPath = make(map[Id]uint64, 1024)
	s.bytesByReferer = make(map[Id]uint64, 1024)
	return s
}

func NewStatRing() (r *StatRing) {
	r = new(StatRing)
	r.lastHour = -1
	return
}

func (r *StatRing) Current() *Stat {
	h := time.Now().Hour()
	if r.lastHour != h {
		r.ring[h] = NewStat()
		r.lastHour = h
	}
	return r.ring[h]
}
