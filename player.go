package main

import (
	"image/color"
	"log"
	"math"
	"time"

	"github.com/Bredgren/geo"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
)

const (
	PlayerSpeed       = 80
	PlayerY           = 20
	PlayerBulletSpeed = 160
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
	yOffset := (size.Y + 4) * 2 // enough room for 2 player imgs + padding
	p.Rect.SetMid(Width*0.25, Height-yOffset)
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
}

func (b *PlayerBullet) update(dt time.Duration) {
	if ebiten.IsKeyPressed(ebiten.KeyUp) && !b.isGoing() {
		b.Rect.SetBottomMid(player.Rect.TopMid())
	}
	if !b.isGoing() {
		b.Rect.SetTopLeft(0, -100) // arbitraty y < 0
		return
	}

	b.Rect.SetTop(b.Rect.Top() - b.Speed*dt.Seconds())
}

func (b *PlayerBullet) draw(dst *ebiten.Image) {
	b.Opts.GeoM.Reset()
	b.Opts.GeoM.Translate(b.Rect.TopLeft())
	dst.DrawImage(b.Img, b.Opts)
}

func (b *PlayerBullet) isGoing() bool {
	return b.Rect.Y > 0
}
