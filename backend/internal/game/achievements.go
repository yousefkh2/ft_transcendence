package game

import "time"

// PlayerStats is everything the achievements are based on. It is computed from the player's
// completed rounds, so achievements can be checked again at any time (also for old rounds).
type PlayerStats struct {
	Rounds int
	Wins   int
	Losses int
	XP     int
	// FastestSolve is the shortest solved round; zero when the player has never solved one.
	FastestSolve time.Duration
	// FlawlessSolve is true when a round was solved without a single wrong check.
	FlawlessSolve bool
	// BestWinStreak is the most solved rounds in a row.
	BestWinStreak int
	// Languages is how many different languages the player has completed rounds in.
	Languages int
	// PlayedWithFriend is true when a completed round was played together with a friend.
	PlayedWithFriend bool
}

// AchievementProgress is one achievement with the player's progress towards it.
type AchievementProgress struct {
	ID       string `json:"id"`
	Progress int    `json:"progress"`
	Target   int    `json:"target"`
}

// Done reports whether the player has reached the target.
func (a AchievementProgress) Done() bool {
	return a.Progress >= a.Target
}

const (
	QuickSolveTime = 2 * time.Minute
	LightningTime  = time.Minute
)

func boolProgress(ok bool) int {
	if ok {
		return 1
	}
	return 0
}

func capped(value, target int) int {
	if value > target {
		return target
	}
	return value
}

func solvedWithin(stats PlayerStats, limit time.Duration) bool {
	return stats.FastestSolve > 0 && stats.FastestSolve <= limit
}

// Achievements lists every achievement in display order with the player's progress.
// The ids are stable: they are stored in user_achievements and translated in the frontend.
func Achievements(stats PlayerStats) []AchievementProgress {
	level := LevelFor(stats.XP).Level
	list := []struct {
		id       string
		progress int
		target   int
	}{
		{"first_round", stats.Rounds, 1},
		{"first_win", stats.Wins, 1},
		{"first_loss", stats.Losses, 1},
		{"quick_solve", boolProgress(solvedWithin(stats, QuickSolveTime)), 1},
		{"lightning", boolProgress(solvedWithin(stats, LightningTime)), 1},
		{"flawless", boolProgress(stats.FlawlessSolve), 1},
		{"streak_3", stats.BestWinStreak, 3},
		{"rounds_10", stats.Rounds, 10},
		{"wins_10", stats.Wins, 10},
		{"level_5", level, 5},
		{"level_10", level, 10},
		{"polyglot", stats.Languages, 2},
		{"team_player", boolProgress(stats.PlayedWithFriend), 1},
	}
	out := make([]AchievementProgress, len(list))
	for i, a := range list {
		out[i] = AchievementProgress{ID: a.id, Progress: capped(a.progress, a.target), Target: a.target}
	}
	return out
}
