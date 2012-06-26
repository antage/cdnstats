package main

import (
	"time"
	"sync"
	"string_table"
)

type Stat struct {
	Bytes uint64
}

type StatWithRequests struct {
	Stat
	Requests uint64
}

type PathStat struct {
	Stat
	Path string_table.Id
}

type RefererStat struct {
	Stat
	Referer string_table.Id
}

type StatByPathAndReferer struct {
	StatWithRequests

	PathStats    []PathStat
	RefererStats []RefererStat

	// aux fields
	lock          sync.RWMutex
	statByPath    map[string_table.Id]Stat
	statByReferer map[string_table.Id]Stat
}

type StatRing struct {
	lock     sync.Mutex
	ring     [24]*StatByPathAndReferer
	lastHour int
}

func NewStatByPathAndReferer() *StatByPathAndReferer {
	s := new(StatByPathAndReferer)
	s.statByPath = make(map[string_table.Id]Stat)
	s.statByReferer = make(map[string_table.Id]Stat)
	return s
}

func NewStatRing() (r *StatRing) {
	r = new(StatRing)
	r.lastHour = -1
	return
}

func (r *StatRing) Current() *StatByPathAndReferer {
	r.lock.Lock()
	defer r.lock.Unlock()

	currentHour := time.Now().Hour()
	if r.lastHour != currentHour {
		r.ring[currentHour] = NewStatByPathAndReferer()
		r.lastHour = currentHour
	}
	return r.ring[currentHour]
}
