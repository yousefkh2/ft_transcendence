package handler

import (
	"errors"
	"net/http"
	
	"transcendence/backend/internal/db"

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
	requesterID, ok := c.Get("userID").(string)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "unauthorized")
	}

	var body sendRequestBody
	if err := c.Bind(&body); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid json")
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
	userID, ok := c.Get("userID").(string)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "unauthorized")
	}

	requests, err := db.ListPendingFriendRequests(c.Request().Context(), h.DB, userID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "internal error")
	}

	return c.JSON(http.StatusOK, requests)
}

func (h *FriendHandler) HandleAcceptRequest(c echo.Context) error {
	return h.respondToRequest(c, "accepted")
}

func (h *FriendHandler) HandleDeclineRequest(c echo.Context) error {
	return h.respondToRequest(c, "declined")
}

func (h *FriendHandler) respondToRequest(c echo.Context, targetStatus string) error {
	userID, ok := c.Get("userID").(string)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "unauthorized")
	}

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

	if err := db.SetFriendRequestStatus(c.Request().Context(), h.DB, requestID, targetStatus);err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "internal error")
	}


	return c.NoContent(http.StatusNoContent)
}

func (h *FriendHandler) HandleListFriends(c echo.Context) error {
	userID, ok := c.Get("userID").(string)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "unauthorized")
	}

	friends, err := db.ListFriends(c.Request().Context(), h.DB, userID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "internal error")
	}

	return c.JSON(http.StatusOK, friends)
}
