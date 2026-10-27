package handler

import (
	"errors"
	"net/http"
	"strings"
	"transcendence/backend/internal/db"
	"transcendence/backend/internal/middleware"
	"transcendence/backend/internal/validate"

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

type updateProfileRequest struct {
	Username         *string `json:"username"`
	Email            *string `json:"email"`
	LearningLanguage *string `json:"learningLanguage"`
	CurrentPassword  string  `json:"currentPassword"`
}

// HandleUpdateProfile changes username and/or email. A new email needs the current password,
// because whoever controls the email can take over the account.
func (h *ProfileHandler) HandleUpdateProfile(c echo.Context) error {
	userID := middleware.UserID(c)
	ctx := c.Request().Context()

	var req updateProfileRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid json")
	}
	if req.Username == nil && req.Email == nil && req.LearningLanguage == nil {
		return echo.NewHTTPError(http.StatusBadRequest, "nothing to update")
	}

	if req.Username != nil {
		username := strings.TrimSpace(*req.Username)
		if msg := validate.Username(username); msg != "" {
			return echo.NewHTTPError(http.StatusBadRequest, msg)
		}
		req.Username = &username
	}
	if req.LearningLanguage != nil && !db.SupportedLobbyLangs[*req.LearningLanguage] {
		return echo.NewHTTPError(http.StatusBadRequest, unsupportedLanguageMsg)
	}
	if req.Email != nil {
		email := strings.TrimSpace(*req.Email)
		if msg := validate.Email(email); msg != "" {
			return echo.NewHTTPError(http.StatusBadRequest, msg)
		}
		req.Email = &email

		current, err := db.GetProfile(ctx, h.DB, userID)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "internal error")
		}
		if email != current.Email {
			if err := h.checkPassword(c, userID, req.CurrentPassword); err != nil {
				return err
			}
		}
	}

	profile, err := db.UpdateProfile(ctx, h.DB, userID, req.Username, req.Email, req.LearningLanguage)
	if err != nil {
		if msg := uniqueViolationMessage(err); msg != "" {
			return echo.NewHTTPError(http.StatusConflict, msg)
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "internal error")
	}
	return c.JSON(http.StatusOK, profile)
}

type changePasswordRequest struct {
	CurrentPassword string `json:"currentPassword"`
	NewPassword     string `json:"newPassword"`
}

func (h *ProfileHandler) HandleChangePassword(c echo.Context) error {
	userID := middleware.UserID(c)

	var req changePasswordRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid json")
	}
	if msg := validate.NewPassword(req.NewPassword); msg != "" {
		return echo.NewHTTPError(http.StatusBadRequest, msg)
	}
	if err := h.checkPassword(c, userID, req.CurrentPassword); err != nil {
		return err
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "internal error")
	}
	if err := db.UpdatePasswordHash(c.Request().Context(), h.DB, userID, string(hash)); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "internal error")
	}
	return c.NoContent(http.StatusNoContent)
}

// checkPassword answers 403 (not 401) on a wrong password, because 401 makes the frontend log out.
func (h *ProfileHandler) checkPassword(c echo.Context, userID, password string) error {
	if password == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "current password is required")
	}
	hash, err := db.GetPasswordHash(c.Request().Context(), h.DB, userID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "internal error")
	}
	if bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) != nil {
		return echo.NewHTTPError(http.StatusForbidden, "incorrect password")
	}
	return nil
}

// HandleHeartbeat is called regularly by the open app, so friends can see the user as online.
func (h *ProfileHandler) HandleHeartbeat(c echo.Context) error {
	if err := db.TouchLastSeen(c.Request().Context(), h.DB, middleware.UserID(c)); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "internal error")
	}
	return c.NoContent(http.StatusNoContent)
}
