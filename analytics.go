package main

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

func extractTop(h map[string]uint64, r []namedValue) {
	for i, _ := range r {
		for k, b := range h {
			if i > 0 && b >= r[i-1].Bytes {
				continue
			}
			if r[i].Bytes < b {
				r[i] = namedValue{Name: k, Bytes: b}
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

				data.ByHour[i].requests = s.requests
				data.ByHour[i].Rps = s.requests / 3600
				data.ByHour[i].bytes = s.bytes
				data.ByHour[i].Bps = s.bytes / 3600

				data.Summary.requests += s.requests
				data.Summary.bytes += s.bytes
			}()
		}
	}
	data.Summary.Rps = data.Summary.requests / (24 * 3600)
	data.Summary.Bps = data.Summary.bytes / (24 * 3600)

	summaryByPath := make(map[string]uint64, 1024)
	summaryByReferer := make(map[string]uint64, 1024)

	// aggregate all data for last 24 hours
	for _, s := range ring.ring {
		if s != nil {
			func() {
				s.lock.RLock()
				defer s.lock.RUnlock()

				for k, b := range s.bytesByPath {
					if cb, ok := summaryByPath[k]; ok {
						summaryByPath[k] = cb + b
					} else {
						summaryByPath[k] = b
					}
				}
				for k, b := range s.bytesByReferer {
					if cb, ok := summaryByReferer[k]; ok {
						summaryByReferer[k] = cb + b
					} else {
						summaryByReferer[k] = b
					}
				}
			}()
		}
	}

	// extract top values
	extractTop(summaryByPath, data.ByPath[:])
	extractTop(summaryByReferer, data.ByReferer[:])

	return data
}
