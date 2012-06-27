package main

import (
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"path"
	"runtime"
	"sort"
	"string_table"
)

func collect(w http.ResponseWriter, r *http.Request) {
	rx <- r
}

func stats(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "path names: %d\n", pathTable.Len())
	fmt.Fprintf(w, "referer names: %d\n", refererTable.Len())
	fmt.Fprintf(w, "updater queue length: %d\n", len(rx))
}

func renderIndex(w http.ResponseWriter, rng *StatRing, title string) {
	t := template.New("base")
	funcs := template.FuncMap{"humanizeSize": humanizeSize}
	template.Must(t.Funcs(funcs).ParseFiles("templates/index.html.template"))

	data := calculateComposedStats(rng)

	data.Title = title

	data.Buckets = ringByBucket.Keys()
	sort.StringSlice(data.Buckets).Sort()
	data.Servers = ringByServer.Keys()
	sort.StringSlice(data.Servers).Sort()

	err := t.ExecuteTemplate(w, "index.html.template", data)
	if err != nil {
		log.Printf("TEMPLATE ERROR: %s", err)
	}
}

func index(w http.ResponseWriter, r *http.Request) {
	defer func() {
		err := recover()
		if err != nil {
			log.Printf("PANIC: %s", err)
			w.Write([]byte("Internal error"))
		}
	}()

	renderIndex(w, ring, "Global")
}

func bucketIndex(w http.ResponseWriter, r *http.Request) {
	defer func() {
		err := recover()
		if err != nil {
			log.Printf("PANIC: %s", err)
			w.Write([]byte("Internal error"))
		}
	}()

	if bucket, ok := stripPrefix(r.URL.Path, "/bucket/"); ok {
		if rng, ok := ringByBucket.Lookup(bucket); ok {
			renderIndex(w, rng, fmt.Sprintf("%s bucket", bucket))
		} else {
			w.WriteHeader(404)
		}
	}
}

func serverIndex(w http.ResponseWriter, r *http.Request) {
	defer func() {
		err := recover()
		if err != nil {
			log.Printf("PANIC: %s", err)
			w.Write([]byte("Internal error"))
		}
	}()

	t := template.New("base")
	funcs := template.FuncMap{"humanizeSize": humanizeSize}
	template.Must(t.Funcs(funcs).ParseFiles("templates/index.html.template"))

	if server, ok := stripPrefix(r.URL.Path, "/server/"); ok {
		if rng, ok := ringByServer.Lookup(server); ok {
			renderIndex(w, rng, fmt.Sprintf("%s server", server))
		} else {
			w.WriteHeader(404)
		}
	}
}

var host = flag.String("h", "127.0.0.1", "host address (default 127.0.0.1)")
var port = flag.Int("p", 9090, "port (default 9090)")
var pathHashtableSize = flag.Int("phts", 1000, "path hashtable size (default 1000)")

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	err := os.Chdir(path.Join(path.Dir(os.Args[0]), ".."))
	if err != nil {
		log.Fatal(err)
	}

	flag.Parse()

	pathTable = string_table.NewPreallocated(*pathHashtableSize)

	http.HandleFunc("/", index)
	http.HandleFunc("/collect", collect)
	http.HandleFunc("/bucket/", bucketIndex)
	http.HandleFunc("/server/", serverIndex)
	http.HandleFunc("/stats", stats)

	http.Handle("/favicon.ico", http.NotFoundHandler())

	go updater(rx)

	err = http.ListenAndServe(fmt.Sprintf("%s:%d", *host, *port), nil)
	if err != nil {
		log.Fatal(err)
	}
}
