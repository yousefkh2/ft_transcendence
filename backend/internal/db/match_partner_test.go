package db_test

import (
	"context"
	"testing"

	"transcendence/backend/internal/db"
)

// Each player's match history names the other player of the round.
func TestGetMatchHistory_ShowsPartner(t *testing.T) {
	pool := testPool(t)
	ctx := context.Background()
	aliceID, alice := createUser(t, pool)
	bobID, bob := createUser(t, pool)
	hostedRound(t, pool, aliceID, bobID, "completed")

	for _, c := range []struct{ userID, wantPartner string }{{aliceID, bob}, {bobID, alice}} {
		history, err := db.GetMatchHistory(ctx, pool, c.userID)
		if err != nil {
			t.Fatalf("GetMatchHistory: %v", err)
		}
		if len(history) != 1 || history[0].PartnerUsername == nil || *history[0].PartnerUsername != c.wantPartner {
			t.Fatalf("expected one round with partner %q, got %+v", c.wantPartner, history)
		}
	}
}

// When the partner deleted their account the round stays, just without a partner name.
func TestGetMatchHistory_DeletedPartnerIsEmpty(t *testing.T) {
	pool := testPool(t)
	ctx := context.Background()
	aliceID, _ := createUser(t, pool)
	bobID, _ := createUser(t, pool)
	hostedRound(t, pool, aliceID, bobID, "completed")

	if err := db.DeleteUser(ctx, pool, bobID); err != nil {
		t.Fatalf("DeleteUser: %v", err)
	}
	history, err := db.GetMatchHistory(ctx, pool, aliceID)
	if err != nil {
		t.Fatalf("GetMatchHistory: %v", err)
	}
	if len(history) != 1 || history[0].PartnerUsername != nil {
		t.Fatalf("expected one round without partner, got %+v", history)
	}
}
