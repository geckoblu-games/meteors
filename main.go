package main

import (
	"github.com/hajimehoshi/ebiten/v2"

	"github.com/orgs/geckoblu-games/meteors/assets"
	"github.com/orgs/geckoblu-games/meteors/game"
)

func main() {
	// deprecated
	// rand.Seed(time.Now().UnixNano())

	g := game.NewGame()

	ebiten.SetWindowSize(game.WindowWidth, game.WindowHeight)
	ebiten.SetWindowTitle("Meteors")
	ebiten.SetCursorMode(ebiten.CursorModeHidden)
	ebiten.SetWindowIcon(assets.WindowIcons())

	err := ebiten.RunGame(g)
	if err != nil {
		panic(err)
	}
}
