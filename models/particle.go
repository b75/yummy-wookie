package models

import (
	"fmt"
	"io"
)

type Particle struct {
	X     float64
	Y     float64
	Vx    float64
	Vy    float64
	Alive bool
}

func (p *Particle) Update(dt float64) {
	if !p.Alive {
		return
	}
	p.X += p.Vx * dt
	p.Y += p.Vy * dt
}

// TODO real transfer format
func (p *Particle) Serialize(w io.Writer) {
	if !p.Alive {
		return
	}
	w.Write([]byte(fmt.Sprintf("p,%.f,%.f,%.3f,%.3f:", p.X, p.Y, p.Vx, p.Vy)))
}
