package db_test

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/jackc/pgx/v5/pgconn"

	"transcendence/backend/internal/db"
)

// expectUniqueViolation checks that err is a unique violation on the given index.
func expectUniqueViolation(t *testing.T, err error, index string) {
	t.Helper()
	var pgErr *pgconn.PgError
	if !errors.As(err, &pgErr) || pgErr.Code != "23505" || pgErr.ConstraintName != index {
		t.Fatalf("expected unique violation on %s, got %v", index, err)
	}
}

func TestUsers_UsernameUniqueIgnoringCase(t *testing.T) {
	pool := testPool(t)
	_, name := createUser(t, pool)

	_, err := pool.Exec(context.Background(),
		`INSERT INTO users (username, email, password_hash) VALUES ($1, $2, 'x')`,
		strings.ToUpper(name), "other_"+name+"@example.com",
	)
	expectUniqueViolation(t, err, "users_username_lower_key")
}

func TestUsers_EmailUniqueIgnoringCase(t *testing.T) {
	pool := testPool(t)
	_, name := createUser(t, pool)

	_, err := pool.Exec(context.Background(),
		`INSERT INTO users (username, email, password_hash) VALUES ($1, $2, 'x')`,
		"other_"+name, strings.ToUpper(name+"@example.com"),
	)
	expectUniqueViolation(t, err, "users_email_lower_key")
}

func TestUsers_LookupByUsernameIgnoresCase(t *testing.T) {
	pool := testPool(t)
	ctx := context.Background()
	userID, name := createUser(t, pool)
	upper := strings.ToUpper(name)

	if id, err := db.ResolveUserID(ctx, pool, upper); err != nil || id != userID {
		t.Fatalf("ResolveUserID(%q) = %q, %v; want %q", upper, id, err, userID)
	}
	p, err := db.GetPublicProfile(ctx, pool, upper)
	if err != nil || p.Username != name {
		t.Fatalf("GetPublicProfile(%q) = %+v, %v; want stored name %q", upper, p, err, name)
	}
}

// A user may change only the case of their own name.
func TestUpdateProfile_CanChangeCaseOfOwnName(t *testing.T) {
	pool := testPool(t)
	userID, name := createUser(t, pool)
	upper := strings.ToUpper(name)

	p, err := db.UpdateProfile(context.Background(), pool, userID, &upper, nil, nil)
	if err != nil || p.Username != upper {
		t.Fatalf("UpdateProfile to %q = %+v, %v", upper, p, err)
	}
}
