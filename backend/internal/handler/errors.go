package handler

import (
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
)

// uniqueViolationMessage turns a unique-constraint error on users into a message for the client.
// It returns "" when err is not such an error.
func uniqueViolationMessage(err error) string {
	var pgErr *pgconn.PgError
	if !errors.As(err, &pgErr) || pgErr.Code != "23505" {
		return ""
	}
	switch pgErr.ConstraintName {
	case "users_username_key":
		return "username already taken"
	case "users_email_key":
		return "email already registered"
	default:
		return "user already exists"
	}
}
