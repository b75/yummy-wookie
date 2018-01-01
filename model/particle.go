package model

import (
	"fmt"
	"io"
)

type Particle struct {
	X      float64
	Y      float64
	Vx     float64
	Vy     float64
	Exists bool
}

func (p *Particle) Update(dt float64) {
	if !p.Exists {
		return
	}

	// position
	p.X += p.Vx * dt
	p.Y += p.Vy * dt
}

// TODO real transfer format
func (p *Particle) Serialize(w io.Writer) {
	if !p.Exists {
		return
	}
	w.Write([]byte(fmt.Sprintf("p,%.f,%.f:", p.X, p.Y)))
}
