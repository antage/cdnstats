package main

import (
	"time"
	"sync"
)

type Stat struct {
	lock     sync.RWMutex
	requests uint64
	bytes    uint64
}

type StatRing struct {
	ring     [24]*Stat
	lastHour int
}

func NewStatRing() (r *StatRing) {
	r = new(StatRing)
	r.lastHour = -1
	return
}

func (r *StatRing) Current() *Stat {
	h := time.Now().Hour()
	if r.lastHour != h {
		r.ring[h] = new(Stat)
		r.lastHour = h
	}
	return r.ring[h]
}
