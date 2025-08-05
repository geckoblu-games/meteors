package game

import (
	"fmt"
	"math"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/orgs/geckoblu-games/meteors/assets"
	"github.com/orgs/geckoblu-games/meteors/assets/sfx"
)

type Player struct {
	Sprite

	speed Vector

	burstImage        *ebiten.Image
	burstVisible      bool
	burstVisibleTimer *Timer

	shootCooldownTimer *Timer
}

func NewPlayer(game *Game) *Player {

	player := Player{

		Sprite: Sprite{
			game:  game,
			image: assets.PlayerSprite,
		},

		burstImage: assets.PlayerBurst,

		shootCooldownTimer: NewTimer(shootCooldownTime),
	}

	player.burstVisibleTimer = AfterFunc(burstAliveTime, func() {
		fmt.Println("Stop burst!")
		player.burstVisible = false
	})

	player.Reset()

	return &player
}

func (p *Player) Reset() {
	pos := Position{
		X: ScreenWidth / 2,
		Y: ScreenHeight / 2,
	}

	p.position = pos
	p.rotation = 0
	p.speed = Vector{}
	p.burstVisible = false
}

func (p *Player) Update() {
	p.shootCooldownTimer.Update()
	p.burstVisibleTimer.Update()

	if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		p.rotation -= playerRotateSpeed
		// fmt.Printf("Player.rotation: %v\n", p.rotation)
	}
	if ebiten.IsKeyPressed(ebiten.KeyRight) {
		p.rotation += playerRotateSpeed
		// fmt.Printf("Player.rotation: %v\n", p.rotation)
	}
	if ebiten.IsKeyPressed(ebiten.KeyUp) {
		p.speed.AddScalar(playerBurstSpeed, p.rotation-90)
		// fmt.Printf("Accellerate: %v\n", p.speed)
		p.burstVisible = true
		p.burstVisibleTimer.Reset()
	}

	if p.shootCooldownTimer.IsReady() && ebiten.IsKeyPressed(ebiten.KeySpace) {
		p.shoot()
	}

	p.position.Add(p.speed)
	p.position.KeepInbound()
}

func (p *Player) Draw(screen *ebiten.Image) {
	p.Sprite.Draw(screen)
	p.drawBurst(screen)
}

func (p *Player) drawBurst(screen *ebiten.Image) {
	if p.burstVisible {
		bounds := p.burstImage.Bounds()
		halfW := float64(bounds.Dx()) / 2
		halfH := float64(bounds.Dy()) / 2

		pbounds := p.image.Bounds()

		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(-halfW, -halfH+float64(pbounds.Dy()))
		rad := float64(p.rotation) * math.Pi / 180.0
		op.GeoM.Rotate(rad)

		op.GeoM.Translate(p.position.X, p.position.Y)

		screen.DrawImage(p.burstImage, op)
		sfx.PlayBurstSound()
	}
}

func (p *Player) shoot() {
	// fmt.Print("Shoot\n")
	p.game.AddBullet(NewBullet(p))
	p.shootCooldownTimer.Reset()
	sfx.PlayShootSound()
}
