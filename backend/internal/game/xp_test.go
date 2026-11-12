package game

import (
	"testing"
	"time"
)

func TestMatchXP(t *testing.T) {
	cases := []struct {
		name     string
		won      bool
		duration time.Duration
		want     int
	}{
		{"instant win", true, 0, WinBaseXP + WinSpeedBonusXP},
		{"win halfway", true, RoundDuration / 2, WinBaseXP + WinSpeedBonusXP/2},
		{"win at the buzzer", true, RoundDuration, WinBaseXP},
		{"win after the grace period", true, RoundDuration + 5*time.Second, WinBaseXP},
		{"negative duration counts as instant", true, -time.Second, WinBaseXP + WinSpeedBonusXP},
		{"fast loss", false, time.Minute, LossXP},
		{"slow loss", false, RoundDuration, LossXP},
	}
	for _, c := range cases {
		if got := MatchXP(c.won, c.duration); got != c.want {
			t.Errorf("%s: MatchXP(%v, %v) = %d, want %d", c.name, c.won, c.duration, got, c.want)
		}
	}
}

func TestMatchXP_FasterWinEarnsMore(t *testing.T) {
	if MatchXP(true, time.Minute) <= MatchXP(true, 4*time.Minute) {
		t.Fatal("a faster win should earn more XP than a slower one")
	}
	if MatchXP(true, RoundDuration) <= MatchXP(false, 0) {
		t.Fatal("even the slowest win should earn more than any loss")
	}
}
