package main

import (
	"log"
	"net/http"
)

/*
 * https://godoc.org/github/gorilla/websocket
 */

func connect(w http.ResponseWriter, rq *http.Request) {
	log.Printf("new connection from %s", rq.RemoteAddr)
	c, err := upgrader.Upgrade(w, rq, nil)
	if err != nil {
		log.Printf("upgrade error: %v", err)
		return
	}
	defer c.Close()

	n, worldChan := world.AddStateListener()
	closeChan := make(chan struct{}, 1)
	readChan := make(chan []byte)

	// reader
	go func(rc chan<- []byte, cc chan<- struct{}) {
		for {
			mt, msg, err := c.ReadMessage()
			if err != nil {
				log.Printf("read error: %v", err)
				cc <- struct{}{}
				break
			}
			if mt == 1 {
				rc <- msg
			}
		}
	}(readChan, closeChan)

talkloop:
	for {
		select {
		case <-closeChan:
			break talkloop
		case msg := <-readChan:
			log.Printf("message: %s", msg)
		case upd := <-worldChan:
			if err = c.WriteJSON(upd); err != nil {
				log.Printf("write error: %v", err)
			}
		}
	}

	world.RemoveStateListener(n)
	log.Printf("connection closed from %s", rq.RemoteAddr)
}
