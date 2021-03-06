package main

import (
	"fmt"
	"image/color"
	_ "image/png"
	"log"
	"time"

	"golang.org/x/image/font/basicfont"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/text"
)

const (
	Width  = 256
	Height = 256
)

var (
	timeScale           = 1.0
	canChangeFullscreen = true
	floor               *ebiten.Image
	fontFace            = basicfont.Face7x13
	score               = 0
	highscore           = 0
	level               = 0
	nextLife            = 2000
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
	// if ebiten.IsKeyPressed(ebiten.KeyZ) {
	// 	timeScale = 0.25
	// } else if ebiten.IsKeyPressed(ebiten.KeyX) {
	// 	timeScale = 8.0
	// } else {
	// 	timeScale = 1.0
	// }
	dt := time.Duration(float64(now.Sub(lastUpdate).Nanoseconds())*timeScale) * time.Nanosecond
	lastUpdate = now

	if /*ebiten.IsKeyPressed(ebiten.KeyR) ||*/ len(aliens.activeAliens()) == 0 {
		level += 1
		resetLevel()
	}

	togglFullscreen()

	// Update things
	mystery.update(dt)
	player.update(dt)
	playerBullet.update(dt)
	aliens.update(dt)
	missiles.update(dt)

	// Check collisons
	mystery.collidePlayerBullet(&playerBullet)
	for i := range shelters {
		shelters[i].collidePlayerBullet(&playerBullet)
		shelters[i].collideMissiles()
	}
	aliens.collidePlayerBullet(&playerBullet)
	for i := range missiles.Missiles {
		if playerBullet.isGoing() && missiles.Missiles[i].Top() > 0 {
			if missiles.Missiles[i].CollideRect(playerBullet.Rect) {
				playerBullet.hitSomething()
				missiles.hitSomething(i)
			}
		}
	}
	player.collideEnemyMissile()
	aliens.collideShelters()
	aliens.collidePlayer()

	if score > highscore {
		highscore = score
	}

	if player.Lives <= 0 {
		level = 0
		resetLevel()
		return nil
	}

	// If we happen to gain an extra life in the same frame that we lose our last one,
	// just end the game to avoid confusion.
	if score >= nextLife {
		nextLife *= 2
		if player.Lives < PlayerMaxLives {
			player.Lives += 1
		}
	}

	if ebiten.IsRunningSlowly() {
		// log.Println("slow")
		return nil
	}

	// Draw
	for i := range shelters {
		shelters[i].draw(screen)
	}
	mystery.draw(screen)
	player.draw(screen)
	aliens.draw(screen)
	playerBullet.draw(screen)
	missiles.draw(screen)

	drawFloor(screen)

	drawPlayerLives(screen)

	// ebitenutil.DebugPrint(screen, fmt.Sprintf("FPS: %0.2f\nTime: %0.2f",
	// 	ebiten.CurrentFPS(), timeScale))

	text.Draw(screen, fmt.Sprintf("High Score %d", highscore), fontFace, 5, 15, color.White)
	text.Draw(screen, fmt.Sprintf("Score %d", score), fontFace, 5, 30, color.White)
	text.Draw(screen, fmt.Sprintf("Level %d", level), fontFace, 5, Height-5, color.White)

	return nil
}

func drawFloor(dst *ebiten.Image) {
	opts := ebiten.DrawImageOptions{}
	opts.GeoM.Translate(0, Height-(player.Rect.H+12))
	dst.DrawImage(floor, &opts)
}

func drawPlayerLives(dst *ebiten.Image) {
	opts := ebiten.DrawImageOptions{}
	opts.ColorM.Scale(0.0, 1.0, 0.0, 1.0)
	w, h := player.Rect.Size()
	padding := 6.0
	opts.GeoM.Translate(Width, Height-h-padding)
	for i := 0; i < player.Lives; i++ {
		opts.GeoM.Translate(-w-padding, 0)
		dst.DrawImage(player.Img, &opts)
	}
}

func resetLevel() {
	if level == 0 {
		score = 0
		nextLife = 2000
	}

	player.resetLevel(level)
	playerBullet.resetLevel(level)
	mystery.resetLevel(level)
	aliens.resetLevel(level)
	missiles.resetLevel(level)

	for i := range shelters {
		shelters[i].resetLevel(level)
	}
}

func init() {
	lastUpdate = time.Now()

	player.init()
	playerBullet.init()
	mystery.init()
	aliens.init()
	missiles.init()

	for i := range shelters {
		shelters[i].init(i)
	}
}

func main() {
	resetLevel()

	floor, _ = ebiten.NewImage(Width, 2, ebiten.FilterNearest)
	floor.Fill(color.NRGBA{0x00, 0xff, 0x00, 0xff})

	if err := ebiten.Run(update, Width, Height, 2, "Invaders"); err != nil {
		log.Fatal(err)
	}
	fmt.Println("bye")
}
