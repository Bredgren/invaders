package main

import (
	"image/color"
	"math"
	"time"

	"github.com/Bredgren/geo"
	"github.com/hajimehoshi/ebiten"
)

const (
	PlayerSpeed       = 100
	PlayerY           = 20
	PlayerBulletSpeed = 160
)

type Player struct {
	Rect geo.Rect
	Img  *ebiten.Image
	Opts *ebiten.DrawImageOptions
}

func (p *Player) init() {
	p.Img = openImg("img/player.png")
	p.Opts = &ebiten.DrawImageOptions{}
	p.Opts.ColorM.Scale(0.0, 1.0, 0.0, 1.0)
	size := geo.VecXYi(p.Img.Size())
	p.Rect = geo.RectWH(size.XY())
	yOffset := (size.Y + 10) * 2 // enough room for 2 player imgs + padding
	p.Rect.SetMid(Width*ShelterX[0], Height-yOffset)
	// The player image doesn't display quite right initially unless left edge is integer aligned
	p.Rect.SetLeft(math.Trunc(p.Rect.Left()))
}

func (p *Player) update(dt time.Duration) {
	vel := 0.0
	if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		vel = -PlayerSpeed
	}
	if ebiten.IsKeyPressed(ebiten.KeyRight) {
		vel = PlayerSpeed
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

type PlayerBullet struct {
	Rect  geo.Rect
	Speed float64
	Img   *ebiten.Image
	Opts  *ebiten.DrawImageOptions
}

func (b *PlayerBullet) init() {
	const w, h = 1, 4
	b.Img, _ = ebiten.NewImage(w, h, ebiten.FilterNearest)
	b.Opts = &ebiten.DrawImageOptions{}
	b.Img.Fill(color.White)
	b.Rect.SetSize(w, h)
	b.Speed = PlayerBulletSpeed
	b.hide()
}

func (b *PlayerBullet) update(dt time.Duration) {
	if ebiten.IsKeyPressed(ebiten.KeyUp) && !b.isGoing() {
		b.Rect.SetBottomMid(player.Rect.TopMid())
	}
	if !b.isGoing() {
		b.hide()
		return
	}

	b.Rect.SetTop(b.Rect.Top() - b.Speed*dt.Seconds())
}

func (b *PlayerBullet) draw(dst *ebiten.Image) {
	if !b.isGoing() {
		return
	}
	b.Opts.GeoM.Reset()
	b.Opts.GeoM.Translate(b.Rect.TopLeft())
	dst.DrawImage(b.Img, b.Opts)
}

func (b *PlayerBullet) isGoing() bool {
	return b.Rect.Y > 0
}

func (b *PlayerBullet) hitSomething() {
	b.hide()
}

func (b *PlayerBullet) hide() {
	b.Rect.SetTopLeft(0, -100) // arbitraty y < 0
}
