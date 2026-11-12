package game

import "time"

// XP for a finished round. Both players of a round get the same amount.
const (
	// WinBaseXP is what every solved round earns, even one solved in the last second.
	WinBaseXP = 50
	// WinSpeedBonusXP is added on top for a fast solve: all of it for an instant solve, none when
	// the round took the whole RoundDuration. It shrinks evenly in between.
	WinSpeedBonusXP = 50
	// LossXP is a small reward for playing a round to the end (time up or out of checks).
	LossXP = 10
)

// MatchXP returns the XP each player earns for a round that lasted duration.
func MatchXP(won bool, duration time.Duration) int {
	if !won {
		return LossXP
	}
	remaining := RoundDuration - duration
	if remaining <= 0 {
		return WinBaseXP
	}
	if remaining > RoundDuration {
		remaining = RoundDuration
	}
	return WinBaseXP + int(int64(WinSpeedBonusXP)*int64(remaining)/int64(RoundDuration))
}
