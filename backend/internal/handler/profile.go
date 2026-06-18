package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"transcendence/backend/internal/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ProfileHandler struct {
	DB *pgxpool.Pool
}

type matchHistoryEntry struct {
	SessionID	string		`json:"sessionId"`
	GameMode	string		`json:"gameMode"`
	Status		string		`json:"status"`
	StartedAt	time.Time	`json:"startedAt"`
	EndedAt		*time.Time	`json:"endedAt,omitempty"`
	Role		string		`json:"role"`
	IsWinner	*bool		`json:"isWinner,omitempty"`
}

func (h *ProfileHandler) HandleMe(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserIDKey).(string)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	var username, email string
	var xp int
	err := h.DB.QueryRow(r.Context(),
		`SELECT username, email, xp FROM users WHERE id = $1`,
	userID,
	).Scan(&username, &email, &xp)
	if err != nil {
		http.Error(w, "internal error oben", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{
		"id":		userID,
		"username":	username,
		"email":	email,
		"xp":		xp,
	})
}

func (h *ProfileHandler) HandleMatchHistory(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserIDKey).(string)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	rows, err := h.DB.Query(r.Context(),
		`SELECT gs.id, gs.game_mode, gs.status, gs.started_at, gs.ended_at, sp.role, sp.is_winner
		FROM session_participants sp
		JOIN game_sessions gs ON gs.id = sp.session_id
		WHERE sp.user_id = $1
		ORDER BY gs.started_at DESC`,
		userID,
	)
	if err != nil {
		http.Error(w, "internal error mitte", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	matches := make([]matchHistoryEntry, 0)
	for rows.Next() {
		var entry matchHistoryEntry
		if err := rows.Scan(
			&entry.SessionID, &entry.GameMode, &entry.Status,
			&entry.StartedAt, &entry.EndedAt, &entry.Role, &entry.IsWinner,
		); err != nil {
			http.Error(w, "internal error unten", http.StatusInternalServerError)
			return
		}
		matches = append(matches, entry)
	}
	if err := rows.Err(); err != nil {
		http.Error(w, "internal error ganz unten", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(matches)
}
