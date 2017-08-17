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

	for i := range m.Missiles {
		m.Missiles[i].SetSize(toFloat(m.MissileImg.Size()))
	}

	m.NextMissle = geo.RandNum(1, 10)
}

func (m *Missiles) resetLevel(level int) {
	m.TimeToMissle = time.Duration(m.NextMissle()) * time.Second

	for i := range m.Missiles {
		m.hitSomething(i)
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
				abovePlayer := make([]*Alien, 0, len(aliens))
				for i := range aliens {
					if aliens[i].Rect.Bottom() < player.Rect.Top()-8 { // 8 is just a bit of buffer
						abovePlayer = append(abovePlayer, aliens[i])
					}
				}
				if len(abovePlayer) == 0 {
					break // No aliens available to shoot
				}
				i := rand.Intn(len(abovePlayer))
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
			m.hitSomething(i)
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

func (m *Missiles) hitSomething(missile int) {
	m.Missiles[missile].SetTop(-100) // arbitratraily off screen
}
