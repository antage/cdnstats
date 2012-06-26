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
)

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
	template.Must(t.Funcs(funcs).ParseFiles("templates/index.html.template"))

	data := calculateComposedStats(ring)
	err := t.ExecuteTemplate(w, "index.html.template", data)
	if err != nil {
		log.Printf("TEMPLATE ERROR: %s", err)
	}
}

func stats(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "path names: %d\n", pathTable.Len())
	fmt.Fprintf(w, "referer names: %d\n", refererTable.Len())
	fmt.Fprintf(w, "updater queue length: %d\n", len(rx))
}

var host = flag.String("h", "127.0.0.1", "host address (default 127.0.0.1)")
var port = flag.Int("p", 9090, "port (default 9090)")

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	err := os.Chdir(path.Join(path.Dir(os.Args[0]), ".."))
	if err != nil {
		log.Fatal(err)
	}

	flag.Parse()

	http.HandleFunc("/", index)
	http.HandleFunc("/collect", collect)
	http.HandleFunc("/stats", stats)

	http.Handle("/favicon.ico", http.NotFoundHandler())

	go updater(rx)

	err = http.ListenAndServe(fmt.Sprintf("%s:%d", *host, *port), nil)
	if err != nil {
		log.Fatal(err)
	}
}
