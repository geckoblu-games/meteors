package game

import (
	"math"
	"math/rand"
)

type Vector struct {
	magnitude float64
	direction int // degrees
}

func NewVector(magnitude float64, direction int) Vector {
	return Vector{
		magnitude: magnitude,
		direction: direction,
	}
}

func NewRandomDirectionVector(magnitude float64) Vector {
	return Vector{
		magnitude: magnitude,
		direction: int(rand.Float64() * 360),
	}
}

func (v Vector) X() float64 {
	rad := float64(v.direction) * math.Pi / 180.0
	return v.magnitude * math.Cos(rad)
}

func (v Vector) Y() float64 {
	rad := float64(v.direction) * math.Pi / 180.0
	return v.magnitude * math.Sin(rad)
}

func (v *Vector) Add(v2 Vector) {
	v.AddScalar(v2.magnitude, v2.direction)
}

func (v *Vector) AddScalar(magnitude float64, direction int) {
	x1, y1 := v.X(), v.Y()

	rad := float64(direction) * math.Pi / 180.0
	x2 := magnitude * math.Cos(rad)
	y2 := magnitude * math.Sin(rad)

	x := x1 + x2
	y := y1 + y2

	m := math.Sqrt(x*x + y*y)

	angleRad := math.Atan2(y, x)           // radians
	angleDeg := angleRad * 180.0 / math.Pi // degrees

	v.magnitude = m
	v.direction = int(angleDeg)
}
