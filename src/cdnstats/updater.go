package main

import (
	"strconv"
	"fmt"
	"net/http"
)

func update(r *http.Request) {
	s := ring.Current()

	s.lock.Lock()
	defer s.lock.Unlock()

	s.Requests++
	b, err := strconv.ParseUint(r.Header.Get("X-Bytes-Sent"), 10, 64)
	if err == nil {
		s.Bytes += b

		referer := normalizeReferer(r.Header.Get("Referer"))
		refererId := refererTable.Store(referer)
		if len(referer) > 0 {
			if sc, ok := s.statByReferer[refererId]; ok {
				s.statByReferer[refererId] = Stat{sc.Bytes + b}
			} else {
				s.statByReferer[refererId] = Stat{b}
			}
		}

		path := fmt.Sprintf("%s:%s", r.FormValue("bucket"), r.FormValue("uri"))
		pathId := pathTable.Store(path)
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
		update(r)
	}
}
