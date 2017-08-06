package main

import (
	"time"

	"github.com/Bredgren/geo"
	"github.com/hajimehoshi/ebiten"
)

var (
	tempo = 1 * time.Second
)

type Aliens struct {
	Bounds   geo.Rect
	AlienImg [3][2]*ebiten.Image
	Opts     *ebiten.DrawImageOptions
	img      int
	counter  time.Duration
}

func (a *Aliens) init() {
	a.AlienImg[0][0] = openImg("img/alien1_00.png")
	a.AlienImg[0][1] = openImg("img/alien1_01.png")
	a.AlienImg[1][0] = openImg("img/alien2_00.png")
	a.AlienImg[1][1] = openImg("img/alien2_01.png")
	a.AlienImg[2][0] = openImg("img/alien3_00.png")
	a.AlienImg[2][1] = openImg("img/alien3_01.png")
	a.Opts = &ebiten.DrawImageOptions{}
	a.Bounds = geo.RectXYWH(50, 100, 100, 50)
	a.img = 0
	a.counter = tempo
}

func (a *Aliens) update(dt time.Duration) {
	a.counter -= dt
	if a.counter <= 0 {
		a.counter = tempo
		a.img = (a.img + 1) % 2
	}
}

func (a *Aliens) draw(dst *ebiten.Image) {
	x, y := a.Bounds.TopLeft()
	a.Opts.GeoM.Reset()
	a.Opts.GeoM.Translate(x, y)
	dst.DrawImage(a.AlienImg[0][a.img], a.Opts)
	a.Opts.GeoM.Translate(20, 0)
	dst.DrawImage(a.AlienImg[1][a.img], a.Opts)
	a.Opts.GeoM.Translate(20, 0)
	dst.DrawImage(a.AlienImg[2][a.img], a.Opts)
}
