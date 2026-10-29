package db_test

import (
	"bytes"
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/jackc/pgx/v5"

	"transcendence/backend/internal/db"
)

func TestAvatar_UploadShowsUpEverywhereAndGoesAwayWithAccount(t *testing.T) {
	pool := testPool(t)
	ctx := context.Background()
	me, _ := createUser(t, pool)
	friend, friendName := createUser(t, pool)
	befriend(t, pool, me, friend)

	if p, err := db.GetProfile(ctx, pool, friend); err != nil || p.AvatarURL != nil {
		t.Fatalf("new user should have no avatar, got %v (%v)", p.AvatarURL, err)
	}

	image := []byte("fake png bytes")
	if err := db.SetAvatar(ctx, pool, friend, image, "image/png"); err != nil {
		t.Fatalf("SetAvatar: %v", err)
	}
	p, err := db.GetProfile(ctx, pool, friend)
	if err != nil || p.AvatarURL == nil || !strings.HasPrefix(*p.AvatarURL, "/api/avatars/"+friend+"?v=") {
		t.Fatalf("profile avatarUrl = %v (%v)", p.AvatarURL, err)
	}
	if pub, err := db.GetPublicProfile(ctx, pool, friendName); err != nil || pub.AvatarURL == nil || *pub.AvatarURL != *p.AvatarURL {
		t.Errorf("public profile avatarUrl = %v (%v), want %s", pub.AvatarURL, err, *p.AvatarURL)
	}
	friends, err := db.ListFriends(ctx, pool, me)
	if err != nil || len(friends) != 1 || friends[0].AvatarURL == nil {
		t.Errorf("friend list should carry the avatarUrl, got %+v (%v)", friends, err)
	}

	got, contentType, err := db.GetAvatar(ctx, pool, friend)
	if err != nil || !bytes.Equal(got, image) || contentType != "image/png" {
		t.Errorf("GetAvatar = %q %q (%v)", got, contentType, err)
	}

	// A new upload replaces the old one and changes the version in the URL.
	if err := db.SetAvatar(ctx, pool, friend, []byte("newer"), "image/jpeg"); err != nil {
		t.Fatalf("SetAvatar again: %v", err)
	}
	if again, _ := db.GetProfile(ctx, pool, friend); again.AvatarURL == nil || *again.AvatarURL == *p.AvatarURL {
		t.Errorf("avatarUrl should change after a new upload, still %v", again.AvatarURL)
	}

	if err := db.DeleteUser(ctx, pool, friend); err != nil {
		t.Fatalf("DeleteUser: %v", err)
	}
	if _, _, err := db.GetAvatar(ctx, pool, friend); !errors.Is(err, pgx.ErrNoRows) {
		t.Errorf("avatar should be deleted with the account, got %v", err)
	}
}

func TestDeleteAvatar(t *testing.T) {
	pool := testPool(t)
	ctx := context.Background()
	me, _ := createUser(t, pool)

	if err := db.SetAvatar(ctx, pool, me, []byte("x"), "image/png"); err != nil {
		t.Fatalf("SetAvatar: %v", err)
	}
	if err := db.DeleteAvatar(ctx, pool, me); err != nil {
		t.Fatalf("DeleteAvatar: %v", err)
	}
	if p, _ := db.GetProfile(ctx, pool, me); p.AvatarURL != nil {
		t.Errorf("avatarUrl should be gone, got %s", *p.AvatarURL)
	}
}
