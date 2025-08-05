package game

import "github.com/hajimehoshi/ebiten/v2"

var alphaMaskCache = map[*ebiten.Image]*AlphaMask{}

type AlphaMask struct {
	Width, Height int
	Data          []bool
}

func GetAlphaMask(img *ebiten.Image) *AlphaMask {
	if mask, ok := alphaMaskCache[img]; ok {
		return mask
	}
	mask := NewAlphaMask(img)
	alphaMaskCache[img] = mask
	return mask
}

func NewAlphaMask(img *ebiten.Image) *AlphaMask {
	bounds := img.Bounds()
	w := bounds.Dx()
	h := bounds.Dy()
	mask := &AlphaMask{
		Width:  w,
		Height: h,
		Data:   make([]bool, w*h),
	}
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			_, _, _, a := img.At(x, y).RGBA()
			mask.Data[y*w+x] = a > 0
		}
	}
	return mask
}

func (m *AlphaMask) At(x, y int) bool {
	if x < 0 || y < 0 || x >= m.Width || y >= m.Height {
		return false
	}
	return m.Data[y*m.Width+x]
}
