package handler

import (
	"errors"
	"net/http"
	"time"

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

type friendRequestEntry struct {
	ID					string		`json:"id"`
	RequesterId			string		`json:"requesterId"`
	RequestedUsername	string		`json:"requestedUsername"`
	CreatedAt			time.Time	`json:"createdAt"`
}

type friendEntry struct {
	FriendID		string		`json:"friendId"`
	FriendUsername	string		`json:"friendUsername"`
	CreatedAt		time.Time	`json:"createdAt"`
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

	var addresseeID string
	err := h.DB.QueryRow(c.Request().Context(),
	`SELECT id FROM users WHERE username = $1`,
	body.Username,
	).Scan(&addresseeID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return echo.NewHTTPError(http.StatusNotFound, "user not found")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "internal error")
	}

	if addresseeID == requesterID {
		return echo.NewHTTPError(http.StatusBadRequest, "cannot send friend request to yourself")
	}

	var id string
	err = h.DB.QueryRow(c.Request().Context(),
	`INSERT INTO friendships (requester_id, addressee_id, status)
	VALUES ($1, $2, 'pending')
	RETURNING id`,
	requesterID, addresseeID,
	).Scan(&id)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return echo.NewHTTPError(http.StatusConflict, "friend request already exits")
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

	rows, err := h.DB.Query(c.Request().Context(),
	`SELECT f.id, f.requester_id, u.username, f.created_at
	FROM friendships f
	JOIN users u ON u.id = f.requester_id
	WHERE f.addressee_id = $1 AND f.status = 'pending'
	ORDER BY f.created_at DESC`,
	userID,
	)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "internal error")
	}
	defer rows.Close()

	requests := make([]friendRequestEntry, 0)
	for rows.Next() {
		var entry friendRequestEntry
		if err := rows.Scan(&entry.ID, &entry.RequesterId, &entry.RequestedUsername, &entry.CreatedAt); err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "internal error")
		}
		requests = append(requests, entry)
	}
	if err := rows.Err(); err != nil {
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

	var addresseeID, status string
	err := h.DB.QueryRow(c.Request().Context(),
	`SELECT addressee_id, status FROM friendships WHERE id = $1`,
	requestID,
	).Scan(&addresseeID, &status)
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

	if _, err := h.DB.Exec(c.Request().Context(),
	`UPDATE friendships SET status = $1 WHERE id = $2`,
	targetStatus, requestID,
	); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "internal error")
	}

	return c.NoContent(http.StatusNoContent)
}

func (h *FriendHandler) HandleListFriends(c echo.Context) error {
	userID, ok := c.Get("userID").(string)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "unauthorized")
	}

	rows, err := h.DB.Query(c.Request().Context(),
	`SELECT
		CASE WHEN f.requester_id = $1 THEN f.addressee_id ELSE f.requester_id END,
		CASE WHEN f.requester_id = $1 THEN ua.username ELSE ur.username END,
		f.created_at
	FROM friendships f
	JOIN users ur ON ur.id = f.requester_id
	JOIN users ua ON ua.id = f.addressee_id
	WHERE (f.requester_id = $1 OR f.addressee_id = $1) AND f.status = 'accepted'
	ORDER BY f.created_at DESC`,
	userID,
	)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "internal error")
	}
	defer rows.Close()

	friends := make([]friendEntry, 0)
	for rows.Next() {
		var entry friendEntry
		if err := rows.Scan(&entry.FriendID, &entry.FriendUsername, &entry.CreatedAt); err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "internal error")
		}
		friends = append(friends, entry)
	}
	if err := rows.Err(); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "internal error")
	}

	return c.JSON(http.StatusOK, friends)
}
