package main

import (
	"log"
	"math"
	"time"

	"github.com/Bredgren/geo"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
)

const (
	PlayerSpeed = 75
	PlayerY     = 20
)

type Player struct {
	Rect  geo.Rect
	Speed float64
	Img   *ebiten.Image
	Opts  *ebiten.DrawImageOptions
}

func (p *Player) init() {
	img, _, err := ebitenutil.NewImageFromFile("img/player.png", ebiten.FilterNearest)
	if err != nil {
		log.Fatal("open player file:", err)
	}
	p.Img = img
	p.Opts = &ebiten.DrawImageOptions{}
	p.Opts.ColorM.Scale(0.0, 1.0, 0.0, 1.0)
	size := geo.VecXYi(p.Img.Size())
	p.Rect = geo.RectWH(size.XY())
	p.Rect.SetMid(Width*0.25, Height-PlayerY)
	// The player image doesn't display quite right initially unless integer aligned
	p.Rect.SetLeft(math.Trunc(p.Rect.Left()))
}

func (p *Player) move(dt time.Duration) {
	vel := 0.0
	if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		vel = PlayerSpeed
	}
	if ebiten.IsKeyPressed(ebiten.KeyRight) {
		vel = -PlayerSpeed
	}
	newX := p.Rect.Left() + vel*dt.Seconds()
	p.Rect.SetLeft(newX)
	// Keep on screen
	p.Rect.Clamp(geo.RectWH(Width, Height))
}

func (p *Player) draw(dst *ebiten.Image) {
	p.Opts.GeoM.Reset()
	p.Opts.GeoM.Translate(p.Rect.TopLeft())
	dst.DrawImage(p.Img, p.Opts)
}
