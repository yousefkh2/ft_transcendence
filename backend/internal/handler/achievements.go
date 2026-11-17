package handler

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"transcendence/backend/internal/db"
	"transcendence/backend/internal/middleware"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
)

const (
	leaderboardDefault = 50
	leaderboardMax     = 100
)

// ProgressHandler serves achievements and the leaderboard.
type ProgressHandler struct {
	DB *pgxpool.Pool
}

// HandleMyAchievements returns all achievements of the user with progress. Checking them here too
// means players also get achievements for rounds played before achievements existed.
func (h *ProgressHandler) HandleMyAchievements(c echo.Context) error {
	list, err := db.SyncAchievements(c.Request().Context(), h.DB, middleware.UserID(c))
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "internal error")
	}
	return c.JSON(http.StatusOK, list)
}

// HandleAchievementsSeen marks the newly unlocked achievements as shown.
func (h *ProgressHandler) HandleAchievementsSeen(c echo.Context) error {
	if err := db.MarkAchievementsSeen(c.Request().Context(), h.DB, middleware.UserID(c)); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "internal error")
	}
	return c.NoContent(http.StatusNoContent)
}

// HandleUserAchievements returns another user's achievements for their profile page.
func (h *ProgressHandler) HandleUserAchievements(c echo.Context) error {
	username := strings.TrimSpace(c.Param("username"))
	if username == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "username is required")
	}
	userID, err := db.ResolveUserID(c.Request().Context(), h.DB, username)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return echo.NewHTTPError(http.StatusNotFound, "user not found")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "internal error")
	}
	list, err := db.SyncAchievements(c.Request().Context(), h.DB, userID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "internal error")
	}
	// "New" is only meaningful to the owner.
	for i := range list {
		list[i].New = false
	}
	return c.JSON(http.StatusOK, list)
}

// HandleLeaderboard ranks players by XP: ?scope=global (default) or ?scope=friends, ?limit=1..100.
func (h *ProgressHandler) HandleLeaderboard(c echo.Context) error {
	friendsOnly := false
	switch c.QueryParam("scope") {
	case "", "global":
	case "friends":
		friendsOnly = true
	default:
		return echo.NewHTTPError(http.StatusBadRequest, "invalid scope")
	}
	limit := leaderboardDefault
	if raw := c.QueryParam("limit"); raw != "" {
		n, err := strconv.Atoi(raw)
		if err != nil || n < 1 || n > leaderboardMax {
			return echo.NewHTTPError(http.StatusBadRequest, "invalid limit")
		}
		limit = n
	}

	board, err := db.GetLeaderboard(c.Request().Context(), h.DB, middleware.UserID(c), friendsOnly, limit)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "internal error")
	}
	return c.JSON(http.StatusOK, board)
}
