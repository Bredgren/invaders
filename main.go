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
	Width  = 256
	Height = 256
)

var (
	timeScale           = 1.0
	canChangeFullscreen = true
)

var (
	player       Player
	playerBullet PlayerBullet
	shelters     = [4]Shelter{}
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

	if ebiten.IsKeyPressed(ebiten.KeyR) {
		reset()
	}

	togglFullscreen()

	player.update(dt)
	playerBullet.update(dt)

	for _, shelter := range shelters {
		shelter.collidePlayerBullet(&playerBullet)
	}

	for _, shelter := range shelters {
		shelter.draw(screen)
	}

	player.draw(screen)
	playerBullet.draw(screen)

	ebitenutil.DebugPrint(screen, fmt.Sprintf("FPS: %0.2f\nTime: %0.2f",
		ebiten.CurrentFPS(), timeScale))

	return nil
}

var (
	lastUpdate time.Time
)

func reset() {
	player.init()
	playerBullet.init()

	for i := 0; i < len(shelters); i++ {
		shelters[i].init(i)
	}
}

func main() {
	reset()

	if err := ebiten.Run(update, Width, Height, 2, "Invaders"); err != nil {
		log.Fatal(err)
	}
	fmt.Println("bye")
}
