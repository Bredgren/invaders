package main

import (
	"log"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
)

func openImg(name string) *ebiten.Image {
	img, _, err := ebitenutil.NewImageFromFile("img/mystery.png", ebiten.FilterNearest)
	if err != nil {
		log.Fatalf("open %s file: %s\n", name, err)
	}
	return img
}
