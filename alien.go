package main

import (
	"image/color"
	"math"
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
	a.Bounds = geo.RectXYWH(40, 64, 0, 0)
	a.speed = 2

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

	tick := int(math.Trunc(a.counter))
	if tick > a.tick {
		y := 0.0
		if a.Bounds.Right() >= Width || a.Bounds.Left() <= 0 {
			y = 16
			a.speed *= -1
		}
		for i := range a.Aliens {
			a.Aliens[i].Rect.X += a.speed
			a.Aliens[i].Rect.Y += y
		}
		a.Bounds.X += a.speed
		a.Bounds.Y += y
	}
	a.tick = tick
}

func (a *Aliens) draw(dst *ebiten.Image) {
	// Debug rect
	img, _ := ebiten.NewImage(int(a.Bounds.W), int(a.Bounds.H), ebiten.FilterLinear)
	img.Fill(color.NRGBA{0xaa, 0xaa, 0xff, 0x44})
	opts := ebiten.DrawImageOptions{}
	opts.GeoM.Translate(a.Bounds.TopLeft())
	dst.DrawImage(img, &opts)

	idx := int(a.counter) % 2
	for _, alien := range a.Aliens {
		a.Opts.GeoM.Reset()
		a.Opts.GeoM.Translate(alien.Rect.TopLeft())
		dst.DrawImage(a.AlienImg[alien.kind][idx], a.Opts)
	}
}
