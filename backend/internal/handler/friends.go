package handler

import (
	"errors"
	"net/http"
	"strings"

	"transcendence/backend/internal/db"
	"transcendence/backend/internal/middleware"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
)

type FriendHandler struct {
	DB *pgxpool.Pool
}

type sendRequestBody struct {
	Username	string	`json:"username"`
}

func (h *FriendHandler) HandleSendRequest(c echo.Context) error {
	requesterID := middleware.UserID(c)

	var body sendRequestBody
	if err := c.Bind(&body); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid json")
	}

	body.Username = strings.TrimSpace(body.Username)
	if body.Username == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "username is required")
	}

	addresseeID, err := db.ResolveUserID(c.Request().Context(), h.DB, body.Username)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return echo.NewHTTPError(http.StatusNotFound, "user not found")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "internal error")
	}

	if addresseeID == requesterID {
		return echo.NewHTTPError(http.StatusBadRequest, "cannot send friend request to yourself")
	}

	id, err := db.CreateFriendRequest(c.Request().Context(), h.DB, requesterID, addresseeID)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return echo.NewHTTPError(http.StatusConflict, "friend request already sent")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "internal error")
	}

	return c.JSON(http.StatusCreated, map[string]string{"id": id})
}

func (h *FriendHandler) HandleListRequests(c echo.Context) error {
	userID := middleware.UserID(c)

	requests, err := db.ListPendingFriendRequests(c.Request().Context(), h.DB, userID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "internal error")
	}

	return c.JSON(http.StatusOK, requests)
}

func (h *FriendHandler) HandleAcceptRequest(c echo.Context) error {
	return h.respondToRequest(c, true)
}

func (h *FriendHandler) HandleDeclineRequest(c echo.Context) error {
	return h.respondToRequest(c, false)
}

func (h *FriendHandler) respondToRequest(c echo.Context, accept bool) error {
	userID := middleware.UserID(c)

	requestID := c.Param("id")

	addresseeID, status, err := db.GetFriendRequest(c.Request().Context(), h.DB, requestID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return echo.NewHTTPError(http.StatusNotFound, "friend request not found")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "internal error")
	}

	if addresseeID != userID {
		return echo.NewHTTPError(http.StatusForbidden, "not your friend request")
	}

	if status != "pending" {
		return echo.NewHTTPError(http.StatusConflict, "friend request already resolved")
	}

	if accept {
		err = db.SetFriendRequestStatus(c.Request().Context(), h.DB, requestID, "accepted")
	} else {
		err = db.DeleteFriendRequest(c.Request().Context(), h.DB, requestID)
	}
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "internal error")
	}


	return c.NoContent(http.StatusNoContent)
}

func (h *FriendHandler) HandleListFriends(c echo.Context) error {
	userID := middleware.UserID(c)

	friends, err := db.ListFriends(c.Request().Context(), h.DB, userID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "internal error")
	}

	return c.JSON(http.StatusOK, friends)
}

func (h *FriendHandler) HandleRemoveFriend(c echo.Context) error {
	userID := middleware.UserID(c)

	removed, err := db.RemoveFriend(c.Request().Context(), h.DB, userID, c.Param("friendId"))
	if err != nil {
		var pgErr *pgconn.PgError
		// 22P02: the id is not a valid UUID, so it cannot belong to a friend.
		if errors.As(err, &pgErr) && pgErr.Code == "22P02" {
			return echo.NewHTTPError(http.StatusNotFound, "friend not found")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "internal error")
	}
	if !removed {
		return echo.NewHTTPError(http.StatusNotFound, "friend not found")
	}

	return c.NoContent(http.StatusNoContent)
}
