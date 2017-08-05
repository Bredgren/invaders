package main

import (
	"fmt"
	_ "image/png"
	"log"
	"time"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
)

const (
	Width  = 255
	Height = 255
)

var (
	timeScale           = 1.0
	canChangeFullscreen = true
)

var (
	player       Player
	playerBullet PlayerBullet
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

	togglFullscreen()

	player.update(dt)
	playerBullet.update(dt)

	player.draw(screen)
	playerBullet.draw(screen)

	ebitenutil.DebugPrint(screen, fmt.Sprintf("FPS: %0.2f\nTime: %0.2f",
		ebiten.CurrentFPS(), timeScale))

	return nil
}

var (
	lastUpdate time.Time
)

func main() {
	player.init()
	playerBullet.init()

	if err := ebiten.Run(update, Width, Height, 2, "Invaders"); err != nil {
		log.Fatal(err)
	}
	fmt.Println("bye")
}
