package db

import (
	"context"
	"time"
	
	"github.com/jackc/pgx/v5/pgxpool"
)

type Profile struct {
	ID			string	`json:"id"`
	Username	string	`json:"username"`
	Email		string	`json:"email"`
	XP			int		`json:"xp"`
}

type MatchHistoryEntry struct {
	SessionID	string		`json:"sessionId"`
	GameMode	string		`json:"gameMode"`
	Status		string		`json:"status"`
	StartedAt	time.Time	`json:"startedAt"`
	EndedAt		*time.Time	`json:"endedAt,omitempty"`
	Role		string		`json:"role"`
	IsWinner	*bool		`json:"isWinner,omitempty"`
}

func GetProfile(ctx context.Context, pool *pgxpool.Pool, userID string) (Profile, error) {
	p := Profile{ID: userID}
	err := pool.QueryRow(ctx,
		`SELECT username, email, xp FROM users WHERE id = $1`,
		userID,
	).Scan(&p.Username, &p.Email, &p.XP)
	return p, err
}

func GetMatchHistory(ctx context.Context, pool *pgxpool.Pool, userID string) ([]MatchHistoryEntry, error) {
	rows, err := pool.Query(ctx,
		`SELECT gs.id, gs.game_mode, gs.status, gs.started_at, gs.ended_at, sp.role, sp.is_winner
		FROM session_participants sp
		JOIN game_sessions gs ON gs.id = sp.session_id
		WHERE sp.user_id = $1
		ORDER BY gs.started_at DESC`,
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	matches := make([]MatchHistoryEntry, 0)
	for rows.Next() {
		var entry MatchHistoryEntry
		if err := rows.Scan(
			&entry.SessionID, &entry.GameMode, &entry.Status,
			&entry.StartedAt, &entry.EndedAt, &entry.Role, &entry.IsWinner,
		); err != nil {
			return nil, err
		}
		matches = append(matches, entry)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return matches, nil
}
