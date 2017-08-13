package main

import (
	"log"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
)

func openImg(name string) *ebiten.Image {
	img, _, err := ebitenutil.NewImageFromFile(name, ebiten.FilterNearest)
	if err != nil {
		log.Fatalf("open %s file: %s\n", name, err)
	}
	return img
}

func toFloat(a, b int) (float64, float64) {
	return float64(a), float64(b)
}
