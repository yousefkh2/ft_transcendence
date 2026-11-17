package db_test

import (
	"context"
	"testing"

	"transcendence/backend/internal/db"
)

func TestChat_FriendsOnlyAndConversationOrder(t *testing.T) {
	pool := testPool(t)
	ctx := context.Background()
	a, _ := createUser(t, pool)
	b, _ := createUser(t, pool)
	stranger, _ := createUser(t, pool)
	befriend(t, pool, a, b)

	if ok, err := db.AreFriends(ctx, pool, b, a); err != nil || !ok {
		t.Fatalf("a and b should be friends: ok=%v err=%v", ok, err)
	}
	if ok, err := db.AreFriends(ctx, pool, a, stranger); err != nil || ok {
		t.Fatalf("a and stranger should not be friends: ok=%v err=%v", ok, err)
	}

	for _, m := range []struct{ from, to, body string }{{a, b, "eins"}, {b, a, "zwei"}, {a, b, "drei"}} {
		if _, err := db.CreateChatMessage(ctx, pool, m.from, m.to, m.body); err != nil {
			t.Fatalf("create message: %v", err)
		}
	}

	all, err := db.ListChatMessages(ctx, pool, b, a, nil, 50)
	if err != nil || len(all) != 3 || all[0].Body != "eins" || all[2].Body != "drei" {
		t.Fatalf("expected eins, zwei, drei in order, got %+v (err %v)", all, err)
	}

	newest, err := db.ListChatMessages(ctx, pool, a, b, nil, 2)
	if err != nil || len(newest) != 2 || newest[0].Body != "zwei" {
		t.Fatalf("limit should keep the newest two, got %+v (err %v)", newest, err)
	}
	older, err := db.ListChatMessages(ctx, pool, a, b, &newest[0].CreatedAt, 50)
	if err != nil || len(older) != 1 || older[0].Body != "eins" {
		t.Fatalf("before should return the older message, got %+v (err %v)", older, err)
	}

	unread, err := db.UnreadChatCounts(ctx, pool, b)
	if err != nil || unread[a] != 2 {
		t.Fatalf("b should have 2 unread from a, got %v (err %v)", unread, err)
	}
	if err := db.MarkChatRead(ctx, pool, b, a); err != nil {
		t.Fatalf("mark read: %v", err)
	}
	unread, _ = db.UnreadChatCounts(ctx, pool, b)
	if unread[a] != 0 {
		t.Fatalf("expected no unread after marking read, got %v", unread)
	}

	// After removing the friendship the old unread messages no longer count.
	if _, err := db.CreateChatMessage(ctx, pool, a, b, "vier"); err != nil {
		t.Fatalf("create message: %v", err)
	}
	if _, err := db.RemoveFriend(ctx, pool, a, b); err != nil {
		t.Fatalf("remove friend: %v", err)
	}
	unread, _ = db.UnreadChatCounts(ctx, pool, b)
	if unread[a] != 0 {
		t.Fatalf("removed friend should not leave unread messages, got %v", unread)
	}
}
