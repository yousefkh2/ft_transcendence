package db_test

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"transcendence/backend/internal/db"
)

// testPool connects to the database in TEST_DATABASE_URL and skips the test when it is not set,
// e.g. TEST_DATABASE_URL=postgres://transcendence:<password>@localhost:5433/transcendence
func testPool(t *testing.T) *pgxpool.Pool {
	t.Helper()
	url := os.Getenv("TEST_DATABASE_URL")
	if url == "" {
		t.Skip("TEST_DATABASE_URL not set, skipping database test")
	}
	pool, err := pgxpool.New(context.Background(), url)
	if err != nil {
		t.Fatalf("connect: %v", err)
	}
	t.Cleanup(pool.Close)
	return pool
}

func createUser(t *testing.T, pool *pgxpool.Pool) string {
	t.Helper()
	name := fmt.Sprintf("dbtest%d", time.Now().UnixNano())
	var id string
	err := pool.QueryRow(context.Background(),
		`INSERT INTO users (username, email, password_hash) VALUES ($1, $2, 'x') RETURNING id`,
		name, name+"@example.com",
	).Scan(&id)
	if err != nil {
		t.Fatalf("create user: %v", err)
	}
	t.Cleanup(func() { pool.Exec(context.Background(), `DELETE FROM users WHERE id = $1`, id) })
	return id
}

func createSession(t *testing.T, pool *pgxpool.Pool, userID, status string, role *string) {
	t.Helper()
	ctx := context.Background()
	var sessionID string
	if err := pool.QueryRow(ctx,
		`INSERT INTO game_sessions (game_mode, status) VALUES ('apartment_setup', $1) RETURNING id`,
		status,
	).Scan(&sessionID); err != nil {
		t.Fatalf("create session: %v", err)
	}
	t.Cleanup(func() { pool.Exec(context.Background(), `DELETE FROM game_sessions WHERE id = $1`, sessionID) })
	if _, err := pool.Exec(ctx,
		`INSERT INTO session_participants (session_id, user_id, role) VALUES ($1, $2, $3)`,
		sessionID, userID, role,
	); err != nil {
		t.Fatalf("add participant: %v", err)
	}
}

func TestGetMatchHistory_ParticipantWithoutRole(t *testing.T) {
	pool := testPool(t)
	userID := createUser(t, pool)
	createSession(t, pool, userID, "completed", nil)

	matches, err := db.GetMatchHistory(context.Background(), pool, userID)
	if err != nil {
		t.Fatalf("GetMatchHistory returned an error for a participant without role: %v", err)
	}
	if len(matches) != 1 {
		t.Fatalf("expected 1 match, got %d", len(matches))
	}
	if matches[0].Role != nil {
		t.Errorf("expected nil role, got %q", *matches[0].Role)
	}
}

func TestGetMatchHistory_ParticipantWithRole(t *testing.T) {
	pool := testPool(t)
	userID := createUser(t, pool)
	role := "on_site"
	createSession(t, pool, userID, "completed", &role)

	matches, err := db.GetMatchHistory(context.Background(), pool, userID)
	if err != nil {
		t.Fatalf("GetMatchHistory: %v", err)
	}
	if len(matches) != 1 || matches[0].Role == nil || *matches[0].Role != role {
		t.Fatalf("expected one match with role %q, got %+v", role, matches)
	}
}

func TestGetMatchHistory_OnlyCompletedSessions(t *testing.T) {
	pool := testPool(t)
	userID := createUser(t, pool)
	createSession(t, pool, userID, "waiting", nil)
	createSession(t, pool, userID, "active", nil)
	createSession(t, pool, userID, "completed", nil)

	matches, err := db.GetMatchHistory(context.Background(), pool, userID)
	if err != nil {
		t.Fatalf("GetMatchHistory: %v", err)
	}
	if len(matches) != 1 || matches[0].Status != "completed" {
		t.Fatalf("expected only the completed session, got %+v", matches)
	}
}
