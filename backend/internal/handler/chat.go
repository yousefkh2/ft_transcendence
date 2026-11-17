package handler

import (
	"encoding/json"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"transcendence/backend/internal/db"
	"transcendence/backend/internal/middleware"
	"transcendence/backend/internal/model"
	"transcendence/backend/internal/validate"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
)

const (
	chatPageDefault = 50
	chatPageMax     = 100
)

var uuidPattern = regexp.MustCompile(`^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`)

// ChatHandler serves the direct messages between friends.
type ChatHandler struct {
	DB *pgxpool.Pool
	// Push delivers a message live to the open chat connections of the given users (hub.PushChat).
	Push func(userIDs []string, message model.ServerMessage)
}

type sendChatBody struct {
	Body string `json:"body"`
}

// friendParam returns the :friendId path parameter after checking that it is a friend of the user.
func (h *ChatHandler) friendParam(c echo.Context, userID string) (string, error) {
	friendID := c.Param("friendId")
	if !uuidPattern.MatchString(friendID) || friendID == userID {
		return "", echo.NewHTTPError(http.StatusNotFound, "friend not found")
	}
	friends, err := db.AreFriends(c.Request().Context(), h.DB, userID, friendID)
	if err != nil {
		return "", echo.NewHTTPError(http.StatusInternalServerError, "internal error")
	}
	if !friends {
		return "", echo.NewHTTPError(http.StatusForbidden, "you can only chat with friends")
	}
	return friendID, nil
}

// HandleSendMessage stores a message to a friend and pushes it live to both users.
func (h *ChatHandler) HandleSendMessage(c echo.Context) error {
	userID := middleware.UserID(c)

	var body sendChatBody
	if err := c.Bind(&body); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid json")
	}
	text := strings.TrimSpace(body.Body)
	if msg := validate.ChatMessage(text); msg != "" {
		return echo.NewHTTPError(http.StatusBadRequest, msg)
	}

	friendID, err := h.friendParam(c, userID)
	if err != nil {
		return err
	}

	message, err := db.CreateChatMessage(c.Request().Context(), h.DB, userID, friendID, text)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "internal error")
	}

	if h.Push != nil {
		if payload, err := json.Marshal(message); err == nil {
			h.Push([]string{userID, friendID}, model.ServerMessage{Type: "chat.message", Payload: payload})
		}
	}

	return c.JSON(http.StatusCreated, message)
}

// HandleListMessages returns the newest messages with a friend, or older ones with ?before=<time>.
func (h *ChatHandler) HandleListMessages(c echo.Context) error {
	userID := middleware.UserID(c)

	limit := chatPageDefault
	if raw := c.QueryParam("limit"); raw != "" {
		n, err := strconv.Atoi(raw)
		if err != nil || n < 1 || n > chatPageMax {
			return echo.NewHTTPError(http.StatusBadRequest, "invalid limit")
		}
		limit = n
	}
	var before *time.Time
	if raw := c.QueryParam("before"); raw != "" {
		parsed, err := time.Parse(time.RFC3339Nano, raw)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "invalid before")
		}
		before = &parsed
	}

	friendID, err := h.friendParam(c, userID)
	if err != nil {
		return err
	}

	messages, err := db.ListChatMessages(c.Request().Context(), h.DB, userID, friendID, before, limit)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "internal error")
	}
	return c.JSON(http.StatusOK, messages)
}

// HandleMarkRead marks the friend's messages to the user as read.
func (h *ChatHandler) HandleMarkRead(c echo.Context) error {
	userID := middleware.UserID(c)

	friendID, err := h.friendParam(c, userID)
	if err != nil {
		return err
	}
	if err := db.MarkChatRead(c.Request().Context(), h.DB, userID, friendID); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "internal error")
	}
	return c.NoContent(http.StatusNoContent)
}

// HandleUnread returns the number of unread messages per friend: {"<friendId>": 3}.
func (h *ChatHandler) HandleUnread(c echo.Context) error {
	counts, err := db.UnreadChatCounts(c.Request().Context(), h.DB, middleware.UserID(c))
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "internal error")
	}
	return c.JSON(http.StatusOK, counts)
}
