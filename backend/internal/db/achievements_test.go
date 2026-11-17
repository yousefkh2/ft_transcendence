package db_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"transcendence/backend/internal/db"
)

// playRound creates a lobby with both users and finishes it through FinishMatch, like the hub does.
func playRound(t *testing.T, pool *pgxpool.Pool, a, b, lang, endReason string, duration time.Duration, failedChecks int) {
	t.Helper()
	ctx := context.Background()
	code := fmt.Sprintf("T%07d", time.Now().UnixNano()%10000000)
	var sessionID string
	if err := pool.QueryRow(ctx,
		`INSERT INTO game_sessions (game_mode, status, code, lang) VALUES ('apartment_setup', 'active', $1, $2) RETURNING id`,
		code, lang,
	).Scan(&sessionID); err != nil {
		t.Fatalf("create session: %v", err)
	}
	t.Cleanup(func() { pool.Exec(context.Background(), `DELETE FROM game_sessions WHERE id = $1`, sessionID) })
	for _, user := range []string{a, b} {
		if _, err := pool.Exec(ctx, `INSERT INTO session_participants (session_id, user_id) VALUES ($1, $2)`, sessionID, user); err != nil {
			t.Fatalf("add participant: %v", err)
		}
	}
	ended := time.Now()
	if err := db.FinishMatch(ctx, pool, code, ended.Add(-duration), ended, endReason, failedChecks); err != nil {
		t.Fatalf("finish match: %v", err)
	}
	time.Sleep(5 * time.Millisecond) // distinct ended_at for the streak order
}

func unlockedIDs(list []db.Achievement) map[string]bool {
	out := map[string]bool{}
	for _, a := range list {
		if a.Unlocked {
			out[a.ID] = true
		}
	}
	return out
}

func TestAchievements_FromPlayedRounds(t *testing.T) {
	pool := testPool(t)
	ctx := context.Background()
	a, _ := createUser(t, pool)
	b, _ := createUser(t, pool)
	befriend(t, pool, a, b)

	playRound(t, pool, a, b, "de", db.EndTimeUp, 5*time.Minute, 2)
	playRound(t, pool, a, b, "de", db.EndSolved, 50*time.Second, 0)
	playRound(t, pool, a, b, "pl", db.EndSolved, 3*time.Minute, 1)
	playRound(t, pool, a, b, "de", db.EndSolved, 4*time.Minute, 3)

	stats, err := db.GetPlayerStats(ctx, pool, a)
	if err != nil {
		t.Fatalf("stats: %v", err)
	}
	if stats.Rounds != 4 || stats.Wins != 3 || stats.Losses != 1 || stats.BestWinStreak != 3 ||
		stats.Languages != 2 || !stats.FlawlessSolve || !stats.PlayedWithFriend ||
		stats.FastestSolve < 49*time.Second || stats.FastestSolve > 51*time.Second {
		t.Fatalf("unexpected stats: %+v", stats)
	}

	list, err := db.SyncAchievements(ctx, pool, a)
	if err != nil {
		t.Fatalf("sync: %v", err)
	}
	got := unlockedIDs(list)
	for _, id := range []string{"first_round", "first_win", "first_loss", "quick_solve", "lightning", "flawless", "streak_3", "polyglot", "team_player"} {
		if !got[id] {
			t.Errorf("expected %s to be unlocked", id)
		}
	}
	for _, id := range []string{"rounds_10", "wins_10", "level_10"} {
		if got[id] {
			t.Errorf("%s should still be locked", id)
		}
	}
	if !list[0].New {
		t.Errorf("a fresh unlock should be new: %+v", list[0])
	}

	if err := db.MarkAchievementsSeen(ctx, pool, a); err != nil {
		t.Fatalf("mark seen: %v", err)
	}
	list, _ = db.SyncAchievements(ctx, pool, a)
	for _, ach := range list {
		if ach.New {
			t.Errorf("%s still new after marking seen", ach.ID)
		}
	}

	// Unlocked achievements are kept even when the stat behind them goes away.
	if _, err := db.RemoveFriend(ctx, pool, a, b); err != nil {
		t.Fatalf("remove friend: %v", err)
	}
	list, _ = db.SyncAchievements(ctx, pool, a)
	if !unlockedIDs(list)["team_player"] {
		t.Error("team_player should stay unlocked after removing the friend")
	}
}

func TestLeaderboard_RanksByXPAndShowsMe(t *testing.T) {
	pool := testPool(t)
	ctx := context.Background()
	top, _ := createUser(t, pool)
	friend, _ := createUser(t, pool)
	me, _ := createUser(t, pool)
	befriend(t, pool, me, friend)
	pool.Exec(ctx, `UPDATE users SET xp = 99999999 WHERE id = $1`, top)
	pool.Exec(ctx, `UPDATE users SET xp = 450 WHERE id = $1`, friend)
	pool.Exec(ctx, `UPDATE users SET xp = 120 WHERE id = $1`, me)

	board, err := db.GetLeaderboard(ctx, pool, me, false, 1)
	if err != nil {
		t.Fatalf("global: %v", err)
	}
	if len(board.Entries) != 1 || board.Entries[0].UserID != top || board.Entries[0].Rank != 1 {
		t.Fatalf("expected the top player first, got %+v", board.Entries)
	}
	if board.Me == nil || board.Me.UserID != me || board.Me.Rank <= 1 || board.Me.Level.Level != 2 {
		t.Fatalf("expected my own row outside the top 1, got %+v", board.Me)
	}

	friends, err := db.GetLeaderboard(ctx, pool, me, true, 50)
	if err != nil {
		t.Fatalf("friends: %v", err)
	}
	if len(friends.Entries) != 2 || friends.Entries[0].UserID != friend || friends.Entries[1].UserID != me || friends.Entries[1].Rank != 2 {
		t.Fatalf("friends board should be friend then me, got %+v", friends.Entries)
	}
}
