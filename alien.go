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

type Aliens struct {
	Bounds   geo.Rect
	AlienImg [3][2]*ebiten.Image
	Opts     *ebiten.DrawImageOptions
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
	a.Bounds = geo.RectXYWH(50, 100, 100, 50)
}

func (a *Aliens) update(dt time.Duration) {
	a.counter += dt.Seconds() * tempo
}

func (a *Aliens) draw(dst *ebiten.Image) {
	x, y := a.Bounds.TopLeft()
	i := int(a.counter) % 2
	a.Opts.GeoM.Reset()
	a.Opts.GeoM.Translate(x, y)
	dst.DrawImage(a.AlienImg[0][i], a.Opts)
	a.Opts.GeoM.Translate(20, 0)
	dst.DrawImage(a.AlienImg[1][i], a.Opts)
	a.Opts.GeoM.Translate(20, 0)
	dst.DrawImage(a.AlienImg[2][i], a.Opts)
}
