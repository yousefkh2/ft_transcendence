package db

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type FriendRequest struct {
	ID					string		`json:"id"`
	RequesterID			string		`json:"requesterId"`
	RequesterUsername	string		`json:"requesterUsername"`
	CreatedAt			time.Time	`json:"createdAt"`
}

type Friend struct {
	FriendID		string		`json:"friendId"`
	FriendUsername	string		`json:"friendUsername"`
	CreatedAt		time.Time	`json:"createdAt"`
}

func ResolveUserID(ctx context.Context, pool *pgxpool.Pool, username string) (string, error) {
	var userID string
	err := pool.QueryRow(ctx,
		`SELECT id FROM users WHERE username = $1`,
		username,
	).Scan(&userID)
	return userID, err
}

func CreateFriendRequest(ctx context.Context, pool *pgxpool.Pool, requesterID, addresseeId string) (string, error) {
	var id string
	err := pool.QueryRow(ctx,
		`INSERT INTO friendships (requester_id, addressee_id, status)
		VALUES ($1, $2, 'pending')
		RETURNING id`,
		requesterID, addresseeId,
	).Scan(&id)
	return id, err
}

func ListPendingFriendRequests(ctx context.Context, pool *pgxpool.Pool, userID string) ([]FriendRequest, error) {
	rows, err := pool.Query(ctx,
		`SELECT f.id, f.requester_id, u.username, f.created_at
		FROM friendships f
		JOIN users u ON u.id = f.requester_id
		WHERE f.addressee_id = $1 AND f.status = 'pending'
		ORDER BY f.created_at DESC`,
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	requests := make([]FriendRequest, 0)
	for rows.Next() {
		var entry FriendRequest
		if err := rows.Scan(&entry.ID, &entry.RequesterID, &entry.RequesterUsername, &entry.CreatedAt); err != nil {
			return nil, err
		}
		requests = append(requests, entry)
	}
	if err := rows.Err(); err != nil {
		return nil ,err
	}
	return requests, nil
}

func GetFriendRequest(ctx context.Context, pool *pgxpool.Pool, requestID string) (addresseeID, status string, err error) {
	err = pool.QueryRow(ctx,
		`SELECT addressee_id, status FROM friendships WHERE id = $1`,
		requestID,
	).Scan(&addresseeID, &status)
	return addresseeID, status, err
}

func SetFriendRequestStatus(ctx context.Context, pool *pgxpool.Pool, requestID, targetStatus string) error {
	_, err := pool.Exec(ctx,
		`UPDATE friendships SET status = $1 WHERE id = $2`,
		targetStatus, requestID,
	)
	return err
}

func ListFriends(ctx context.Context, pool *pgxpool.Pool, userID string) ([]Friend, error) {
	rows, err := pool.Query(ctx,
		`SELECT
			CASE WHEN f.requester_id = $1 THEN f.addressee_id ELSE f.requester_id END,
			CASE WHEN f.requester_id = $1 THEN ua.username ELSE ur.username END,
			f.created_at
		FROM friendships f
		JOIN users ur ON ur.id = f.requester_id
		JOIN users ua ON ua.id = f.addressee_id
		WHERE (f.requester_id = $1 OR f.addressee_id = $1) AND f.status = 'accepted'
		ORDER BY f.created_at DESC`,
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	friends := make([]Friend, 0)
	for rows.Next() {
		var entry Friend
		if err := rows.Scan(&entry.FriendID, &entry.FriendUsername, &entry.CreatedAt); err != nil {
			return nil, err
		}
		friends = append(friends, entry)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return friends, nil
}
