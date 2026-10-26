package handler

import (
	"errors"
	"net/http"
	"strings"
	"transcendence/backend/internal/db"
	"transcendence/backend/internal/middleware"

	"github.com/jackc/pgx/v5"
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

func (h *ProfileHandler) HandleUserProfile(c echo.Context) error {
	username := strings.TrimSpace(c.Param("username"))
	if username == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "username is required")
	}

	profile, err := db.GetPublicProfile(c.Request().Context(), h.DB, username)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return echo.NewHTTPError(http.StatusNotFound, "user not found")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "internal error")
	}

	return c.JSON(http.StatusOK, profile)
}
