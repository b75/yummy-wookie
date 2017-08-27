package main

import (
	"bytes"
	"sync"
	"time"

	"github.com/b75/yummy-wookie/models"
)

const (
	UPDATE_INTERVAL = 200 * time.Millisecond
)

type World struct {
	StateListeners          map[int]chan []byte
	NextStateListenerNumber int
	StateListenerMutex      *sync.Mutex
	LastUpdate              int64
	LoopNumber              int64
	NextStateUpdate         []byte
	Particles               []*models.Particle
}

func NewWorld() *World {
	w := &World{
		StateListeners:     make(map[int]chan []byte),
		StateListenerMutex: &sync.Mutex{},
		Particles:          []*models.Particle{},
		LastUpdate:         time.Now().UnixNano(),
	}

	for i := 0; i < 50; i++ {
		p := &models.Particle{
			X:     10,
			Y:     5 + float64(i),
			Vx:    1,
			Vy:    0,
			Alive: true,
		}

		w.Particles = append(w.Particles, p)
	}

	return w
}

func (w *World) AddStateListener() (int, chan []byte) {
	w.StateListenerMutex.Lock()
	defer w.StateListenerMutex.Unlock()

	n := w.NextStateListenerNumber
	w.NextStateListenerNumber++

	c := make(chan []byte)
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

func (w *World) Run() {
	for {
		w.Update()
		w.Broadcast()
		time.Sleep(UPDATE_INTERVAL)
	}
}

func (w *World) Broadcast() {
	w.StateListenerMutex.Lock()
	defer w.StateListenerMutex.Unlock()

	for _, c := range w.StateListeners {
		c <- w.NextStateUpdate
	}
}

func (w *World) Update() {
	dt := float64(time.Now().UnixNano()-w.LastUpdate) * 10e-10 // seconds

	buf := &bytes.Buffer{}
	for _, p := range w.Particles {
		if !p.Alive {
			continue
		}
		p.Update(dt)
		p.Serialize(buf)
	}

	w.NextStateUpdate = buf.Bytes()
	w.LastUpdate = time.Now().UnixNano()
	w.LoopNumber++
}
