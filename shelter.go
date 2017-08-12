package main

import (
	"image/color"
	"math/rand"

	"github.com/Bredgren/geo"
	"github.com/hajimehoshi/ebiten"
)

const (
	ShelterBottomY = 210
	ShelterW       = 20
	ShelterH       = 15
)

const (
	unused   = 0
	inactive = 1
	active   = 2
)

var (
	shelters = [4]Shelter{}
	ShelterX = [4]float64{0.2, 0.4, 0.6, 0.8}
)

type Shelter struct {
	Rect      geo.Rect
	subRects  [ShelterW * ShelterH]geo.Rect
	subStates [ShelterW * ShelterH]int
	Img       *ebiten.Image
	Opts      *ebiten.DrawImageOptions
}

func (s *Shelter) init(num int) {
	s.Img, _ = ebiten.NewImage(1, 1, ebiten.FilterNearest)
	s.Img.Fill(color.NRGBA{0x00, 0xff, 0x00, 0xff})
	s.Opts = &ebiten.DrawImageOptions{}

	s.Rect = geo.RectWH(ShelterW, ShelterH)
	s.Rect.SetBottomMid(ShelterX[num]*Width, ShelterBottomY)

	s.initSubRects()
}

func (s *Shelter) resetLevel(level int) {

}

func (s *Shelter) initSubRects() {
	for y := 0; y < ShelterH; y++ {
		for x := 0; x < ShelterW; x++ {
			if (y == 0 && (x < 4 || x > 15)) ||
				y == 1 && (x < 3 || x > 16) ||
				y == 2 && (x < 2 || x > 17) ||
				y == 3 && (x < 1 || x > 18) ||
				y == 12 && (x > 6 && x < 13) ||
				y == 13 && (x > 5 && x < 14) ||
				y == 14 && (x > 4 && x < 15) {
				continue
			}
			i := y*ShelterW + x
			s.subRects[i] = geo.RectXYWH(s.Rect.X+float64(x), s.Rect.Y+float64(y), 1, 1)
			s.subStates[i] = active
		}
	}
}

func (s *Shelter) draw(dst *ebiten.Image) {
	for i, subRect := range s.subRects {
		if s.subStates[i] != active {
			continue
		}
		s.Opts.GeoM.Reset()
		s.Opts.GeoM.Translate(subRect.TopLeft())
		dst.DrawImage(s.Img, s.Opts)
	}
}

func (s *Shelter) collidePlayerBullet(b *PlayerBullet) {
	if !s.Rect.CollideRect(b.Rect) {
		return
	}

	activeRects := make([]geo.Rect, 0, len(s.subRects))
	for i := range s.subStates {
		if s.subStates[i] == active {
			activeRects = append(activeRects, s.subRects[i])
		}
	}

	if _, collide := b.Rect.CollideRectList(activeRects); !collide {
		return
	}

	explosionArea := geo.CircleXYR(b.Rect.MidX(), b.Rect.MidY(), 3)

	hit := explosionArea.CollideRectListAll(s.subRects[:])
	activeHits := make([]int, 0, len(hit))
	for _, i := range hit {
		if s.subStates[i] == active {
			activeHits = append(activeHits, i)
		}
	}

	// Remove 5/8 of the rects in the explosion radius
	countToRemove := len(activeHits) * 5 / 8
	if countToRemove < 2 {
		countToRemove = len(activeHits)
	}

	// "shuffle" the activeHits slice and pull out countToRemove indices from it
	order := rand.Perm(len(activeHits))
	for i := 0; i < countToRemove; i++ {
		iToRemove := activeHits[order[i]]
		s.subStates[iToRemove] = inactive
	}

	b.hitSomething()
}
