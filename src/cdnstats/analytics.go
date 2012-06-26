package main

import (
	"string_table"
)

type DisplayableStat struct {
	requests uint64
	Rps      uint64
	bytes    uint64
	Bps      uint64
}

type namedValue struct {
	Name  string
	Bytes uint64
}

type ComposedDisplaybleStats struct {
	Summary DisplayableStat
	ByHour  [24]DisplayableStat

	ByPath    [50]namedValue
	ByReferer [50]namedValue
}

func extractTop(h map[string_table.Id]Stat, table *string_table.StringTable, r []namedValue) {
	for i, _ := range r {
		for k, s := range h {
			if i > 0 && s.Bytes >= r[i-1].Bytes {
				continue
			}
			if r[i].Bytes < s.Bytes {
				if name, ok := table.Lookup(k); ok {
					r[i] = namedValue{Name: name, Bytes: s.Bytes}
				}
			}
		}
	}
}

func calculateComposedStats(r *StatRing) *ComposedDisplaybleStats {
	data := new(ComposedDisplaybleStats)
	for i, s := range ring.ring {
		if s != nil {
			func() {
				s.lock.RLock()
				defer s.lock.RUnlock()

				data.ByHour[i].requests = s.Requests
				data.ByHour[i].Rps = s.Requests / 3600
				data.ByHour[i].bytes = s.Bytes
				data.ByHour[i].Bps = s.Bytes / 3600

				data.Summary.requests += s.Requests
				data.Summary.bytes += s.Bytes
			}()
		}
	}
	data.Summary.Rps = data.Summary.requests / (24 * 3600)
	data.Summary.Bps = data.Summary.bytes / (24 * 3600)

	summaryByPath := make(map[string_table.Id]Stat)
	summaryByReferer := make(map[string_table.Id]Stat)

	// aggregate all data for last 24 hours
	for _, s := range ring.ring {
		if s != nil {
			func() {
				s.lock.RLock()
				defer s.lock.RUnlock()

				for k, sc := range s.statByPath {
					if ss, ok := summaryByPath[k]; ok {
						summaryByPath[k] = Stat{ss.Bytes + sc.Bytes}
					} else {
						summaryByPath[k] = sc
					}
				}
				for k, sc := range s.statByReferer {
					if ss, ok := summaryByReferer[k]; ok {
						summaryByReferer[k] = Stat{ss.Bytes + sc.Bytes}
					} else {
						summaryByReferer[k] = sc
					}
				}
			}()
		}
	}

	// extract top values
	extractTop(summaryByPath, pathTable, data.ByPath[:])
	extractTop(summaryByReferer, refererTable, data.ByReferer[:])

	return data
}
