package db_test

import (
	"context"
	"testing"
	"time"

	"transcendence/backend/internal/db"
)

func TestIsOnline(t *testing.T) {
	now := time.Now()
	ago := func(d time.Duration) *time.Time { ts := now.Add(-d); return &ts }
	cases := []struct {
		name     string
		lastSeen *time.Time
		want     bool
	}{
		{"never seen", nil, false},
		{"just now", ago(0), true},
		{"one missed heartbeat", ago(60 * time.Second), true},
		{"window passed", ago(db.OnlineWindow), false},
		{"hours ago", ago(4 * time.Hour), false},
	}
	for _, c := range cases {
		if got := db.IsOnline(c.lastSeen, now); got != c.want {
			t.Errorf("%s: IsOnline = %v, want %v", c.name, got, c.want)
		}
	}
}

func TestListFriends_ShowsOnlineStatus(t *testing.T) {
	pool := testPool(t)
	ctx := context.Background()
	me, _ := createUser(t, pool)
	active, _ := createUser(t, pool)
	away, _ := createUser(t, pool)
	never, _ := createUser(t, pool)
	for _, f := range []string{active, away, never} {
		befriend(t, pool, me, f)
	}

	if err := db.TouchLastSeen(ctx, pool, active); err != nil {
		t.Fatalf("TouchLastSeen: %v", err)
	}
	if _, err := pool.Exec(ctx, `UPDATE users SET last_seen_at = NOW() - INTERVAL '10 minutes' WHERE id = $1`, away); err != nil {
		t.Fatalf("set last seen: %v", err)
	}

	friends, err := db.ListFriends(ctx, pool, me)
	if err != nil {
		t.Fatalf("ListFriends: %v", err)
	}
	byID := map[string]db.Friend{}
	for _, f := range friends {
		byID[f.FriendID] = f
	}
	if f := byID[active]; !f.Online || f.LastSeenAt == nil {
		t.Errorf("friend with a fresh heartbeat should be online, got %+v", f)
	}
	if f := byID[away]; f.Online || f.LastSeenAt == nil {
		t.Errorf("friend seen 10 minutes ago should be offline with a last-seen time, got %+v", f)
	}
	if f := byID[never]; f.Online || f.LastSeenAt != nil {
		t.Errorf("friend who was never seen should be offline without a last-seen time, got %+v", f)
	}
}
