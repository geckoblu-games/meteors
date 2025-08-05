package game

import (
	"testing"
)

func Test_Add_Vector(t *testing.T) {
	pos := Position{10, 10}
	vec := NewVector(10, 0)

	t.Logf("Position: %v\n", pos)
	pos.Add(vec)
	t.Logf("Position: %v\n", pos)
}
