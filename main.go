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

var ring = NewStatRing()
var rx = make(chan *http.Request, 1024)

func collect(w http.ResponseWriter, r *http.Request) {
	rx <- r
}

func index(w http.ResponseWriter, r *http.Request) {
	defer func() {
		err := recover()
		if err != nil {
			log.Printf("PANIC: %s", err)
			w.Write([]byte("Internal error"))
		}
	}()

	t := template.New("base")
	funcs := template.FuncMap{"humanizeSize": humanizeSize}
	template.Must(t.Funcs(funcs).ParseFiles("views/index.html.template"))

	data := calculateComposedStats(ring)
	err := t.ExecuteTemplate(w, "index.html.template", data)
	if err != nil {
		log.Printf("TEMPLATE ERROR: %s", err)
	}
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

				referer := normalizeReferer(r.Header.Get("Referer"))
				if len(referer) > 0 {
					if cb, ok := s.bytesByReferer[referer]; ok {
						s.bytesByReferer[referer] = cb + b
					} else {
						s.bytesByReferer[referer] = b
					}
				}

				path := fmt.Sprintf("%s:%s", r.FormValue("bucket"), r.FormValue("uri"))
				if len(path) > 0 {
					if cb, ok := s.bytesByPath[path]; ok {
						s.bytesByPath[path] = cb + b
					} else {
						s.bytesByPath[path] = b
					}
				}
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
