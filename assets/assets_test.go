package assets

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_GetMeteorImage(t *testing.T) {
	for color := 1; color < 3; color++ {
		for size := 1; size < 4; size++ {
			if size == 1 {
				for variant := 1; variant < 5; variant++ {
					GetMeteorImage(color, size, variant)
				}
			} else {
				for variant := 1; variant < 3; variant++ {
					image := GetMeteorImage(color, size, variant)
					assert.NotNil(t, image)
				}
			}
		}
	}

	image1 := GetMeteorImage(1, 1, 1)
	image2 := GetMeteorImage(1, 1, 1)
	assert.Equal(t, image1, image2)
}
