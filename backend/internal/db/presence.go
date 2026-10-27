package db

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// OnlineWindow is how recent the last heartbeat must be to count as online. The frontend sends one
// every 30 seconds; the extra margin covers browsers that slow down timers in background tabs.
const OnlineWindow = 90 * time.Second

// IsOnline reports whether a user with this last heartbeat counts as online at the given time.
func IsOnline(lastSeen *time.Time, now time.Time) bool {
	return lastSeen != nil && now.Sub(*lastSeen) < OnlineWindow
}

// TouchLastSeen records a heartbeat for the user.
func TouchLastSeen(ctx context.Context, pool *pgxpool.Pool, userID string) error {
	_, err := pool.Exec(ctx, `UPDATE users SET last_seen_at = NOW() WHERE id = $1`, userID)
	return err
}
