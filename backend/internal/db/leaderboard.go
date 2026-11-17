package db

import (
	"context"
	"time"

	"transcendence/backend/internal/game"

	"github.com/jackc/pgx/v5/pgxpool"
)

// LeaderboardEntry is one ranked player.
type LeaderboardEntry struct {
	Rank         int            `json:"rank"`
	UserID       string         `json:"userId"`
	Username     string         `json:"username"`
	AvatarURL    *string        `json:"avatarUrl,omitempty"`
	XP           int            `json:"xp"`
	Level        game.LevelInfo `json:"level"`
	Wins         int            `json:"wins"`
	Rounds       int            `json:"rounds"`
	Achievements int            `json:"achievements"`
	IsMe         bool           `json:"isMe"`
}

type Leaderboard struct {
	Entries []LeaderboardEntry `json:"entries"`
	// Me is the user's own row, also when they are not in Entries.
	Me *LeaderboardEntry `json:"me,omitempty"`
}

// Players are ranked by XP, then wins; players with the same XP and wins share a rank.
// With friendsOnly, only the user and their friends are ranked.
const leaderboardQuery = `
WITH players AS (
	SELECT u.id, u.username, u.xp,
		(SELECT updated_at FROM user_avatars WHERE user_id = u.id) AS avatar_at,
		(SELECT COUNT(*) FROM session_participants sp JOIN game_sessions gs ON gs.id = sp.session_id
			WHERE sp.user_id = u.id AND gs.status = 'completed' AND sp.is_winner) AS wins,
		(SELECT COUNT(*) FROM session_participants sp JOIN game_sessions gs ON gs.id = sp.session_id
			WHERE sp.user_id = u.id AND gs.status = 'completed') AS rounds,
		(SELECT COUNT(*) FROM user_achievements ua WHERE ua.user_id = u.id) AS achievements
	FROM users u
	WHERE NOT $2::boolean OR u.id = $1 OR EXISTS (
		SELECT 1 FROM friendships f WHERE f.status = 'accepted'
		AND ((f.requester_id = $1 AND f.addressee_id = u.id) OR (f.addressee_id = $1 AND f.requester_id = u.id))
	)
),
ranked AS (
	SELECT *, RANK() OVER (ORDER BY xp DESC, wins DESC) AS rank FROM players
)
SELECT rank, id, username, avatar_at, xp, wins, rounds, achievements FROM ranked
WHERE rank <= $3 OR id = $1
ORDER BY rank, username`

// GetLeaderboard returns the top `limit` ranks, plus the user's own row.
func GetLeaderboard(ctx context.Context, pool *pgxpool.Pool, userID string, friendsOnly bool, limit int) (Leaderboard, error) {
	rows, err := pool.Query(ctx, leaderboardQuery, userID, friendsOnly, limit)
	if err != nil {
		return Leaderboard{}, err
	}
	defer rows.Close()

	board := Leaderboard{Entries: make([]LeaderboardEntry, 0)}
	for rows.Next() {
		var e LeaderboardEntry
		var avatarAt *time.Time
		if err := rows.Scan(&e.Rank, &e.UserID, &e.Username, &avatarAt, &e.XP, &e.Wins, &e.Rounds, &e.Achievements); err != nil {
			return Leaderboard{}, err
		}
		e.AvatarURL = avatarURL(e.UserID, avatarAt)
		e.Level = game.LevelFor(e.XP)
		e.IsMe = e.UserID == userID
		if e.IsMe {
			me := e
			board.Me = &me
		}
		if e.Rank <= limit {
			board.Entries = append(board.Entries, e)
		}
	}
	return board, rows.Err()
}
