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
	Width  = 320
	Height = 240
)

var (
	timeScale = 1.0
)

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
