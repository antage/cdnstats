package main

import (
	"log"
	"net/http"
    "runtime"
)

func collect(w http.ResponseWriter, r *http.Request) {
    log.Println(r)
}

func main() {
    runtime.GOMAXPROCS(runtime.NumCPU())
    http.HandleFunc("/collect", collect)
    http.Handle("/favicon.ico", http.NotFoundHandler())
    http.Handle("/", http.FileServer(http.Dir("assets")))
	err := http.ListenAndServe(":9090", nil)
	if err != nil {
		log.Fatal(err)
	}
}
