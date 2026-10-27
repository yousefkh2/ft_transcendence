package db

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// avatarURL is the public address of a user's picture, or nil without one. The version changes with
// every upload, so browsers can cache a picture forever and still get the new one after a change.
func avatarURL(userID string, updatedAt *time.Time) *string {
	if updatedAt == nil {
		return nil
	}
	url := fmt.Sprintf("/api/avatars/%s?v=%d", userID, updatedAt.UnixMilli())
	return &url
}

func SetAvatar(ctx context.Context, pool *pgxpool.Pool, userID string, image []byte, contentType string) error {
	_, err := pool.Exec(ctx,
		`INSERT INTO user_avatars (user_id, image, content_type) VALUES ($1, $2, $3)
		ON CONFLICT (user_id) DO UPDATE SET image = $2, content_type = $3, updated_at = NOW()`,
		userID, image, contentType,
	)
	return err
}

func GetAvatar(ctx context.Context, pool *pgxpool.Pool, userID string) (image []byte, contentType string, err error) {
	err = pool.QueryRow(ctx,
		`SELECT image, content_type FROM user_avatars WHERE user_id = $1`,
		userID,
	).Scan(&image, &contentType)
	return image, contentType, err
}

func DeleteAvatar(ctx context.Context, pool *pgxpool.Pool, userID string) error {
	_, err := pool.Exec(ctx, `DELETE FROM user_avatars WHERE user_id = $1`, userID)
	return err
}
