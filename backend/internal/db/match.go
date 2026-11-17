package db

import (
	"context"
	"errors"
	"log"
	"time"

	"transcendence/backend/internal/game"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// How a round ended. Only a solved round counts as a win.
const (
	EndSolved      = "solved"
	EndTimeUp      = "time_up"
	EndOutOfChecks = "out_of_checks"
)

// FinishMatch saves a finished round: result, XP for both players and newly reached achievements.
func FinishMatch(ctx context.Context, pool *pgxpool.Pool, code string, startedAt, endedAt time.Time, endReason string, failedChecks int) error {
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
		`UPDATE game_sessions SET status = 'completed', started_at = $1, ended_at = $2, end_reason = $3, failed_checks = $5
		WHERE id = $4 AND status <> 'completed'`,
		startedAt, endedAt, endReason, sessionID, failedChecks,
	)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return nil
	}

	xp := game.MatchXP(won, endedAt.Sub(startedAt))
	if _, err := tx.Exec(ctx,
		`UPDATE session_participants SET is_winner = $1, score_delta = $2 WHERE session_id = $3`,
		won, xp, sessionID,
	); err != nil {
		return err
	}

	if _, err := tx.Exec(ctx,
		`UPDATE users SET xp = xp + $1
		WHERE id IN (SELECT user_id FROM session_participants WHERE session_id = $2)`,
		xp, sessionID,
	); err != nil {
		return err
	}

	var players []string
	rows, err := tx.Query(ctx, `SELECT user_id FROM session_participants WHERE session_id = $1`, sessionID)
	if err != nil {
		return err
	}
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			rows.Close()
			return err
		}
		players = append(players, id)
	}
	rows.Close()
	if err := rows.Err(); err != nil {
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		return err
	}

	// The round is saved; a failure here only delays the achievements until the next sync.
	for _, id := range players {
		if _, err := SyncAchievements(ctx, pool, id); err != nil {
			log.Printf("achievement sync failed for user %s: %v", id, err)
		}
	}
	return nil
}
