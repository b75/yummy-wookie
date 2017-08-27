package main

import (
	"fmt"
	"sync"
	"time"
)

type StateUpdate struct {
	Msg string
}

type World struct {
	StateListeners          map[int]chan *StateUpdate
	NextStateListenerNumber int
	StateListenerMutex      *sync.Mutex
}

func NewWorld() *World {
	w := &World{
		StateListeners:     make(map[int]chan *StateUpdate),
		StateListenerMutex: &sync.Mutex{},
	}

	return w
}

func (w *World) AddStateListener() (int, chan *StateUpdate) {
	w.StateListenerMutex.Lock()
	defer w.StateListenerMutex.Unlock()

	n := w.NextStateListenerNumber
	w.NextStateListenerNumber++

	c := make(chan *StateUpdate)
	w.StateListeners[n] = c

	return n, c
}

func (w *World) RemoveStateListener(n int) {
	w.StateListenerMutex.Lock()
	defer w.StateListenerMutex.Unlock()

	c, ok := w.StateListeners[n]
	if !ok {
		return
	}

	delete(w.StateListeners, n)
	close(c)
}

func (w *World) Broadcast() {
	w.StateListenerMutex.Lock()
	defer w.StateListenerMutex.Unlock()

	upd := &StateUpdate{
		Msg: fmt.Sprintf("world time is %d, I have %d listeners", time.Now().Second(), len(w.StateListeners)),
	}

	for _, c := range w.StateListeners {
		c <- upd
	}
}

func (w *World) Run() {
	for {
		w.Broadcast()
		time.Sleep(time.Second)
	}
}
