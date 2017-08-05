package main

import (
	"fmt"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
)

const (
	Width  = 320
	Height = 240
)

func update(screen *ebiten.Image) error {
	ebitenutil.DebugPrint(screen, fmt.Sprintf("FPS: %0.2f", ebiten.CurrentFPS()))
	return nil
}

func main() {
	if err := ebiten.Run(update, Width, Height, 2, "Invaders"); err != nil {
		panic(err)
	}
	fmt.Println("bye")
}
