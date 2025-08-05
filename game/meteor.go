package game

import (
	"math"
	"math/rand"

	"github.com/orgs/geckoblu-games/meteors/assets"
	"github.com/orgs/geckoblu-games/meteors/assets/sfx"
)

type Meteor struct {
	Sprite

	speed         Vector
	rotationSpeed int // Degrees for Ticks

	mcolor   int // 1: brown, 2 grey
	msize    int // 1: big, 2: med, 3: small
	mvariant int // 1 or 2 [3 or 4 only for big]
}

func NewMeteor(game *Game) *Meteor {

	pos := randomPosition(game)

	rotationSpeed := meteorRotationSpeedMin + int(rand.Float64()*(meteorRotationSpeedMax-meteorRotationSpeedMin))
	if rotationSpeed == 0 {
		// reduce the possibility to have a 0 rotationSpeed
		rotationSpeed = meteorRotationSpeedMin + int(rand.Float64()*(meteorRotationSpeedMax-meteorRotationSpeedMin))
	}

	mcolor, msize, mvariant := randomColorSizeVariant()
	image := assets.GetMeteorImage(mcolor, msize, mvariant)

	return &Meteor{
		Sprite: Sprite{
			game:     game,
			image:    image,
			position: pos,
		},

		speed:         NewRandomDirectionVector(1),
		rotationSpeed: rotationSpeed,

		mcolor:   mcolor,
		msize:    msize,
		mvariant: mvariant,
	}
}

func NewHalfMeteor(original *Meteor) *Meteor {

	mcolor := original.mcolor
	msize := original.msize + 1
	mvariant := rand.Intn(2) + 1

	image := assets.GetMeteorImage(mcolor, msize, mvariant)

	return &Meteor{
		Sprite: Sprite{
			game:     original.game,
			image:    image,
			position: original.position,
		},

		speed:         NewRandomDirectionVector(1),
		rotationSpeed: original.rotationSpeed,

		mcolor:   mcolor,
		msize:    msize,
		mvariant: mvariant,
	}
}

func randomPosition(game *Game) Position {
	targetX := game.player.position.X
	targetY := game.player.position.Y

	angle := rand.Float64() * 2 * math.Pi
	r := ScreenWidth / 2.0

	pos := Position{
		X: targetX + math.Cos(angle)*r,
		Y: targetY + math.Sin(angle)*r,
	}

	pos.KeepInbound()

	return pos
}

func randomColorSizeVariant() (int, int, int) {

	// To have a not uniform distribution of sizes
	sizes := []int{1, 1, 1, 1, 2, 2, 3}

	mcolor := rand.Intn(2) + 1
	msize := sizes[rand.Intn(len(sizes))]
	var mvariant int
	if msize == 1 {
		mvariant = rand.Intn(4) + 1
	} else {
		mvariant = rand.Intn(2) + 1
	}

	return mcolor, msize, mvariant
}

func (m *Meteor) Update() {

	m.position.Add(m.speed)
	m.position.KeepInbound()

	m.rotation += m.rotationSpeed
}

// func (m *Meteor) Draw(screen *ebiten.Image) {
// 	m.Sprite.Draw(screen)
// }

func (m *Meteor) Explode() {
	if m.msize < 3 { // small just disappear
		// m.game.RemoveMeteor(m)
		m1 := NewHalfMeteor(m)
		m.game.AddMeteor(m1)
		m2 := NewHalfMeteor(m)
		m.game.AddMeteor(m2)
	}
	sfx.PlayExplosion(m.msize)
	m.game.RemoveMeteor(m)
}

func (m *Meteor) ScoreValue() int {
	return 5
}
