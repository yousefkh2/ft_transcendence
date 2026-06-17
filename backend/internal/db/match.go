package db

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type MatchResult struct {
	GameMode		string
	StartedAt		time.Time
	EndedAt			time.Time
	Won				bool // true if round_completed, false if round_expired
	Participants	map[string]string
}

func SaveMatch(ctx context.Context, pool *pgxpool.Pool, result MatchResult) error {
	var sessionID string
	err := pool.QueryRow(ctx,
		`INSERT INTO game_sessions (game_mode, status, started_at, ended_at)
		VALUES ($1, 'completed', $2, $3)
		RETURNING id`,
		result.GameMode, result.StartedAt, result.EndedAt,
	).Scan(&sessionID)
	if err != nil {
		return err
	}

	for role, userID := range result.Participants {
		_, err := pool.Exec(ctx,
			`INSERT INTO session_participants (session_id, user_id, role, is_winner)
			VALUES ($1, $2, $3, $4)`,
			sessionID, userID, role, result.Won,
		)
		if err != nil {
			return err
		}
	}

	return nil
}
