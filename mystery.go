package main

import (
	"log"
	"time"

	"github.com/Bredgren/geo"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
)

const (
	MysterySpeed = 50
	MysteryY     = 40
)

type Mystery struct {
	Rect geo.Rect
	Img  *ebiten.Image
	Opts *ebiten.DrawImageOptions
}

func (m *Mystery) init() {
	img, _, err := ebitenutil.NewImageFromFile("img/mystery.png", ebiten.FilterNearest)
	if err != nil {
		log.Fatal("open mystery file:", err)
	}
	m.Img = img
	m.Opts = &ebiten.DrawImageOptions{}
	m.Opts.ColorM.Scale(1.0, 0.0, 0.0, 1.0)
	size := geo.VecXYi(m.Img.Size())
	m.Rect = geo.RectWH(size.XY())
	m.hide()
}

func (m *Mystery) update(dt time.Duration) {
	newX := m.Rect.Left() - MysterySpeed*dt.Seconds()
	newX = geo.Mod(newX, Width)
	m.Rect.SetLeft(newX)
}

func (m *Mystery) draw(dst *ebiten.Image) {
	m.Opts.GeoM.Reset()
	m.Opts.GeoM.Translate(m.Rect.TopLeft())
	dst.DrawImage(m.Img, m.Opts)
}

func (m *Mystery) hide() {
	m.Rect.SetTopLeft(-100, MysteryY)
}
