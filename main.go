package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"runtime"
)

func collect(w http.ResponseWriter, r *http.Request) {
	log.Println(r)
}

var host = flag.String("h", "127.0.0.1", "host address (default 127.0.0.1)")
var port = flag.Int("p", 9090, "port (default 9090)")

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	flag.Parse()
	http.HandleFunc("/collect", collect)
	http.Handle("/favicon.ico", http.NotFoundHandler())
	http.Handle("/", http.FileServer(http.Dir("assets")))
	err := http.ListenAndServe(fmt.Sprintf("%s:%d", *host, *port), nil)
	if err != nil {
		log.Fatal(err)
	}
}
