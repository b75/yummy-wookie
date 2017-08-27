package main

import (
	"net/http"
)

func init() {
	http.HandleFunc("/game", gameHandler)
}

func gameHandler(w http.ResponseWriter, rq *http.Request) {
	p := &struct{}{}
	executeTemplate(w, "game.html", p)
}
