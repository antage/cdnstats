package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"runtime"
	"html/template"
	"strconv"
)

type stat struct {
	requests uint64
	Rps      uint64
	bytes    uint64
	Bps      uint64
}

type composedStats struct {
	Summary stat
	ByHour  [24]stat
}

var ring = NewStatRing()
var rx = make(chan *http.Request, 1024)

func collect(w http.ResponseWriter, r *http.Request) {
	rx <- r
}

func index(w http.ResponseWriter, r *http.Request) {
	t := template.Must(template.ParseFiles("views/index.html.template"))

	var data composedStats
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
	t.Execute(w, data)
}

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
			}
		}()
	}
}

var host = flag.String("h", "127.0.0.1", "host address (default 127.0.0.1)")
var port = flag.Int("p", 9090, "port (default 9090)")

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	flag.Parse()

	http.HandleFunc("/", index)
	http.HandleFunc("/collect", collect)

	http.Handle("/assets/", http.FileServer(http.Dir("assets")))
	http.Handle("/favicon.ico", http.NotFoundHandler())

	go updater(rx)

	err := http.ListenAndServe(fmt.Sprintf("%s:%d", *host, *port), nil)
	if err != nil {
		log.Fatal(err)
	}
}
