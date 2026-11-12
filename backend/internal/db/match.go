package db

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// How a round ended. Only a solved round counts as a win.
const (
	EndSolved      = "solved"
	EndTimeUp      = "time_up"
	EndOutOfChecks = "out_of_checks"
)

func FinishMatch(ctx context.Context, pool *pgxpool.Pool, code string, startedAt, endedAt time.Time, endReason string) error {
	won := endReason == EndSolved

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

	// A round is saved only once: a second result for the same lobby is ignored.
	tag, err := tx.Exec(ctx,
		`UPDATE game_sessions SET status = 'completed', started_at = $1, ended_at = $2, end_reason = $3
		WHERE id = $4 AND status <> 'completed'`,
		startedAt, endedAt, endReason, sessionID,
	)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return nil
	}

	if _, err := tx.Exec(ctx,
		`UPDATE session_participants SET is_winner = $1 WHERE session_id = $2`,
		won, sessionID,
	); err != nil {
		return err
	}

	return tx.Commit(ctx)
}
