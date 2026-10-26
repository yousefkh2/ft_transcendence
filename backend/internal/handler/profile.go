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
	"golang.org/x/crypto/bcrypt"
)

type ProfileHandler struct {
	DB *pgxpool.Pool
}

func (h *ProfileHandler) HandleMe(c echo.Context) error {
	userID := middleware.UserID(c)

	profile, err := db.GetProfile(c.Request().Context(), h.DB, userID)
	if err != nil {
		// A valid token for an account that was deleted: the frontend logs out on 401.
		if errors.Is(err, pgx.ErrNoRows) {
			return echo.NewHTTPError(http.StatusUnauthorized, "unauthorized")
		}
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

type deleteAccountRequest struct {
	Password string `json:"password"`
}

// HandleDeleteAccount asks for the password again, so an open session alone cannot delete the account.
// A wrong password answers 403, not 401, because 401 means "your session is invalid" to the frontend.
func (h *ProfileHandler) HandleDeleteAccount(c echo.Context) error {
	userID := middleware.UserID(c)

	var req deleteAccountRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid json")
	}
	if req.Password == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "password is required")
	}

	hash, err := db.GetPasswordHash(c.Request().Context(), h.DB, userID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "internal error")
	}
	if bcrypt.CompareHashAndPassword([]byte(hash), []byte(req.Password)) != nil {
		return echo.NewHTTPError(http.StatusForbidden, "incorrect password")
	}

	if err := db.DeleteUser(c.Request().Context(), h.DB, userID); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "internal error")
	}
	return c.NoContent(http.StatusNoContent)
}
