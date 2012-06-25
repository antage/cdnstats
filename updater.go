package main

import (
	"log"
	"strconv"
	"fmt"
	"net/http"
)

func updater(source chan *http.Request) {
	for r := range source {
		func() {
			log.Println(r)
			s := ring.Current()

			s.lock.Lock()
			defer s.lock.Unlock()

			s.requests++
			b, err := strconv.ParseUint(r.Header.Get("X-Bytes-Sent"), 10, 64)
			if err == nil {
				s.bytes += b

				referer := normalizeReferer(r.Header.Get("Referer"))
				refererId := refererMapper.NameToId(referer)
				if len(referer) > 0 {
					if cb, ok := s.bytesByReferer[refererId]; ok {
						s.bytesByReferer[refererId] = cb + b
					} else {
						s.bytesByReferer[refererId] = b
					}
				}

				path := fmt.Sprintf("%s:%s", r.FormValue("bucket"), r.FormValue("uri"))
				pathId := pathMapper.NameToId(path)
				if len(path) > 0 {
					if cb, ok := s.bytesByPath[pathId]; ok {
						s.bytesByPath[pathId] = cb + b
					} else {
						s.bytesByPath[pathId] = b
					}
				}
			}
		}()
	}
}
