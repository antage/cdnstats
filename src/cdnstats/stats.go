package main

import (
	"sort"
	"string_table"
	"sync"
	"time"
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

type PathStatSlice []PathStat
type RefererStatSlice []RefererStat

type StatByPathAndReferer struct {
	StatWithRequests

	PathStats    PathStatSlice
	RefererStats RefererStatSlice

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
		if r.lastHour >= 0 {
			go postProcess(r.ring[r.lastHour])
		}
		r.ring[currentHour] = NewStatByPathAndReferer()
		r.lastHour = currentHour
	}
	return r.ring[currentHour]
}

func postProcess(s *StatByPathAndReferer) {
	s.lock.Lock()
	defer s.lock.Unlock()

	s.PathStats = make([]PathStat, len(s.statByPath))
	i := 0
	for p, cs := range s.statByPath {
		s.PathStats[i] = PathStat{Stat{cs.Bytes}, p}
		i++
	}
	// clear map
	s.statByPath = nil
	// sort by .Bytes
	sort.Sort(s.PathStats)
	// truncate slice to 1024 elements
	if len(s.PathStats) > 1024 {
		newSlice := make([]PathStat, 1024)
		copy(newSlice, s.PathStats[len(s.PathStats)-1024:])
		s.PathStats = newSlice
	}

	s.RefererStats = make([]RefererStat, len(s.statByReferer))
	i = 0
	for r, cs := range s.statByReferer {
		s.RefererStats[i] = RefererStat{Stat{cs.Bytes}, r}
		i++
	}
	// clear map
	s.statByReferer = nil
	// sort by .Bytes
	sort.Sort(s.RefererStats)
	// truncate slice to 1024 elements
	if len(s.RefererStats) > 1024 {
		newSlice := make([]RefererStat, 1024)
		copy(newSlice, s.RefererStats[len(s.RefererStats)-1024:])
		s.RefererStats = newSlice
	}

}

// sort.Interface implementation
func (a PathStatSlice) Len() int {
	return len(a)
}

func (a PathStatSlice) Less(i, j int) bool {
	return a[i].Stat.Bytes < a[j].Stat.Bytes
}

func (a PathStatSlice) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

// sort.Interface implementation
func (a RefererStatSlice) Len() int {
	return len(a)
}

func (a RefererStatSlice) Less(i, j int) bool {
	return a[i].Bytes < a[j].Bytes
}

func (a RefererStatSlice) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}
