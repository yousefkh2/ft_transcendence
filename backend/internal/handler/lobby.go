package handler

import (
	"errors"
	"net/http"

	"transcendence/backend/internal/db"
	"transcendence/backend/internal/middleware"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
)

type LobbyHandler struct {
	DB *pgxpool.Pool
}

type createLobbyRequest struct {
	GameMode	string `json:"gameMode"`
}

const defaultGameMode = "apartment_setup"

func (h *LobbyHandler) HandleCreateLobby(c echo.Context) error {
	userID := middleware.UserID(c)

	var req createLobbyRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid json")
	}
	gameMode := req.GameMode
	if gameMode == "" {
		gameMode = defaultGameMode
	}

	lobby, err := db.CreateLobby(c.Request().Context(), h.DB, userID, gameMode)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "internal error")
	}

	return c.JSON(http.StatusCreated, lobby)
}

func (h *LobbyHandler) HandleListLobbies(c echo.Context) error {
	lobbies, err := db.ListOpenLobbies(c.Request().Context(), h.DB)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "internal error")
	}

	return c.JSON(http.StatusOK, lobbies)
}

func (h *LobbyHandler) HandleJoinLobby(c echo.Context) error {
	userID := middleware.UserID(c)
	code := c.Param("code")


	lobby, err := db.JoinLobby(c.Request().Context(), h.DB, code, userID)
	if err != nil {
		switch {
		case errors.Is(err, db.ErrLobbyNotFound):
			return echo.NewHTTPError(http.StatusNotFound, "lobby not found")
		case errors.Is(err, db.ErrLobbyNotJoinAble):
			return echo.NewHTTPError(http.StatusConflict, "lobby is not joinable")
		case errors.Is(err, db.ErrLobbyFull):
			return echo.NewHTTPError(http.StatusConflict, "lobby is full")
		case errors.Is(err, db.ErrAlreadyJoined)
			return echo.NewHTTPError(http.StatusConflict, "already joioned this lobby")
		default:
			return echo.NewHTTPError(http.StatusInternalServerError, "internal error")
		}
	}
	
	return c.JSON(http.StatusOK, lobby)
}

func (h *LobbyHandler) HandleLeaveLobby(c echo.Context) error {
	userID := middleware.UserID(c)
	code := c.Param("code")

	if err := db.LeaveLobby(c.Request().Context(), h.DB, code, userID); err != nil {
		if errors.Is(err, db.ErrLobbyNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, "lobby not found")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "internal error")
	}
	
	return c.NoContent(http.StatusNoContent)
}
