package main

import (
	"time"

	"github.com/Bredgren/geo"
	"github.com/hajimehoshi/ebiten"
)

const (
	MysterySpeed = 50
	MysteryY     = 40
)

var (
	mystery Mystery
)

type Mystery struct {
	Rect        geo.Rect
	Img         *ebiten.Image
	Opts        *ebiten.DrawImageOptions
	nextGoTime  time.Duration
	getNextTime geo.NumGen
}

func (m *Mystery) init() {
	m.Img = openImg("img/mystery.png")
	m.Opts = &ebiten.DrawImageOptions{}
	m.Opts.ColorM.Scale(1.0, 0.0, 0.0, 1.0)
	size := geo.VecXYi(m.Img.Size())
	m.Rect = geo.RectWH(size.XY())
	m.getNextTime = geo.RandNum(20, 60)
}

func (m *Mystery) resetLevel(level int) {
	m.hide()
	m.nextGoTime = 5 * time.Second
}

func (m *Mystery) update(dt time.Duration) {
	m.nextGoTime -= dt
	if !m.isGoing() && m.nextGoTime < 0 {
		m.Rect.X = Width + 50
		m.nextGoTime = time.Duration(int(m.getNextTime())) * time.Second
	}
	if !m.isGoing() {
		return
	}

	newX := m.Rect.Left() - MysterySpeed*dt.Seconds()
	m.Rect.SetLeft(newX)
}

func (m *Mystery) draw(dst *ebiten.Image) {
	m.Opts.GeoM.Reset()
	m.Opts.GeoM.Translate(m.Rect.TopLeft())
	dst.DrawImage(m.Img, m.Opts)
}

func (m *Mystery) collidePlayerBullet(b *PlayerBullet) {
	if !m.Rect.CollideRect(b.Rect) {
		return
	}
	m.hide()
	b.hitSomething()
	m.nextGoTime = 5 * time.Second
}

func (m *Mystery) isGoing() bool {
	return m.Rect.X > -50
}

func (m *Mystery) hide() {
	m.Rect.SetTopLeft(-100, MysteryY)
}
