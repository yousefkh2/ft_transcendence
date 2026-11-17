package game

import (
	"testing"
	"time"
)

func TestLevelFor(t *testing.T) {
	cases := []struct {
		xp                   int
		level, into, forNext int
	}{
		{0, 1, 0, 100},
		{-5, 1, 0, 100},
		{99, 1, 99, 100},
		{100, 2, 0, 200},
		{299, 2, 199, 200},
		{300, 3, 0, 300},
		{1000, 5, 0, 500},
		{4500, 10, 0, 1000},
	}
	for _, c := range cases {
		got := LevelFor(c.xp)
		if got.Level != c.level || got.XPIntoLevel != c.into || got.XPForNext != c.forNext {
			t.Errorf("LevelFor(%d) = %+v, want level %d into %d next %d", c.xp, got, c.level, c.into, c.forNext)
		}
	}
}

func achievementMap(stats PlayerStats) map[string]AchievementProgress {
	out := map[string]AchievementProgress{}
	for _, a := range Achievements(stats) {
		out[a.ID] = a
	}
	return out
}

func TestAchievements_NewPlayerHasNone(t *testing.T) {
	for _, a := range Achievements(PlayerStats{}) {
		if a.Done() {
			t.Errorf("%s should not be unlocked for a new player", a.ID)
		}
	}
}

func TestAchievements_Unlocks(t *testing.T) {
	got := achievementMap(PlayerStats{
		Rounds: 12, Wins: 4, Losses: 8, XP: 1000,
		FastestSolve: 90 * time.Second, BestWinStreak: 3, Languages: 1, PlayedWithFriend: true,
	})
	want := map[string]bool{
		"first_round": true, "first_win": true, "first_loss": true,
		"quick_solve": true, "lightning": false, "flawless": false,
		"streak_3": true, "rounds_10": true, "wins_10": false,
		"level_5": true, "level_10": false, "polyglot": false, "team_player": true,
	}
	for id, done := range want {
		if got[id].Done() != done {
			t.Errorf("%s: done=%v, want %v (%+v)", id, got[id].Done(), done, got[id])
		}
	}
	// Progress is capped at the target and shows partial progress below it.
	if got["rounds_10"].Progress != 10 || got["wins_10"].Progress != 4 || got["level_10"].Progress != 5 {
		t.Errorf("unexpected progress: %+v %+v %+v", got["rounds_10"], got["wins_10"], got["level_10"])
	}
}

func TestAchievements_SolveTimesAtTheLimit(t *testing.T) {
	got := achievementMap(PlayerStats{Wins: 1, FastestSolve: LightningTime})
	if !got["lightning"].Done() || !got["quick_solve"].Done() {
		t.Errorf("a solve exactly at the limit should count: %+v", got)
	}
}
