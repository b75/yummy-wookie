package model

import (
	"fmt"
	"io"
	"math"
	"strings"

	"github.com/b75/yummy-wookie/gameutil"
)

type Human struct {
	X         float64
	Y         float64
	Vx        float64
	Vy        float64
	Dir       float64 // radians
	Thrusting int
	Exists    bool
	Actions   []string
}

func (h *Human) Update(dt float64) {
	if !h.Exists {
		return
	}

	// actions
	for _, act := range h.Actions {
		switch act {
		case "tl":
			h.Dir -= 0.2
		case "tr":
			h.Dir += 0.2
		case "mf":
			h.Vx += 40.0 * dt * math.Cos(h.Dir)
			h.Vy += 40.0 * dt * math.Sin(h.Dir)
			h.Thrusting = 3
		case "mb":
			h.Vx -= 32.0 * dt * math.Cos(h.Dir)
			h.Vy -= 32.0 * dt * math.Sin(h.Dir)
			h.Thrusting = 3
		}
	}
	if len(h.Actions) != 0 {
		h.Actions = h.Actions[:0] // reuse underlying array
	}

	// position
	h.X += h.Vx * dt
	h.Y += h.Vy * dt

	// velocity
	dampX := h.Vx - h.Vx*0.001
	dampY := h.Vy - h.Vy*0.001
	h.Vx -= dampX * dt
	h.Vy -= dampY * dt
	if h.Thrusting <= 0 && gameutil.VectorLength(h.Vx, h.Vy) < 30.0 {
		h.Vx = 0
		h.Vy = 0
	}
	if h.Thrusting > 0 {
		h.Thrusting--
	}
}

func (h *Human) Serialize(w io.Writer) {
	if !h.Exists {
		return
	}
	w.Write([]byte(fmt.Sprintf("h,%.f,%.f,%.3f:", h.X, h.Y, h.Dir)))
}

// TODO player -> human mapping
func (h *Human) AddActions(msg string) {
	if !h.Exists {
		return
	}
	for _, act := range strings.Split(msg, ":") {
		if len(act) != 0 {
			h.Actions = append(h.Actions, act)
		}
	}
}
