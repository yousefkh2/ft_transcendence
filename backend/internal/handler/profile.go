package handler

import (
	"net/http"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
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

func (h *ProfileHandler) HandleMe(c echo.Context) error {
	userID, ok := c.Get("userID").(string)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "unauthorized")
	}

	var username, email string
	var xp int
	err := h.DB.QueryRow(c.Request().Context(),
		`SELECT username, email, xp FROM users WHERE id = $1`,
	userID,
	).Scan(&username, &email, &xp)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "internal error")
	}

	return c.JSON(http.StatusOK, map[string]any{
		"id":		userID,
		"username":	username,
		"email":	email,
		"xp":		xp,
	})
}

func (h *ProfileHandler) HandleMatchHistory(c echo.Context) error {
	userID, ok := c.Get("userID").(string)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "unauthorized")
	}

	rows, err := h.DB.Query(c.Request().Context(),
		`SELECT gs.id, gs.game_mode, gs.status, gs.started_at, gs.ended_at, sp.role, sp.is_winner
		FROM session_participants sp
		JOIN game_sessions gs ON gs.id = sp.session_id
		WHERE sp.user_id = $1
		ORDER BY gs.started_at DESC`,
		userID,
	)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "internal error")
	}
	defer rows.Close()

	matches := make([]matchHistoryEntry, 0)
	for rows.Next() {
		var entry matchHistoryEntry
		if err := rows.Scan(
			&entry.SessionID, &entry.GameMode, &entry.Status,
			&entry.StartedAt, &entry.EndedAt, &entry.Role, &entry.IsWinner,
		); err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "internal error")
		}
		matches = append(matches, entry)
	}
	if err := rows.Err(); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "internal error")
	}

	return c.JSON(http.StatusOK, matches)
}
