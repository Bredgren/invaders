package main

import (
	"fmt"
	"image/color"
	_ "image/png"
	"log"
	"time"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
)

const (
	Width  = 256
	Height = 256
)

var (
	timeScale           = 1.0
	canChangeFullscreen = true
	floor               *ebiten.Image
)

var (
	player       Player
	playerBullet PlayerBullet
	shelters     = [4]Shelter{}
	mystery      Mystery
)

func togglFullscreen() {
	if ebiten.IsKeyPressed(ebiten.KeyF) {
		if canChangeFullscreen {
			ebiten.SetFullscreen(!ebiten.IsFullscreen())
			canChangeFullscreen = false
		}
	} else {
		canChangeFullscreen = true
	}
}

var (
	lastUpdate time.Time
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
	dt := time.Duration(float64(now.Sub(lastUpdate).Nanoseconds())*timeScale) * time.Nanosecond
	lastUpdate = now

	if ebiten.IsKeyPressed(ebiten.KeyR) {
		reset()
	}

	togglFullscreen()

	if ebiten.IsRunningSlowly() {
		log.Println("slow")
		return nil
	}

	// Update things
	mystery.update(dt)
	player.update(dt)
	playerBullet.update(dt)

	// Check collisons
	mystery.collidePlayerBullet(&playerBullet)
	for _, shelter := range shelters {
		shelter.collidePlayerBullet(&playerBullet)
	}

	// Draw
	for _, shelter := range shelters {
		shelter.draw(screen)
	}
	mystery.draw(screen)
	player.draw(screen)
	playerBullet.draw(screen)

	drawFloor(screen)

	ebitenutil.DebugPrint(screen, fmt.Sprintf("FPS: %0.2f\nTime: %0.2f",
		ebiten.CurrentFPS(), timeScale))

	return nil
}

func drawFloor(dst *ebiten.Image) {
	opts := ebiten.DrawImageOptions{}
	opts.GeoM.Translate(0, Height-(player.Rect.H+12))
	dst.DrawImage(floor, &opts)
}

func reset() {
	player.init()
	playerBullet.init()
	mystery.init()

	for i := 0; i < len(shelters); i++ {
		shelters[i].init(i)
	}
}

func main() {
	lastUpdate = time.Now()
	reset()

	floor, _ = ebiten.NewImage(Width, 2, ebiten.FilterNearest)
	floor.Fill(color.NRGBA{0x00, 0xff, 0x00, 0xff})

	if err := ebiten.Run(update, Width, Height, 2, "Invaders"); err != nil {
		log.Fatal(err)
	}
	fmt.Println("bye")
}
