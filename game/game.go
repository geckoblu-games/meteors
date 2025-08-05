package game

import (
	"fmt"
	"image/color"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/orgs/geckoblu-games/meteors/assets"
	"github.com/orgs/geckoblu-games/meteors/assets/sfx"
)

const (
	WindowWidth  = 800
	WindowHeight = 600

	TPS          = 60 // Ticks Per Second
	ScreenRatio  = 1.5
	ScreenWidth  = WindowWidth * ScreenRatio
	ScreenHeight = WindowHeight * ScreenRatio

	meteorSpawnTime        = 5 * time.Second
	meteorRotationSpeedMin = -2 // Degrees for Ticks
	meteorRotationSpeedMax = 2  // Degrees for Ticks

	shootCooldownTime = time.Duration(0.5 * float64(time.Second))
	bulletSpeed       = 120 / TPS
	bulletAliveTime   = time.Second * 3
	playerRotateSpeed = 2 // Degrees for Ticks
	playerBurstSpeed  = 0.1
	burstAliveTime    = time.Duration(0.2 * float64(time.Second))

	blinkingTiker = time.Duration(0.2 * float64(time.Second))
	blinkingTime  = time.Second * 1
)

const (
	GAMESTATUS_INTRO = iota
	GAMESTATUS_BLINKING
	GAMESTATUS_RUN
	GAMESTATUS_GAMEOVER
	GAMESTATUS_PAUSED
)

type Game struct {
	player           *Player
	meteorSpawnTimer *Timer
	meteors          []*Meteor
	bullets          []*Bullet
	scoreFace        text.Face
	status           int
	score            int
	lives            int

	blinking       bool
	blinkingTicker *Ticker
}

func NewGame() *Game {

	g := &Game{
		meteorSpawnTimer: NewTimer(meteorSpawnTime),
		scoreFace: &text.GoTextFace{
			Source: assets.KenvectorFaceSource,
			Size:   24 * ScreenRatio,
		},
		lives: 3,
	}

	g.player = NewPlayer(g)

	// g.status = GAMESTATUS_GAMEOVER

	g.blinkingTicker = NewTicker(blinkingTiker, func() {
		// fmt.Printf("Blink\n")
		g.blinking = !g.blinking
	})

	return g
}

func (g *Game) Update() error {
	if inpututil.IsKeyJustPressed(ebiten.KeyQ) {
		return ebiten.Termination
	}
	switch g.status {
	case GAMESTATUS_INTRO:
		return g.Update_Intro()
	case GAMESTATUS_BLINKING:
		return g.Update_Blinking()
	case GAMESTATUS_RUN:
		return g.Update_Run()
	case GAMESTATUS_GAMEOVER:
		return g.Update_Gameover()
	case GAMESTATUS_PAUSED:
		return g.Update_Paused()
	}

	return nil
}

func (g *Game) Update_Intro() error {
	if len(g.meteors) < 10 {
		for i := 0; i < 10-len(g.meteors); i++ {
			m := NewMeteor(g)
			m.position = NewRandomPosition()
			g.AddMeteor(m)
		}
	}
	for _, m := range g.meteors {
		m.Update()
	}
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		g.startNewGame()
	}
	return nil
}

func (g *Game) Update_Blinking() error {
	g.blinkingTicker.Update()
	return nil
}

func (g *Game) Update_Gameover() error {
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		g.startNewGame()
	}
	return nil
}

func (g *Game) Update_Paused() error {
	if inpututil.IsKeyJustPressed(ebiten.KeyP) ||
		inpututil.IsKeyJustPressed(ebiten.KeyEscape) ||
		inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
		fmt.Print("Game Unpaused\n")
		g.status = GAMESTATUS_RUN
	}
	return nil
}

func (g *Game) Update_Run() error {
	if inpututil.IsKeyJustPressed(ebiten.KeyP) ||
		inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		fmt.Print("Game Paused\n")
		g.status = GAMESTATUS_PAUSED
	}

	g.meteorSpawnTimer.Update()
	if g.meteorSpawnTimer.IsReady() {
		// Execute the spawn!
		// fmt.Printf("spawn %v!\n", len(assets.MeteorSprites))
		g.AddMeteor(NewMeteor(g))

		g.meteorSpawnTimer.Decrement(5, 100)
		g.meteorSpawnTimer.Reset()
	}

	for _, m := range g.meteors {
		m.Update()
	}

	g.player.Update()

	for _, b := range g.bullets {
		b.Update()
	}

	for _, m := range g.meteors {
		for _, b := range g.bullets {
			//if m.Collider().Intersects(b.Collider()) {
			if CheckCollision(&m.Sprite, &b.Sprite) {
				// A meteor collided with a bullet
				// fmt.Printf("A meteor collided with a bullet\n")
				g.score += m.ScoreValue()
				g.RemoveMeteor(m)
				m.Explode()
				g.RemoveBullet(b)
			}
		}
	}

	for _, m := range g.meteors {
		if CheckCollision(&m.Sprite, &g.player.Sprite) {
			// A meteor collided with the player
			fmt.Printf("A meteor collided with the player\n")
			sfx.PlayExplosion(1)
			if g.lives > 0 {
				g.startNewRound()
			} else {
				g.status = GAMESTATUS_GAMEOVER
			}
		}
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	switch g.status {
	case GAMESTATUS_INTRO:
		g.Draw_Intro(screen)
	case GAMESTATUS_BLINKING:
		g.Draw_Blinking(screen)
	case GAMESTATUS_RUN:
		g.Draw_Run(screen)
	case GAMESTATUS_PAUSED:
		g.Draw_Paused(screen)
	case GAMESTATUS_GAMEOVER:
		g.Draw_Gameover(screen)
	}
}

func (g *Game) Draw_Intro(screen *ebiten.Image) {
	for _, m := range g.meteors {
		m.Draw(screen)
	}

	titleFontSize := 60.0 * ScreenRatio
	fontSize := 20.0 * ScreenRatio
	titleTexts := "METEORS"
	// texts := "\n\n\n\n\n\nCONTROLS:\n\n←  Turn left\n→  Turn right\n\u2191  Burst\n<SPACE> Fire\n\n\n\n\nPRESS <SPACE> TO START"
	texts := "PRESS <SPACE> TO START"

	op := &text.DrawOptions{}
	op.GeoM.Translate(ScreenWidth/2, 3*titleFontSize)
	op.ColorScale.ScaleWithColor(color.White)
	op.LineSpacing = titleFontSize
	op.PrimaryAlign = text.AlignCenter
	text.Draw(screen, titleTexts, &text.GoTextFace{
		Source: assets.KenvectorFaceSource,
		Size:   titleFontSize,
	}, op)

	op = &text.DrawOptions{}
	op.GeoM.Translate(ScreenWidth/2, 700)
	op.ColorScale.ScaleWithColor(color.White)
	op.LineSpacing = fontSize
	op.PrimaryAlign = text.AlignCenter
	text.Draw(screen, texts, &text.GoTextFace{
		Source: assets.KenvectorFaceSource,
		Size:   fontSize,
	}, op)

}

func (g *Game) Draw_Blinking(screen *ebiten.Image) {

	if g.blinking {
		g.player.DrawEx(screen)
	} else {
		g.player.Draw(screen)
	}

	g.drawScoreAndLives(screen)
}

func (g *Game) Draw_Gameover(screen *ebiten.Image) {
	for _, m := range g.meteors {
		m.Draw(screen)
	}

	g.drawScoreAndLives(screen)

	titleFontSize := 60.0 * ScreenRatio
	fontSize := 20.0 * ScreenRatio
	titleTexts := "GAME OVER"
	texts := "PRESS <SPACE> TO START\nOR\n <q> TO QUIT"

	op := &text.DrawOptions{}
	op.GeoM.Translate(ScreenWidth/2, 3*titleFontSize)
	op.ColorScale.ScaleWithColor(color.White)
	op.LineSpacing = titleFontSize
	op.PrimaryAlign = text.AlignCenter
	text.Draw(screen, titleTexts, &text.GoTextFace{
		Source: assets.KenvectorFaceSource,
		Size:   titleFontSize,
	}, op)

	op = &text.DrawOptions{}
	op.GeoM.Translate(ScreenWidth/2, 500)
	op.ColorScale.ScaleWithColor(color.White)
	op.LineSpacing = fontSize
	op.PrimaryAlign = text.AlignCenter
	text.Draw(screen, texts, &text.GoTextFace{
		Source: assets.KenvectorFaceSource,
		Size:   fontSize,
	}, op)
}

func (g *Game) Draw_Paused(screen *ebiten.Image) {
	g.Draw_Run(screen)
}

func (g *Game) Draw_Run(screen *ebiten.Image) {

	for _, m := range g.meteors {
		m.Draw(screen)
	}

	for _, b := range g.bullets {
		b.Draw(screen)
	}

	g.player.Draw(screen)

	g.drawScoreAndLives(screen)
}

func (g *Game) drawScoreAndLives(screen *ebiten.Image) {
	// Draw score
	msg := fmt.Sprintf("%05d", g.score)
	opt := &text.DrawOptions{}
	opt.GeoM.Translate(100, ScreenHeight-100)
	opt.ColorScale.ScaleWithColor(color.White)
	text.Draw(screen, msg, g.scoreFace, opt)

	// Draw lives
	for i := 1; i <= g.lives; i++ {
		lx := ScreenWidth - 100 - 50*float64(i)
		opi := &ebiten.DrawImageOptions{}
		opi.GeoM.Translate(lx, ScreenHeight-90)
		screen.DrawImage(assets.PlayerLive, opi)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return ScreenWidth, ScreenHeight
}

func (g *Game) AddBullet(bullet *Bullet) {
	g.bullets = append(g.bullets, bullet)
}

func (g *Game) RemoveBullet(target *Bullet) {
	for i, b := range g.bullets {
		if b == target {
			// Remove by re-slicing (preserves order)
			g.bullets = append(g.bullets[:i], g.bullets[i+1:]...)
			return
		}
	}
}

func (g *Game) AddMeteor(meteor *Meteor) {
	g.meteors = append(g.meteors, meteor)
}

func (g *Game) RemoveMeteor(target *Meteor) {
	for i, m := range g.meteors {
		if m == target {
			// Remove by re-slicing (preserves order)
			g.meteors = append(g.meteors[:i], g.meteors[i+1:]...)
			return
		}
	}
}

func (g *Game) startNewGame() {
	fmt.Printf("Start new game\n")

	g.score = 0
	g.lives = 3
	g.meteors = g.meteors[:0]
	g.bullets = g.bullets[:0]
	g.player.Reset()

	g.startNewRound()
}

func (g *Game) startNewRound() {
	g.meteors = g.meteors[:0]
	g.bullets = g.bullets[:0]
	g.player.Reset()

	g.status = GAMESTATUS_BLINKING
	time.AfterFunc(blinkingTime, func() {
		fmt.Println("Stop blinking")
		if g.lives > 0 {
			g.lives -= 1
		}
		// player.burstVisible = false
		g.status = GAMESTATUS_RUN
	})
}
