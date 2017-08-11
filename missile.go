package main

import (
	"math"
	"math/rand"
	"time"

	"github.com/Bredgren/geo"
	"github.com/hajimehoshi/ebiten"
)

var (
	missiles     Missiles
	missileSpeed = 40.0
)

type Missiles struct {
	MissileImg   *ebiten.Image
	Opts         *ebiten.DrawImageOptions
	Missiles     [3]geo.Rect
	NextMissle   geo.NumGen
	TimeToMissle time.Duration
}

func (m *Missiles) init() {
	m.MissileImg = openImg("img/missile.png")
	m.Opts = &ebiten.DrawImageOptions{}

	m.NextMissle = geo.RandNum(1, 10)
	m.TimeToMissle = time.Duration(m.NextMissle()) * time.Second

	for i := range m.Missiles {
		r := &m.Missiles[i]
		r.SetTop(-100) // arbitratraily off screen
	}
}

func (m *Missiles) update(dt time.Duration) {
	m.TimeToMissle -= dt
	if m.TimeToMissle <= 0 {
		m.TimeToMissle = time.Duration(m.NextMissle()) * time.Second

		for i := range m.Missiles {
			r := &m.Missiles[i]
			if r.Y < 0 { // on screen is considered active, don't use it then
				aliens := aliens.activeAliens()
				i := rand.Intn(len(aliens))
				a := aliens[i]
				r.SetTopMid(a.Rect.BottomMid())
				// Align to avoid visual artifacts
				r.SetTopLeft(math.Trunc(r.Left()), math.Trunc(r.Top()))
				break // Only activate one
			}
		}
	}

	for i := range m.Missiles {
		r := &m.Missiles[i]
		if r.Y < 0 {
			continue
		}
		r.Move(0, missileSpeed*dt.Seconds())

		if r.Top() > Width {
			r.SetTop(-100) // arbitratraily off screen
		}
	}
}

func (m *Missiles) draw(dst *ebiten.Image) {
	for i := range m.Missiles {
		r := &m.Missiles[i]
		if r.Y < 0 {
			continue
		}
		m.Opts.GeoM.Reset()
		m.Opts.GeoM.Translate(r.TopLeft())
		dst.DrawImage(m.MissileImg, m.Opts)
	}
}
