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
}

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
	err := pool.QueryRow(ctx,
		`SELECT username, email, xp, created_at FROM users WHERE id = $1`,
		userID,
	).Scan(&p.Username, &p.Email, &p.XP, &p.CreatedAt)
	return p, err
}

// PublicProfile is what other logged-in users may see: no email, only completed rounds count.
type PublicProfile struct {
	Username	string		`json:"username"`
	XP			int			`json:"xp"`
	CreatedAt	time.Time	`json:"createdAt"`
	RoundsPlayed	int		`json:"roundsPlayed"`
	Wins		int			`json:"wins"`
	Losses		int			`json:"losses"`
}

func GetPublicProfile(ctx context.Context, pool *pgxpool.Pool, username string) (PublicProfile, error) {
	var p PublicProfile
	err := pool.QueryRow(ctx,
		`SELECT u.username, u.xp, u.created_at,
			COUNT(gs.id) FILTER (WHERE gs.status = 'completed'),
			COUNT(gs.id) FILTER (WHERE gs.status = 'completed' AND sp.is_winner),
			COUNT(gs.id) FILTER (WHERE gs.status = 'completed' AND NOT sp.is_winner)
		FROM users u
		LEFT JOIN session_participants sp ON sp.user_id = u.id
		LEFT JOIN game_sessions gs ON gs.id = sp.session_id
		WHERE u.username = $1
		GROUP BY u.id`,
		username,
	).Scan(&p.Username, &p.XP, &p.CreatedAt, &p.RoundsPlayed, &p.Wins, &p.Losses)
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
