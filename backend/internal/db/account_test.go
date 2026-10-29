package db_test

import (
	"context"
	"errors"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"transcendence/backend/internal/db"
)

// hostedRound creates a round hosted by hostID with both users as participants.
func hostedRound(t *testing.T, pool *pgxpool.Pool, hostID, guestID, status string) string {
	t.Helper()
	ctx := context.Background()
	var sessionID string
	if err := pool.QueryRow(ctx,
		`INSERT INTO game_sessions (game_mode, status, host_user_id) VALUES ('apartment_setup', $1, $2) RETURNING id`,
		status, hostID,
	).Scan(&sessionID); err != nil {
		t.Fatalf("create session: %v", err)
	}
	t.Cleanup(func() { pool.Exec(context.Background(), `DELETE FROM game_sessions WHERE id = $1`, sessionID) })
	for _, uid := range []string{hostID, guestID} {
		if _, err := pool.Exec(ctx,
			`INSERT INTO session_participants (session_id, user_id, role, is_winner) VALUES ($1, $2, 'on_site', true)`,
			sessionID, uid,
		); err != nil {
			t.Fatalf("add participant: %v", err)
		}
	}
	return sessionID
}

func TestDeleteUser_RemovesAccountAndFriendships(t *testing.T) {
	pool := testPool(t)
	ctx := context.Background()
	a, _ := createUser(t, pool)
	b, _ := createUser(t, pool)
	befriend(t, pool, a, b)

	if err := db.DeleteUser(ctx, pool, a); err != nil {
		t.Fatalf("DeleteUser: %v", err)
	}
	if _, err := db.GetProfile(ctx, pool, a); !errors.Is(err, pgx.ErrNoRows) {
		t.Fatalf("deleted user should be gone, got %v", err)
	}
	if n := friendCount(t, pool, b); n != 0 {
		t.Fatalf("friendship should be gone, friend still has %d friends", n)
	}
}

func TestDeleteUser_KeepsFinishedRoundsForTheOtherPlayer(t *testing.T) {
	pool := testPool(t)
	ctx := context.Background()
	host, _ := createUser(t, pool)
	guest, _ := createUser(t, pool)
	hostedRound(t, pool, host, guest, "completed")

	if err := db.DeleteUser(ctx, pool, host); err != nil {
		t.Fatalf("DeleteUser: %v", err)
	}
	matches, err := db.GetMatchHistory(ctx, pool, guest)
	if err != nil || len(matches) != 1 {
		t.Fatalf("guest should keep the finished round, got %d (err %v)", len(matches), err)
	}
}

func TestDeleteUser_RemovesOpenLobbiesTheUserHosted(t *testing.T) {
	pool := testPool(t)
	ctx := context.Background()
	host, _ := createUser(t, pool)
	guest, _ := createUser(t, pool)
	sessionID := hostedRound(t, pool, host, guest, "waiting")

	if err := db.DeleteUser(ctx, pool, host); err != nil {
		t.Fatalf("DeleteUser: %v", err)
	}
	var count int
	if err := pool.QueryRow(ctx, `SELECT COUNT(*) FROM game_sessions WHERE id = $1`, sessionID).Scan(&count); err != nil {
		t.Fatalf("count: %v", err)
	}
	if count != 0 {
		t.Fatalf("an open lobby of a deleted host should be removed")
	}
}
