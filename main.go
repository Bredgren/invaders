package main

import (
	"fmt"
	_ "image/png"
	"log"
	"time"

	"github.com/Bredgren/geo"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
)

const (
	Width  = 320
	Height = 240
)

const (
	PlayerSpeed = 75
	PlayerY     = 20
)

var (
	timeScale = 1.0
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

var (
	player Player
)

func update(screen *ebiten.Image) error {
	now := time.Now()
	if ebiten.IsKeyPressed(ebiten.KeyZ) {
		timeScale = 0.5
	} else if ebiten.IsKeyPressed(ebiten.KeyX) {
		timeScale = 2.0
	} else {
		timeScale = 1.0
	}
	dt := time.Duration(float64(lastUpdate.Sub(now).Nanoseconds())*timeScale) * time.Nanosecond
	lastUpdate = now

	player.move(dt)

	player.draw(screen)
	ebitenutil.DebugPrint(screen, fmt.Sprintf("FPS: %0.2f\nTime: %0.2f",
		ebiten.CurrentFPS(), timeScale))
	return nil
}

var (
	lastUpdate time.Time
)

func main() {
	player.init()

	if err := ebiten.Run(update, Width, Height, 2, "Invaders"); err != nil {
		log.Fatal(err)
	}
	fmt.Println("bye")
}
