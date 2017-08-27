package main

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var world *World = NewWorld()

var upgrader = websocket.Upgrader{
	CheckOrigin: func(rq *http.Request) bool {
		return true // TODO
	},
}

func main() {
	go world.Run()
	http.HandleFunc("/connect", connect)
	log.Print("listening on 8081")
	log.Fatal(http.ListenAndServe("localhost:8081", nil))
}
