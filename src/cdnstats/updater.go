package main

import (
	"net/http"
	"strconv"
)

func update(r *http.Request, rng *StatRing) {
	s := rng.Current()

	s.lock.Lock()
	defer s.lock.Unlock()

	s.Requests++
	b, err := strconv.ParseUint(r.Header.Get("X-Bytes-Sent"), 10, 64)
	if err == nil {
		s.Bytes += b

		referer := normalizeReferer(r.Header.Get("Referer"))

		// copy referer to avoid memory leaks
		refererCopy := make([]byte, len(referer))
		copy(refererCopy, referer)

		refererId := refererTable.Store(string(refererCopy))
		if len(referer) > 0 {
			if sc, ok := s.statByReferer[refererId]; ok {
				s.statByReferer[refererId] = Stat{sc.Bytes + b}
			} else {
				s.statByReferer[refererId] = Stat{b}
			}
		}

		path := r.FormValue("uri")

		// copy path to avoid memory leaks
		pathCopy := make([]byte, len(path))
		copy(pathCopy, path)

		pathId := pathTable.Store(string(pathCopy))
		if len(path) > 0 {
			if sc, ok := s.statByPath[pathId]; ok {
				s.statByPath[pathId] = Stat{sc.Bytes + b}
			} else {
				s.statByPath[pathId] = Stat{b}
			}
		}
	}
}

func updater(source chan *http.Request) {
	for r := range source {
		// update global ring
		go update(r, ring)

		// update bucket ring
		bucket := r.FormValue("bucket")
		go update(r, ringByBucket.LookupOrCreate(bucket))

		// update server ring
		server := r.FormValue("s")
		go update(r, ringByServer.LookupOrCreate(server))
	}
}
