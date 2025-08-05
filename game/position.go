package game

import "math/rand"

type Position struct {
	X float64
	Y float64
}

func (p *Position) Add(v Vector) {
	p.X += v.X()
	p.Y += v.Y()
}

func (p *Position) KeepInbound() {
	if p.X < 0 {
		p.X += ScreenWidth
	}
	if p.X > ScreenWidth {
		p.X -= ScreenWidth
	}
	if p.Y < 0 {
		p.Y += ScreenHeight
	}
	if p.Y > ScreenHeight {
		p.Y -= ScreenHeight
	}
}

func (p *Position) CheckInbound() bool {
	if p.X < 0 || p.Y < 0 {
		return false
	}
	if p.X > ScreenWidth {
		return false
	}
	if p.Y > ScreenHeight {
		return false
	}
	return true
}

func NewRandomPosition() Position {
	return Position{
		X: rand.Float64() * ScreenWidth,
		Y: rand.Float64() * ScreenHeight,
	}
}
