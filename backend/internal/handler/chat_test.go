package handler_test

import (
	"net/http"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"

	"transcendence/backend/internal/handler"
)

// Invalid messages and friend ids are rejected before the database is used, so DB is nil here.
func TestSendMessage_RejectsInvalidInput(t *testing.T) {
	h := &handler.ChatHandler{}
	cases := []struct {
		name, friendID, body string
		wantCode             int
		wantMsg              string
	}{
		{"empty message", "11111111-1111-1111-1111-111111111111", `{"body":""}`, http.StatusBadRequest, "Please enter a message."},
		{"whitespace only", "11111111-1111-1111-1111-111111111111", `{"body":"   \n "}`, http.StatusBadRequest, "Please enter a message."},
		{"too long", "11111111-1111-1111-1111-111111111111", `{"body":"` + strings.Repeat("x", 1001) + `"}`, http.StatusBadRequest, "at most 1000 characters"},
		{"friend id not a uuid", "nope", `{"body":"hallo"}`, http.StatusNotFound, "friend not found"},
		{"message to yourself", "00000000-0000-0000-0000-000000000000", `{"body":"hallo"}`, http.StatusNotFound, "friend not found"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			code, msg := asUser(t, func(ctx echo.Context) error {
				ctx.SetParamNames("friendId")
				ctx.SetParamValues(c.friendID)
				return h.HandleSendMessage(ctx)
			}, c.body)
			if code != c.wantCode || !strings.Contains(msg, c.wantMsg) {
				t.Errorf("got %d %s, want %d %q", code, msg, c.wantCode, c.wantMsg)
			}
		})
	}
}

// 1000 characters is allowed, counted as characters and not bytes.
func TestSendMessage_LimitCountsCharacters(t *testing.T) {
	h := &handler.ChatHandler{}
	code, msg := asUser(t, func(ctx echo.Context) error {
		ctx.SetParamNames("friendId")
		ctx.SetParamValues("not-a-uuid")
		return h.HandleSendMessage(ctx)
	}, `{"body":"`+strings.Repeat("ü", 1000)+`"}`)
	// Passes the length check and only fails on the friend id.
	if code != http.StatusNotFound {
		t.Errorf("got %d %s, want 404 after passing the length check", code, msg)
	}
}
