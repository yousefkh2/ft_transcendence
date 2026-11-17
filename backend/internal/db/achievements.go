package db

import (
	"context"
	"time"

	"transcendence/backend/internal/game"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Achievement is one achievement as the API returns it.
type Achievement struct {
	game.AchievementProgress
	Unlocked   bool       `json:"unlocked"`
	UnlockedAt *time.Time `json:"unlockedAt,omitempty"`
	// New is true until the user has been shown that it was unlocked.
	New bool `json:"new"`
}

// playerStatsQuery computes game.PlayerStats from the completed rounds of user $1.
const playerStatsQuery = `
WITH r AS (
	SELECT gs.id, gs.lang, gs.started_at, gs.ended_at, gs.failed_checks, COALESCE(sp.is_winner, false) AS won
	FROM session_participants sp
	JOIN game_sessions gs ON gs.id = sp.session_id
	WHERE sp.user_id = $1 AND gs.status = 'completed' AND gs.ended_at IS NOT NULL
),
runs AS (
	-- Every loss starts a new group, so the wins inside one group are a streak.
	SELECT won, SUM(CASE WHEN won THEN 0 ELSE 1 END) OVER (ORDER BY ended_at, id) AS grp FROM r
)
SELECT
	(SELECT COUNT(*) FROM r),
	(SELECT COUNT(*) FROM r WHERE won),
	(SELECT COUNT(*) FROM r WHERE NOT won),
	(SELECT xp FROM users WHERE id = $1),
	(SELECT MIN(EXTRACT(EPOCH FROM ended_at - started_at)) FROM r WHERE won),
	EXISTS (SELECT 1 FROM r WHERE won AND failed_checks = 0),
	(SELECT COALESCE(MAX(n), 0) FROM (SELECT COUNT(*) AS n FROM runs WHERE won GROUP BY grp) streaks),
	(SELECT COUNT(DISTINCT lang) FROM r),
	EXISTS (
		SELECT 1 FROM r
		JOIN session_participants other ON other.session_id = r.id AND other.user_id <> $1
		JOIN friendships f ON f.status = 'accepted'
			AND ((f.requester_id = $1 AND f.addressee_id = other.user_id)
				OR (f.addressee_id = $1 AND f.requester_id = other.user_id))
	)`

// GetPlayerStats returns the stats the achievements are based on.
func GetPlayerStats(ctx context.Context, pool *pgxpool.Pool, userID string) (game.PlayerStats, error) {
	var s game.PlayerStats
	var fastest *float64
	err := pool.QueryRow(ctx, playerStatsQuery, userID).Scan(
		&s.Rounds, &s.Wins, &s.Losses, &s.XP, &fastest,
		&s.FlawlessSolve, &s.BestWinStreak, &s.Languages, &s.PlayedWithFriend,
	)
	if fastest != nil {
		s.FastestSolve = time.Duration(*fastest * float64(time.Second))
	}
	return s, err
}

// SyncAchievements stores every achievement the user has reached but not unlocked yet, and returns
// the full list with progress. Unlocked achievements stay unlocked, even if a stat drops later
// (for example after removing a friend).
func SyncAchievements(ctx context.Context, pool *pgxpool.Pool, userID string) ([]Achievement, error) {
	stats, err := GetPlayerStats(ctx, pool, userID)
	if err != nil {
		return nil, err
	}
	progress := game.Achievements(stats)

	var reached []string
	for _, a := range progress {
		if a.Done() {
			reached = append(reached, a.ID)
		}
	}
	if len(reached) > 0 {
		if _, err := pool.Exec(ctx,
			`INSERT INTO user_achievements (user_id, achievement_id)
			SELECT $1, unnest($2::text[])
			ON CONFLICT DO NOTHING`,
			userID, reached,
		); err != nil {
			return nil, err
		}
	}

	rows, err := pool.Query(ctx,
		`SELECT achievement_id, unlocked_at, seen_at IS NULL FROM user_achievements WHERE user_id = $1`,
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	type unlock struct {
		at    time.Time
		isNew bool
	}
	unlocked := map[string]unlock{}
	for rows.Next() {
		var id string
		var u unlock
		if err := rows.Scan(&id, &u.at, &u.isNew); err != nil {
			return nil, err
		}
		unlocked[id] = u
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	list := make([]Achievement, len(progress))
	for i, p := range progress {
		a := Achievement{AchievementProgress: p}
		if u, ok := unlocked[p.ID]; ok {
			at := u.at
			a.Unlocked, a.UnlockedAt, a.New = true, &at, u.isNew
			a.Progress = a.Target
		}
		list[i] = a
	}
	return list, nil
}

// MarkAchievementsSeen records that the user was shown their newly unlocked achievements.
func MarkAchievementsSeen(ctx context.Context, pool *pgxpool.Pool, userID string) error {
	_, err := pool.Exec(ctx,
		`UPDATE user_achievements SET seen_at = NOW() WHERE user_id = $1 AND seen_at IS NULL`,
		userID,
	)
	return err
}
