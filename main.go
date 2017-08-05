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

type Player struct {
	Pos   float64
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
	p.Pos = Width / 2
}

func (p *Player) move(dt time.Duration) {
	vel := 0.0
	if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		vel = PlayerSpeed
	}
	if ebiten.IsKeyPressed(ebiten.KeyRight) {
		vel = -PlayerSpeed
	}
	p.Pos += vel * dt.Seconds()

	// Keep on screen
	w, _ := p.Img.Size()
	p.Pos = geo.Clamp(p.Pos, 0, float64(Width-w))
}

func (p *Player) draw(dst *ebiten.Image) {
	p.Opts.GeoM.Reset()
	p.Opts.GeoM.Translate(p.Pos, Height-PlayerY)
	dst.DrawImage(p.Img, p.Opts)
}

var (
	player Player
)

func update(screen *ebiten.Image) error {
	now := time.Now()
	dt := lastUpdate.Sub(now)
	lastUpdate = now

	player.move(dt)

	player.draw(screen)
	ebitenutil.DebugPrint(screen, fmt.Sprintf("FPS: %0.2f", ebiten.CurrentFPS()))
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
