package sfx

import (
	"embed"

	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/wav"
)

const sampleRate = 11025

//go:embed *.wav
var sounds embed.FS

var audioContext = audio.NewContext(sampleRate)

func mustLoadSound(name string) *audio.Player {
	data, err := sounds.Open(name)
	if err != nil {
		panic(err)
	}
	defer data.Close()

	stream, err := wav.DecodeWithoutResampling(data)
	if err != nil {
		panic(err)
	}

	// Takes the stream and turns it into a player
	singlePlayer, err := audioContext.NewPlayer(stream)
	if err != nil {
		panic(err)
	}

	return singlePlayer
}

func playSound(name string) {
	player := mustLoadSound(name)
	player.SetVolume(0.1)
	player.Play()
}

func PlayShootSound() {
	playSound("fire.wav")
}

func PlayBurstSound() {
	playSound("thrust.wav")
}

func PlayExplosion(size int) {
	switch size {
	case 1:
		playSound("bangLarge.wav")
	case 2:
		playSound("bangMedium.wav")
	default:
		playSound("bangSmall.wav")
	}
}
