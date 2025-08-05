package game

import (
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/colorm"
)

type Sprite struct {
	game  *Game
	image *ebiten.Image

	position Position
	rotation int // degrees
}

func (s *Sprite) Draw(screen *ebiten.Image) {

	bounds := s.image.Bounds()
	halfW := float64(bounds.Dx()) / 2
	halfH := float64(bounds.Dy()) / 2

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(-halfW, -halfH)
	rad := float64(s.rotation) * math.Pi / 180.0
	op.GeoM.Rotate(rad)

	op.GeoM.Translate(s.position.X, s.position.Y)

	screen.DrawImage(s.image, op)
}

func (s *Sprite) DrawEx(screen *ebiten.Image) {

	bounds := s.image.Bounds()
	halfW := float64(bounds.Dx()) / 2
	halfH := float64(bounds.Dy()) / 2

	op := &colorm.DrawImageOptions{}
	op.GeoM.Translate(-halfW, -halfH)
	rad := float64(s.rotation) * math.Pi / 180.0
	op.GeoM.Rotate(rad)

	op.GeoM.Translate(s.position.X, s.position.Y)

	cm := colorm.ColorM{}
	// cm.Scale(-1, -1, -1, 1)
	cm.Translate(0, 0, 1, 0)

	colorm.DrawImage(screen, s.image, cm, op)
}

func (s *Sprite) GetAlphaMask() *AlphaMask {
	return GetAlphaMask(s.image)
}
