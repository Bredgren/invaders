package main

import (
	"image/color"
	"math/rand"
	"sort"

	"github.com/Bredgren/geo"
	"github.com/hajimehoshi/ebiten"
)

const (
	ShelterBottomY = 210
	ShelterW       = 20
	ShelterH       = 15
)

var (
	ShelterX = [4]float64{0.2, 0.4, 0.6, 0.8}
)

type Shelter struct {
	Rect     geo.Rect
	subRects []geo.Rect
	Img      *ebiten.Image
	Opts     *ebiten.DrawImageOptions
}

func (s *Shelter) init(num int) {
	s.Img, _ = ebiten.NewImage(1, 1, ebiten.FilterNearest)
	s.Img.Fill(color.NRGBA{0x00, 0xff, 0x00, 0xff})
	s.Opts = &ebiten.DrawImageOptions{}

	s.Rect = geo.RectWH(ShelterW, ShelterH)
	s.Rect.SetBottomMid(ShelterX[num]*Width, ShelterBottomY)

	s.subRects = shelterPix(s.Rect)
}

func shelterPix(r geo.Rect) []geo.Rect {
	res := make([]geo.Rect, 0, ShelterW*ShelterH)
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
			res = append(res, geo.RectXYWH(r.X+float64(x), r.Y+float64(y), 1, 1))
		}
	}
	return res
}

func (s *Shelter) draw(dst *ebiten.Image) {
	for _, subRect := range s.subRects {
		s.Opts.GeoM.Reset()
		s.Opts.GeoM.Translate(subRect.TopLeft())
		dst.DrawImage(s.Img, s.Opts)
	}
}

func (s *Shelter) collidePlayerBullet(b *PlayerBullet) {
	if !s.Rect.CollideRect(b.Rect) {
		return
	}

	if _, collide := b.Rect.CollideRectList(s.subRects); !collide {
		return
	}

	explosionArea := geo.CircleXYR(b.Rect.MidX(), b.Rect.MidY(), 3)

	hit := explosionArea.CollideRectListAll(s.subRects)
	// Remove 5/8 of the rects in the explosion radius
	countToRemove := len(hit) * 5 / 8
	if countToRemove < 2 {
		countToRemove = len(hit)
	}
	// "shuffle" the hit slice and pull out countToRemove indices from it
	order := rand.Perm(len(hit))
	toRemove := make([]int, countToRemove)
	for i := 0; i < countToRemove; i++ {
		toRemove[i] = hit[order[i]]
	}

	// Remove the items from s.subRects, backwards so that the indicies are correct
	sort.Sort(sort.Reverse(sort.IntSlice(toRemove)))
	for _, iToRemove := range toRemove {
		s.subRects = append(s.subRects[:iToRemove], s.subRects[iToRemove+1:]...)
	}

	b.hitSomething()
}
