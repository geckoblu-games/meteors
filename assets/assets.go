package assets

import (
	"embed"
	"fmt"
	"image"
	_ "image/png" // required for PNG decoding

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

// Load in memory all the files under current directory
//
//go:embed *.png *.ttf meteors/*.png
var assets embed.FS

var assetsCache = make(map[string]*ebiten.Image)

func GetMeteorImage(color, size, variant int) *ebiten.Image {
	meteorImageName := fmt.Sprintf("meteors/meteor_%d_%d_%d.png", color, size, variant)
	// fmt.Println(meteorImageName)
	image, ok := assetsCache[meteorImageName]
	if !ok {
		image = mustLoadImage(meteorImageName)
		assetsCache[meteorImageName] = image
	}
	return image
}

var PlayerSprite = mustLoadImage("player.png")
var PlayerBurst = mustLoadImage("burst.png")
var PlayerLive = mustLoadImage("playerLife.png")

var BulletSprite = mustLoadImage("bullet.png")

var KenvectorFaceSource = mustLoadFaceSource("kenvector_future.ttf")

func mustLoadImage(name string) *ebiten.Image {
	f, err := assets.Open(name)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	img, _, err := image.Decode(f)
	if err != nil {
		panic(err)
	}

	return ebiten.NewImageFromImage(img)
}

// func mustLoadImages(path string) []*ebiten.Image {
// 	matches, err := fs.Glob(assets, path)
// 	if err != nil {
// 		panic(err)
// 	}

// 	images := make([]*ebiten.Image, len(matches))
// 	for i, match := range matches {
// 		images[i] = mustLoadImage(match)
// 	}

// 	return images
// }

func mustLoadFaceSource(name string) *text.GoTextFaceSource {
	f, err := assets.Open(name)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	s, err := text.NewGoTextFaceSource(f)
	if err != nil {
		panic(err)
	}

	return s
}

func WindowIcons() []image.Image {
	icon16 := mustLoadImage("icon16.png")
	icon32 := mustLoadImage("icon32.png")
	icon48 := mustLoadImage("icon48.png")

	return []image.Image{icon16, icon32, icon48}
}
