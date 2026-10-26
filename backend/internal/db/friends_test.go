package db_test

import (
	"context"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"

	"transcendence/backend/internal/db"
)

func befriend(t *testing.T, pool *pgxpool.Pool, requesterID, addresseeID string) {
	t.Helper()
	ctx := context.Background()
	id, err := db.CreateFriendRequest(ctx, pool, requesterID, addresseeID)
	if err != nil {
		t.Fatalf("create request: %v", err)
	}
	if err := db.SetFriendRequestStatus(ctx, pool, id, "accepted"); err != nil {
		t.Fatalf("accept request: %v", err)
	}
}

func friendCount(t *testing.T, pool *pgxpool.Pool, userID string) int {
	t.Helper()
	friends, err := db.ListFriends(context.Background(), pool, userID)
	if err != nil {
		t.Fatalf("list friends: %v", err)
	}
	return len(friends)
}

func TestRemoveFriend_WorksFromBothSides(t *testing.T) {
	pool := testPool(t)
	ctx := context.Background()
	a, _ := createUser(t, pool)
	b, _ := createUser(t, pool)
	c, _ := createUser(t, pool)
	befriend(t, pool, a, b)
	befriend(t, pool, a, c)

	// b was the addressee, but can still remove the friendship.
	removed, err := db.RemoveFriend(ctx, pool, b, a)
	if err != nil || !removed {
		t.Fatalf("addressee removing friend: removed=%v err=%v", removed, err)
	}
	removed, err = db.RemoveFriend(ctx, pool, a, c)
	if err != nil || !removed {
		t.Fatalf("requester removing friend: removed=%v err=%v", removed, err)
	}
	if n := friendCount(t, pool, a); n != 0 {
		t.Fatalf("expected no friends left, got %d", n)
	}
}

func TestRemoveFriend_NotFriends(t *testing.T) {
	pool := testPool(t)
	a, _ := createUser(t, pool)
	b, _ := createUser(t, pool)

	removed, err := db.RemoveFriend(context.Background(), pool, a, b)
	if err != nil || removed {
		t.Fatalf("expected removed=false without error, got removed=%v err=%v", removed, err)
	}
}

func TestRemoveFriend_DoesNotTouchPendingRequests(t *testing.T) {
	pool := testPool(t)
	ctx := context.Background()
	a, _ := createUser(t, pool)
	b, _ := createUser(t, pool)
	if _, err := db.CreateFriendRequest(ctx, pool, a, b); err != nil {
		t.Fatalf("create request: %v", err)
	}

	removed, err := db.RemoveFriend(ctx, pool, a, b)
	if err != nil || removed {
		t.Fatalf("pending request must not count as a friend: removed=%v err=%v", removed, err)
	}
	pending, err := db.ListPendingFriendRequests(ctx, pool, b)
	if err != nil || len(pending) != 1 {
		t.Fatalf("pending request should still exist, got %d (err %v)", len(pending), err)
	}
}

func TestRemoveFriend_CanBeAddedAgainAfterwards(t *testing.T) {
	pool := testPool(t)
	ctx := context.Background()
	a, _ := createUser(t, pool)
	b, _ := createUser(t, pool)
	befriend(t, pool, a, b)

	if _, err := db.RemoveFriend(ctx, pool, a, b); err != nil {
		t.Fatalf("remove: %v", err)
	}
	if _, err := db.CreateFriendRequest(ctx, pool, b, a); err != nil {
		t.Fatalf("sending a new request after removal should work, got %v", err)
	}
}

func TestDeleteFriendRequest_SenderCanAskAgain(t *testing.T) {
	pool := testPool(t)
	ctx := context.Background()
	a, _ := createUser(t, pool)
	b, _ := createUser(t, pool)
	id, err := db.CreateFriendRequest(ctx, pool, a, b)
	if err != nil {
		t.Fatalf("create request: %v", err)
	}

	if err := db.DeleteFriendRequest(ctx, pool, id); err != nil {
		t.Fatalf("decline: %v", err)
	}
	if pending, _ := db.ListPendingFriendRequests(ctx, pool, b); len(pending) != 0 {
		t.Fatalf("declined request should be gone, still %d pending", len(pending))
	}
	if _, err := db.CreateFriendRequest(ctx, pool, a, b); err != nil {
		t.Fatalf("sender should be able to ask again after a decline, got %v", err)
	}
}

func TestDeleteFriendRequest_KeepsAcceptedFriendship(t *testing.T) {
	pool := testPool(t)
	ctx := context.Background()
	a, _ := createUser(t, pool)
	b, _ := createUser(t, pool)
	id, err := db.CreateFriendRequest(ctx, pool, a, b)
	if err != nil {
		t.Fatalf("create request: %v", err)
	}
	if err := db.SetFriendRequestStatus(ctx, pool, id, "accepted"); err != nil {
		t.Fatalf("accept: %v", err)
	}

	if err := db.DeleteFriendRequest(ctx, pool, id); err != nil {
		t.Fatalf("delete: %v", err)
	}
	if n := friendCount(t, pool, a); n != 1 {
		t.Fatalf("an accepted friendship must not be deleted as a request, friends: %d", n)
	}
}
