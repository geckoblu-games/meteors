package game

import (
	"math"

	"github.com/orgs/geckoblu-games/meteors/assets"
)

type Bullet struct {
	Sprite

	speed            Vector
	bulletAliveTimer *Timer
}

func NewBullet(p *Player) *Bullet {
	image := assets.BulletSprite

	pbounds := p.image.Bounds()
	sbounds := image.Bounds()

	rotation := p.rotation - 90
	rad := float64(rotation) * math.Pi / 180.0
	position := p.position
	position.Y += (float64(pbounds.Dy()+sbounds.Dy()) / 2) * math.Sin(rad)
	position.X += (float64(pbounds.Dy()+sbounds.Dy()) / 2) * math.Cos(rad)

	bspeed := bulletSpeed + p.speed.magnitude

	bullet := &Bullet{

		Sprite: Sprite{
			game:     p.game,
			image:    image,
			position: position,
			rotation: p.rotation,
		},

		speed: NewVector(bspeed, rotation),
	}

	bullet.bulletAliveTimer = AfterFunc(bulletAliveTime, func() {
		// fmt.Println("timer2 fired!")
		bullet.game.RemoveBullet(bullet)
	})

	return bullet

}

func (b *Bullet) Update() {
	b.position.Add(b.speed)
	b.position.KeepInbound()
	b.bulletAliveTimer.Update()
}

// func (b *Bullet) Draw(screen *ebiten.Image) {
// 	b.Sprite.Draw(screen)
// }
