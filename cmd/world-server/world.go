package main

import (
	"bytes"
	"log"
	"sync"
	"time"

	"github.com/b75/yummy-wookie/model"
)

const (
	UPDATE_INTERVAL time.Duration = 50 * time.Millisecond
)

type World struct {
	StateListeners          map[int]chan []byte
	NextStateListenerNumber int
	StateListenerMutex      *sync.Mutex
	LastUpdate              int64
	LoopNumber              int64
	TotalUpdateTime         int64 // microseconds
	MaxUpdateTime           int64 // nanoseconds
	NumUpdates              int64
	TotalBroadcastTime      int64 // microseconds
	MaxBroadcastTime        int64 // nanoseconds
	NumBroadcasts           int64
	NextStateUpdate         []byte
	Particles               []*model.Particle
	Player                  *model.Human // TODO player to human mapping
}

func NewWorld() *World {
	w := &World{
		StateListeners:     make(map[int]chan []byte),
		StateListenerMutex: &sync.Mutex{},
		Particles:          []*model.Particle{},
		LastUpdate:         time.Now().UnixNano(),
	}

	for i := 0; i < 50; i++ {
		p := &model.Particle{
			X:      0,
			Y:      float64(i),
			Vx:     0,
			Vy:     0,
			Exists: true,
		}

		w.Particles = append(w.Particles, p)
	}

	w.Player = &model.Human{
		X:       0,
		Y:       0,
		Vx:      0,
		Vy:      0,
		Dir:     0,
		Exists:  true,
		Actions: []string{},
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
	i := 0
	for {
		w.Update()
		w.Broadcast()
		time.Sleep(UPDATE_INTERVAL)
		i++
		if i >= 3000 {
			if w.NumUpdates != 0 {
				var avg float64 = (float64(w.TotalUpdateTime) / float64(w.NumUpdates)) * 0.001
				var max float64 = float64(w.MaxUpdateTime) * 1e-6
				log.Printf("update time: avg %.2f ms, max %.2f ms", avg, max)
			}
			if w.NumBroadcasts != 0 {
				var avg float64 = (float64(w.TotalBroadcastTime) / float64(w.NumBroadcasts)) * 0.001
				var max float64 = float64(w.MaxBroadcastTime) * 1e-6
				log.Printf("broadcast time: avg %.2f ms, max %.2f ms", avg, max)
			}
			i = 0
		}
	}
}

func (w *World) Broadcast() {
	now := time.Now()
	w.StateListenerMutex.Lock()
	defer w.StateListenerMutex.Unlock()

	for _, c := range w.StateListeners {
		c <- w.NextStateUpdate
	}
	w.NumBroadcasts++
	broadcastTime := time.Since(now).Nanoseconds()
	w.TotalBroadcastTime += broadcastTime / 1000
	if broadcastTime > w.MaxBroadcastTime {
		w.MaxBroadcastTime = broadcastTime
	}
}

func (w *World) Update() {
	now := time.Now()
	dt := float64(now.UnixNano()-w.LastUpdate) * 10e-10 // seconds

	buf := &bytes.Buffer{}
	for _, p := range w.Particles {
		if !p.Exists {
			continue
		}
		p.Update(dt)
		p.Serialize(buf)
	}
	if p := w.Player; p != nil {
		p.Update(dt)
		p.Serialize(buf)
	}

	w.NextStateUpdate = buf.Bytes()
	w.LastUpdate = time.Now().UnixNano()
	w.LoopNumber++
	w.NumUpdates++
	updateTime := time.Since(now).Nanoseconds()
	w.TotalUpdateTime += updateTime / 1000
	if updateTime > w.MaxUpdateTime {
		w.MaxUpdateTime = updateTime
	}
}
