package main

import (
	"net/http"
)

func init() {
	http.HandleFunc("/index", indexHandler)
}

func indexHandler(w http.ResponseWriter, rq *http.Request) {
	p := &struct{}{}
	executeTemplate(w, "index.html", p)
}
