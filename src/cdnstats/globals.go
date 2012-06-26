package main

import (
	"net/http"
	"string_table"
)

var ring = NewStatRing()
var rx = make(chan *http.Request, 1024)

var pathTable = string_table.New()
var refererTable = string_table.New()
