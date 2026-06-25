package handler

import (
	"net/http"
	"transcendence/backend/internal/db"
	"transcendence/backend/internal/middleware"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
)

type ProfileHandler struct {
	DB *pgxpool.Pool
}

func (h *ProfileHandler) HandleMe(c echo.Context) error {
	userID := middleware.UserID(c)

	profile, err := db.GetProfile(c.Request().Context(), h.DB, userID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "internal error")
	}

	return c.JSON(http.StatusOK, profile)
}

func (h *ProfileHandler) HandleMatchHistory(c echo.Context) error {
	userID := middleware.UserID(c)

	matches, err := db.GetMatchHistory(c.Request().Context(), h.DB, userID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "internal error")
	}

	return c.JSON(http.StatusOK, matches)
}
