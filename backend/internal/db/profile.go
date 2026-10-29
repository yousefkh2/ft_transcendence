package db

import (
	"context"
	"time"
	
	"github.com/jackc/pgx/v5/pgxpool"
)

type Profile struct {
	ID			string	`json:"id"`
	Username	string	`json:"username"`
	Email		string	`json:"email"`
	XP			int		`json:"xp"`
	CreatedAt	time.Time	`json:"createdAt"`
	LearningLanguage	string	`json:"learningLanguage"`
	AvatarURL	*string	`json:"avatarUrl,omitempty"`
}

// avatarUpdatedAt is a subquery for the picture version of the user in users.id (NULL without a picture).
const avatarUpdatedAt = `(SELECT updated_at FROM user_avatars WHERE user_id = users.id)`

type MatchHistoryEntry struct {
	SessionID	string		`json:"sessionId"`
	GameMode	string		`json:"gameMode"`
	Status		string		`json:"status"`
	StartedAt	time.Time	`json:"startedAt"`
	EndedAt		*time.Time	`json:"endedAt,omitempty"`
	Role		*string		`json:"role,omitempty"`
	IsWinner	*bool		`json:"isWinner,omitempty"`
}

func GetProfile(ctx context.Context, pool *pgxpool.Pool, userID string) (Profile, error) {
	p := Profile{ID: userID}
	var avatarAt *time.Time
	err := pool.QueryRow(ctx,
		`SELECT username, email, xp, created_at, learning_lang, `+avatarUpdatedAt+` FROM users WHERE id = $1`,
		userID,
	).Scan(&p.Username, &p.Email, &p.XP, &p.CreatedAt, &p.LearningLanguage, &avatarAt)
	p.AvatarURL = avatarURL(userID, avatarAt)
	return p, err
}

// PublicProfile is what other logged-in users may see: no email, only completed rounds count.
type PublicProfile struct {
	Username	string		`json:"username"`
	XP			int			`json:"xp"`
	CreatedAt	time.Time	`json:"createdAt"`
	LearningLanguage	string	`json:"learningLanguage"`
	AvatarURL	*string		`json:"avatarUrl,omitempty"`
	LastSeenAt	*time.Time	`json:"lastSeenAt,omitempty"`
	Online		bool		`json:"online"`
	RoundsPlayed	int		`json:"roundsPlayed"`
	Wins		int			`json:"wins"`
	Losses		int			`json:"losses"`
}

func GetPublicProfile(ctx context.Context, pool *pgxpool.Pool, username string) (PublicProfile, error) {
	var p PublicProfile
	var userID string
	var avatarAt *time.Time
	err := pool.QueryRow(ctx,
		`SELECT u.id, u.username, u.xp, u.created_at, u.learning_lang, u.last_seen_at,
			(SELECT updated_at FROM user_avatars WHERE user_id = u.id),
			COUNT(gs.id) FILTER (WHERE gs.status = 'completed'),
			COUNT(gs.id) FILTER (WHERE gs.status = 'completed' AND sp.is_winner),
			COUNT(gs.id) FILTER (WHERE gs.status = 'completed' AND NOT sp.is_winner)
		FROM users u
		LEFT JOIN session_participants sp ON sp.user_id = u.id
		LEFT JOIN game_sessions gs ON gs.id = sp.session_id
		WHERE u.username = $1
		GROUP BY u.id`,
		username,
	).Scan(&userID, &p.Username, &p.XP, &p.CreatedAt, &p.LearningLanguage, &p.LastSeenAt, &avatarAt,
		&p.RoundsPlayed, &p.Wins, &p.Losses)
	p.Online = IsOnline(p.LastSeenAt, time.Now())
	p.AvatarURL = avatarURL(userID, avatarAt)
	return p, err
}

func GetMatchHistory(ctx context.Context, pool *pgxpool.Pool, userID string) ([]MatchHistoryEntry, error) {
	rows, err := pool.Query(ctx,
		`SELECT gs.id, gs.game_mode, gs.status, gs.started_at, gs.ended_at, sp.role, sp.is_winner
		FROM session_participants sp
		JOIN game_sessions gs ON gs.id = sp.session_id
		WHERE sp.user_id = $1
		AND gs.status = 'completed'
		ORDER BY gs.started_at DESC`,
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	matches := make([]MatchHistoryEntry, 0)
	for rows.Next() {
		var entry MatchHistoryEntry
		if err := rows.Scan(
			&entry.SessionID, &entry.GameMode, &entry.Status,
			&entry.StartedAt, &entry.EndedAt, &entry.Role, &entry.IsWinner,
		); err != nil {
			return nil, err
		}
		matches = append(matches, entry)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return matches, nil
}

// DeleteUser deletes an account and everything that belongs to it (friendships, round participations,
// flashcards cascade). Finished rounds the user hosted are kept, so they stay in the other players'
// match history: only the host reference is cleared before the user row goes away.
func DeleteUser(ctx context.Context, pool *pgxpool.Pool, userID string) error {
	tx, err := pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	if _, err := tx.Exec(ctx,
		`UPDATE game_sessions SET host_user_id = NULL WHERE host_user_id = $1 AND status = 'completed'`,
		userID,
	); err != nil {
		return err
	}
	if _, err := tx.Exec(ctx, `DELETE FROM users WHERE id = $1`, userID); err != nil {
		return err
	}
	return tx.Commit(ctx)
}

func GetPasswordHash(ctx context.Context, pool *pgxpool.Pool, userID string) (string, error) {
	var hash string
	err := pool.QueryRow(ctx, `SELECT password_hash FROM users WHERE id = $1`, userID).Scan(&hash)
	return hash, err
}

// UpdateProfile changes username, email and/or learning language. A nil value keeps the current one.
func UpdateProfile(ctx context.Context, pool *pgxpool.Pool, userID string, username, email, learningLang *string) (Profile, error) {
	p := Profile{ID: userID}
	var avatarAt *time.Time
	err := pool.QueryRow(ctx,
		`UPDATE users SET username = COALESCE($2, username), email = COALESCE($3, email),
			learning_lang = COALESCE($4, learning_lang)
		WHERE id = $1
		RETURNING username, email, xp, created_at, learning_lang, `+avatarUpdatedAt,
		userID, username, email, learningLang,
	).Scan(&p.Username, &p.Email, &p.XP, &p.CreatedAt, &p.LearningLanguage, &avatarAt)
	p.AvatarURL = avatarURL(userID, avatarAt)
	return p, err
}

func UpdatePasswordHash(ctx context.Context, pool *pgxpool.Pool, userID, hash string) error {
	_, err := pool.Exec(ctx, `UPDATE users SET password_hash = $2 WHERE id = $1`, userID, hash)
	return err
}
