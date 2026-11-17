package game

// Levels grow with XP. Each level needs 100 XP more than the one before:
// level 2 at 100 XP, level 3 at 300, level 4 at 600, level 5 at 1000, level 10 at 4500.
const levelStepXP = 100

// LevelInfo is a player's level and how far they are into it, for the progress bar.
type LevelInfo struct {
	Level int `json:"level"`
	// XPIntoLevel is the XP earned since reaching the current level.
	XPIntoLevel int `json:"xpIntoLevel"`
	// XPForNext is how much XP the current level takes in total to reach the next one.
	XPForNext int `json:"xpForNext"`
}

// XPForLevel returns the total XP needed to reach level (level 1 needs 0).
func XPForLevel(level int) int {
	if level <= 1 {
		return 0
	}
	return levelStepXP * (level - 1) * level / 2
}

// LevelFor returns the level for a total amount of XP.
func LevelFor(xp int) LevelInfo {
	if xp < 0 {
		xp = 0
	}
	level := 1
	for XPForLevel(level+1) <= xp {
		level++
	}
	start := XPForLevel(level)
	return LevelInfo{
		Level:       level,
		XPIntoLevel: xp - start,
		XPForNext:   XPForLevel(level+1) - start,
	}
}
