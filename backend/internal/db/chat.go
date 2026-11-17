package db

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// ChatMessage is one direct message between two friends.
type ChatMessage struct {
	ID          string     `json:"id"`
	SenderID    string     `json:"senderId"`
	RecipientID string     `json:"recipientId"`
	Body        string     `json:"body"`
	CreatedAt   time.Time  `json:"createdAt"`
	ReadAt      *time.Time `json:"readAt,omitempty"`
}

// AreFriends reports whether the two users have an accepted friendship, in either direction.
func AreFriends(ctx context.Context, pool *pgxpool.Pool, userID, otherID string) (bool, error) {
	var friends bool
	err := pool.QueryRow(ctx,
		`SELECT EXISTS (
			SELECT 1 FROM friendships
			WHERE status = 'accepted'
			AND ((requester_id = $1 AND addressee_id = $2) OR (requester_id = $2 AND addressee_id = $1))
		)`,
		userID, otherID,
	).Scan(&friends)
	return friends, err
}

// CreateChatMessage stores a message and returns it with its id and time.
func CreateChatMessage(ctx context.Context, pool *pgxpool.Pool, senderID, recipientID, body string) (ChatMessage, error) {
	message := ChatMessage{SenderID: senderID, RecipientID: recipientID, Body: body}
	err := pool.QueryRow(ctx,
		`INSERT INTO direct_messages (sender_id, recipient_id, body)
		VALUES ($1, $2, $3)
		RETURNING id, created_at`,
		senderID, recipientID, body,
	).Scan(&message.ID, &message.CreatedAt)
	return message, err
}

// ListChatMessages returns up to limit messages between the two users, oldest first.
// With before set, only messages older than that time are returned (to load earlier messages).
func ListChatMessages(ctx context.Context, pool *pgxpool.Pool, userID, otherID string, before *time.Time, limit int) ([]ChatMessage, error) {
	// The newest `limit` messages are picked in the inner query, then put back in chronological order.
	rows, err := pool.Query(ctx,
		`SELECT id, sender_id, recipient_id, body, created_at, read_at FROM (
			SELECT id, sender_id, recipient_id, body, created_at, read_at
			FROM direct_messages
			WHERE LEAST(sender_id, recipient_id) = LEAST($1::uuid, $2::uuid)
			AND GREATEST(sender_id, recipient_id) = GREATEST($1::uuid, $2::uuid)
			AND ($3::timestamptz IS NULL OR created_at < $3)
			ORDER BY created_at DESC
			LIMIT $4
		) newest
		ORDER BY created_at ASC`,
		userID, otherID, before, limit,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	messages := make([]ChatMessage, 0)
	for rows.Next() {
		var m ChatMessage
		if err := rows.Scan(&m.ID, &m.SenderID, &m.RecipientID, &m.Body, &m.CreatedAt, &m.ReadAt); err != nil {
			return nil, err
		}
		messages = append(messages, m)
	}
	return messages, rows.Err()
}

// MarkChatRead marks every message the other user sent to userID as read.
func MarkChatRead(ctx context.Context, pool *pgxpool.Pool, userID, otherID string) error {
	_, err := pool.Exec(ctx,
		`UPDATE direct_messages SET read_at = NOW()
		WHERE recipient_id = $1 AND sender_id = $2 AND read_at IS NULL`,
		userID, otherID,
	)
	return err
}

// UnreadChatCounts returns, per sender, how many messages userID has not read yet.
// Only current friends are counted, so a removed friend does not leave a badge behind.
func UnreadChatCounts(ctx context.Context, pool *pgxpool.Pool, userID string) (map[string]int, error) {
	rows, err := pool.Query(ctx,
		`SELECT m.sender_id, COUNT(*)
		FROM direct_messages m
		JOIN friendships f ON f.status = 'accepted'
			AND ((f.requester_id = m.sender_id AND f.addressee_id = m.recipient_id)
				OR (f.requester_id = m.recipient_id AND f.addressee_id = m.sender_id))
		WHERE m.recipient_id = $1 AND m.read_at IS NULL
		GROUP BY m.sender_id`,
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	counts := make(map[string]int)
	for rows.Next() {
		var senderID string
		var count int
		if err := rows.Scan(&senderID, &count); err != nil {
			return nil, err
		}
		counts[senderID] = count
	}
	return counts, rows.Err()
}
