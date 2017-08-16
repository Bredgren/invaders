package main

import (
	"math"
	"time"

	"github.com/Bredgren/geo"
	"github.com/hajimehoshi/ebiten"
)

const (
	rows = 5
	cols = 11
)

var (
	// alien actions per second
	tempo  = 1.0
	aliens Aliens
)

type Alien struct {
	Rect geo.Rect
	kind int
}

type Aliens struct {
	Bounds   geo.Rect
	AlienImg [3][2]*ebiten.Image
	Opts     *ebiten.DrawImageOptions
	Aliens   [rows * cols]Alien
	counter  float64
	tick     int
	speed    float64
}

func (a *Aliens) init() {
	a.AlienImg[0][0] = openImg("img/alien1_00.png")
	a.AlienImg[0][1] = openImg("img/alien1_01.png")
	a.AlienImg[1][0] = openImg("img/alien2_00.png")
	a.AlienImg[1][1] = openImg("img/alien2_01.png")
	a.AlienImg[2][0] = openImg("img/alien3_00.png")
	a.AlienImg[2][1] = openImg("img/alien3_01.png")
	a.Opts = &ebiten.DrawImageOptions{}
}

func (a *Aliens) resetLevel(level int) {
	a.speed = 2

	a.Bounds = geo.RectXYWH(40, 64, 0, 0)

	xSpacing, ySpacing := 16.0, 16.0
	x, y := a.Bounds.TopLeft()

	for row := 0; row < rows; row++ {
		kind := alienKindForRow(row)
		offset := float64(kind)
		wi, hi := a.AlienImg[kind][0].Size()
		w, h := float64(wi), float64(hi)
		for col := 0; col < cols; col++ {
			a.Aliens[cols*row+col] = Alien{
				Rect: geo.RectXYWH(x+offset, y, w, h),
				kind: kind,
			}
			x += xSpacing
		}
		x = a.Bounds.Left()
		y += ySpacing
	}

	a.reCalcBounds()

	tempo = 1.0
}

func alienKindForRow(row int) int {
	switch row {
	case 0:
		return 2
	case 1, 2:
		return 1
	}
	return 0
}

func (a *Aliens) activeAliens() []*Alien {
	active := make([]*Alien, 0, len(a.Aliens))
	for i := range a.Aliens {
		alien := &a.Aliens[i]
		if alien.kind < 0 {
			continue
		}
		active = append(active, alien)
	}
	return active
}

func (a *Aliens) update(dt time.Duration) {
	a.counter += dt.Seconds() * tempo

	tick := int(math.Trunc(a.counter))
	if tick > a.tick {
		y := 0.0
		if a.Bounds.Right() >= Width || a.Bounds.Left() <= 0 {
			y = 8
			a.speed *= -1
		}
		for i := range a.Aliens {
			alien := &a.Aliens[i]
			if alien.kind < 0 {
				continue
			}
			alien.Rect.X += a.speed
			alien.Rect.Y += y
		}
		a.Bounds.X += a.speed
		a.Bounds.Y += y
	}
	a.tick = tick
}

func (a *Aliens) draw(dst *ebiten.Image) {
	// Debug rect
	// img, _ := ebiten.NewImage(int(a.Bounds.W), int(a.Bounds.H), ebiten.FilterLinear)
	// img.Fill(color.NRGBA{0xaa, 0xaa, 0xff, 0x44})
	// opts := ebiten.DrawImageOptions{}
	// opts.GeoM.Translate(a.Bounds.TopLeft())
	// dst.DrawImage(img, &opts)

	idx := int(a.counter) % 2
	for i := range a.Aliens {
		alien := &a.Aliens[i]
		if alien.kind < 0 {
			continue
		}
		a.Opts.GeoM.Reset()
		a.Opts.GeoM.Translate(alien.Rect.TopLeft())
		dst.DrawImage(a.AlienImg[alien.kind][idx], a.Opts)
	}
}

func (a *Aliens) collidePlayerBullet(b *PlayerBullet) {
	if !a.Bounds.CollideRect(b.Rect) {
		return
	}

	for i := range a.Aliens {
		alien := &a.Aliens[i]
		if alien.kind < 0 {
			continue
		}
		if alien.Rect.CollideRect(b.Rect) {
			score += (alien.kind + 1) * 10
			alien.kind = -1
			b.hitSomething()
			a.reCalcBounds()
			tempo = a.getTempo()
			break
		}
	}
}

func (a *Aliens) getTempo() float64 {
	remaining := float64(a.remaining())
	total := float64(rows * cols)
	destroyed := total - remaining
	switch {
	case destroyed < total*0.2:
		return 1.0
	case destroyed < total*0.4:
		return 1.5
	case destroyed < total*0.7:
		return 2.5
	case destroyed < total*0.9:
		return 5.0
	case destroyed < total-1:
		return 10.0
	case destroyed < total:
		return 30.0
	}
	return 1.0
}

func (a *Aliens) remaining() int {
	count := 0
	for i := range a.Aliens {
		if a.Aliens[i].kind >= 0 {
			count++
		}
	}
	return count
}

func (a *Aliens) reCalcBounds() {
	if a.remaining() == 0 {
		return
	}

	// Find first alien's rect for a starting point, then union with all others
	a.Bounds.X = -100
	for i := range a.Aliens {
		alien := &a.Aliens[i]
		if alien.kind < 0 {
			continue
		}
		if a.Bounds.X == -100 {
			a.Bounds = alien.Rect
		} else {
			a.Bounds.Union(alien.Rect)
		}
	}
}

func (a *Aliens) collideShelters() {
	for i := range a.Aliens {
		alien := &a.Aliens[i]
		if alien.kind < 0 || alien.Rect.Bottom() < ShelterBottomY-ShelterH {
			continue
		}
		for i := range shelters {
			shelters[i].clearRect(alien.Rect)
		}
	}
}
