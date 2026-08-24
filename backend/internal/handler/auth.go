package handler

import (
	"errors"
	"net/http"

	"transcendence/backend/internal/auth"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
	DB *pgxpool.Pool
}

type registerRequest struct {
	Username	string `json:"username"`
	Email		string `json:"email"`
	Password	string `json:"password"`
}

func (h *AuthHandler) HandleRegister(c echo.Context) error {
	var req registerRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid json")
	}

	if req.Username == "" || req.Email == "" || req.Password == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "username, email and password are required")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "internal error")
	}

	var userID string
	err = h.DB.QueryRow(c.Request().Context(),
		`INSERT INTO users (username, email, password_hash)
		VALUES ($1, $2, $3)
		RETURNING id`,
		req.Username, req.Email, string(hash),
	).Scan(&userID)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			switch pgErr.ConstraintName {
			case "users_username_key":
				return echo.NewHTTPError(http.StatusConflict, "username already taken")
			case "users_email_key":
				return echo.NewHTTPError(http.StatusConflict, "email already registered")
			default:
				return echo.NewHTTPError(http.StatusConflict, "user already exists")
			}
		}
			return echo.NewHTTPError(http.StatusInternalServerError, "internal error")
	}
	
	return c.JSON(http.StatusCreated, map[string]string{"id": userID})
}

type loginRequest struct {
	Email		string `json:"email"`
	Password	string `json:"password"`
}

func (h *AuthHandler) HandleLogin(c echo.Context) error {
	var req loginRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid json")
	}

	

	var userID, passwordHash string
	err := h.DB.QueryRow(c.Request().Context(),
		`SELECT id, password_hash FROM users WHERE email = $1`,
		req.Email,
	).Scan(&userID, &passwordHash)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return echo.NewHTTPError(http.StatusUnauthorized, "invalid credentials")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "internal server error")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(req.Password)); err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "invalid credentials")
	}

	token, err := auth.CreateJWT(userID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "internal server error")
	}

	return c.JSON(http.StatusOK, map[string]string{"token": token})
}
