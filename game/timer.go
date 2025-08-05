package game

import (
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

type Timer struct {
	currentTicks int
	targetTicks  int
	expired      bool
	timerFn      func()
}

func NewTimer(d time.Duration) *Timer {
	return &Timer{
		currentTicks: 0,
		targetTicks:  int(d.Milliseconds()) * ebiten.TPS() / 1000,
		expired:      false,
		timerFn:      nil,
	}
}

func AfterFunc(d time.Duration, timerFn func()) *Timer {
	return &Timer{
		currentTicks: 0,
		targetTicks:  int(d.Milliseconds()) * ebiten.TPS() / 1000,
		timerFn:      timerFn,
	}
}

func (t *Timer) Update() {
	if t.currentTicks < t.targetTicks {
		t.currentTicks++
	} else {
		if !t.expired && t.timerFn != nil {
			t.timerFn()
		}
		t.expired = true
	}
}

func (t *Timer) IsReady() bool {
	return t.expired
}

func (t *Timer) Reset() {
	t.currentTicks = 0
	t.expired = false
}

func (t *Timer) Decrement(decrement, limit int) {
	t.targetTicks -= decrement
	if t.targetTicks < limit {
		t.targetTicks = limit
	}
	// fmt.Printf("Timer.Decrement %v\n", t.targetTicks)
}

type Ticker struct {
	currentTicks int
	targetTicks  int
	tickerFn     func()
}

func NewTicker(d time.Duration, tickerFn func()) *Ticker {
	return &Ticker{
		currentTicks: 0,
		targetTicks:  int(d.Milliseconds()) * ebiten.TPS() / 1000,
		tickerFn:     tickerFn,
	}
}

func (t *Ticker) Update() {
	if t.currentTicks < t.targetTicks {
		t.currentTicks++
	} else {
		t.currentTicks = 0
		t.tickerFn()
	}
}
