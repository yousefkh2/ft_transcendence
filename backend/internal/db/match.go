package db

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func FinishMatch(ctx context.Context, pool *pgxpool.Pool, code string, startedAt, endedAt time.Time, won bool) error {
	tx, err := pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	var sessionID string
	err = tx.QueryRow(ctx,
		`SELECT id FROM game_sessions WHERE code = $1 FOR UPDATE`,
		code,
	).Scan(&sessionID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrLobbyNotFound
		}
		return err
	}

	if _, err := tx.Exec(ctx,
		`UPDATE session_participants SET is_winner = $1 WHERE session_id = $2`,
		won,sessionID,
	); err != nil {
		return err
	}

	return tx.Commit(ctx)
}
