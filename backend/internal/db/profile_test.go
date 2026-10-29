package db_test

import (
	"context"
	"errors"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"
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

func createUser(t *testing.T, pool *pgxpool.Pool) (string, string) {
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
	return id, name
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
	userID, _ := createUser(t, pool)
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
	userID, _ := createUser(t, pool)
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
	userID, _ := createUser(t, pool)
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

func createFinishedRound(t *testing.T, pool *pgxpool.Pool, userID string, won bool) {
	t.Helper()
	ctx := context.Background()
	var sessionID string
	if err := pool.QueryRow(ctx,
		`INSERT INTO game_sessions (game_mode, status, ended_at) VALUES ('apartment_setup', 'completed', NOW()) RETURNING id`,
	).Scan(&sessionID); err != nil {
		t.Fatalf("create session: %v", err)
	}
	t.Cleanup(func() { pool.Exec(context.Background(), `DELETE FROM game_sessions WHERE id = $1`, sessionID) })
	if _, err := pool.Exec(ctx,
		`INSERT INTO session_participants (session_id, user_id, role, is_winner) VALUES ($1, $2, 'on_site', $3)`,
		sessionID, userID, won,
	); err != nil {
		t.Fatalf("add participant: %v", err)
	}
}

func TestGetPublicProfile_CountsOnlyCompletedRounds(t *testing.T) {
	pool := testPool(t)
	userID, username := createUser(t, pool)
	createFinishedRound(t, pool, userID, true)
	createFinishedRound(t, pool, userID, true)
	createFinishedRound(t, pool, userID, false)
	createSession(t, pool, userID, "active", nil)

	p, err := db.GetPublicProfile(context.Background(), pool, username)
	if err != nil {
		t.Fatalf("GetPublicProfile: %v", err)
	}
	if p.Username != username || p.RoundsPlayed != 3 || p.Wins != 2 || p.Losses != 1 || p.CreatedAt.IsZero() {
		t.Fatalf("unexpected profile %+v", p)
	}
}

func TestGetPublicProfile_NewUserHasZeroStats(t *testing.T) {
	pool := testPool(t)
	_, username := createUser(t, pool)

	p, err := db.GetPublicProfile(context.Background(), pool, username)
	if err != nil {
		t.Fatalf("GetPublicProfile: %v", err)
	}
	if p.RoundsPlayed != 0 || p.Wins != 0 || p.Losses != 0 {
		t.Fatalf("expected zero stats, got %+v", p)
	}
}

func TestGetPublicProfile_UnknownUser(t *testing.T) {
	pool := testPool(t)

	_, err := db.GetPublicProfile(context.Background(), pool, "no_such_user_xyz")
	if !errors.Is(err, pgx.ErrNoRows) {
		t.Fatalf("expected pgx.ErrNoRows, got %v", err)
	}
}

func TestUpdateProfile_ChangesOnlyGivenFields(t *testing.T) {
	pool := testPool(t)
	ctx := context.Background()
	userID, oldName := createUser(t, pool)
	newName := oldName + "x"

	p, err := db.UpdateProfile(ctx, pool, userID, &newName, nil, nil)
	if err != nil {
		t.Fatalf("UpdateProfile: %v", err)
	}
	if p.Username != newName || p.Email != oldName+"@example.com" {
		t.Fatalf("expected new username and unchanged email, got %+v", p)
	}
}

func TestUpdateProfile_TakenUsernameFails(t *testing.T) {
	pool := testPool(t)
	_, taken := createUser(t, pool)
	userID, _ := createUser(t, pool)

	if _, err := db.UpdateProfile(context.Background(), pool, userID, &taken, nil, nil); err == nil {
		t.Fatalf("expected a unique violation when taking another user's name")
	}
}

func TestUpdatePasswordHash(t *testing.T) {
	pool := testPool(t)
	ctx := context.Background()
	userID, _ := createUser(t, pool)

	if err := db.UpdatePasswordHash(ctx, pool, userID, "newhash"); err != nil {
		t.Fatalf("UpdatePasswordHash: %v", err)
	}
	if hash, err := db.GetPasswordHash(ctx, pool, userID); err != nil || hash != "newhash" {
		t.Fatalf("expected stored hash, got %q (err %v)", hash, err)
	}
}

func TestLearningLanguage_DefaultAndUpdate(t *testing.T) {
	pool := testPool(t)
	ctx := context.Background()
	userID, username := createUser(t, pool)

	if p, err := db.GetProfile(ctx, pool, userID); err != nil || p.LearningLanguage != "de" {
		t.Fatalf("new users should learn German by default, got %q (err %v)", p.LearningLanguage, err)
	}
	lang := "pl"
	if p, err := db.UpdateProfile(ctx, pool, userID, nil, nil, &lang); err != nil || p.LearningLanguage != "pl" {
		t.Fatalf("UpdateProfile language: got %q (err %v)", p.LearningLanguage, err)
	}
	if p, err := db.GetPublicProfile(ctx, pool, username); err != nil || p.LearningLanguage != "pl" {
		t.Fatalf("public profile should show the learning language, got %q (err %v)", p.LearningLanguage, err)
	}
}
