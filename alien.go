package main

import (
	"time"

	"github.com/Bredgren/geo"
	"github.com/hajimehoshi/ebiten"
)

var (
	// alien actions per second
	tempo = 1.0
)

type Alien struct {
	Rect geo.Rect
	kind int
}

type Aliens struct {
	Bounds   geo.Rect
	AlienImg [3][2]*ebiten.Image
	Opts     *ebiten.DrawImageOptions
	Aliens   [11 * 5]Alien
	counter  float64
}

func (a *Aliens) init() {
	a.AlienImg[0][0] = openImg("img/alien1_00.png")
	a.AlienImg[0][1] = openImg("img/alien1_01.png")
	a.AlienImg[1][0] = openImg("img/alien2_00.png")
	a.AlienImg[1][1] = openImg("img/alien2_01.png")
	a.AlienImg[2][0] = openImg("img/alien3_00.png")
	a.AlienImg[2][1] = openImg("img/alien3_01.png")
	a.Opts = &ebiten.DrawImageOptions{}
	a.Bounds = geo.RectXYWH(40, 64, 0, 0)

	xSpacing, ySpacing := 16.0, 16.0
	x, y := a.Bounds.TopLeft()

	for row := 0; row < 5; row++ {
		kind := alienKindForRow(row)
		offset := float64(kind)
		wi, hi := a.AlienImg[kind][0].Size()
		w, h := float64(wi), float64(hi)
		for col := 0; col < 11; col++ {
			a.Aliens[11*row+col] = Alien{
				Rect: geo.RectXYWH(x+offset, y, w, h),
				kind: kind,
			}
			x += xSpacing
		}
		x = a.Bounds.Left()
		y += ySpacing
	}

	for _, alien := range a.Aliens {
		a.Bounds.Union(alien.Rect)
	}
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

func (a *Aliens) update(dt time.Duration) {
	a.counter += dt.Seconds() * tempo
}

func (a *Aliens) draw(dst *ebiten.Image) {
	idx := int(a.counter) % 2
	for _, alien := range a.Aliens {
		a.Opts.GeoM.Reset()
		a.Opts.GeoM.Translate(alien.Rect.TopLeft())
		dst.DrawImage(a.AlienImg[alien.kind][idx], a.Opts)
	}
}
