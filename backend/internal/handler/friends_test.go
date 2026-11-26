package handler_test

import (
	"net/http"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"

	"transcendence/backend/internal/handler"
)

// A request id that is not a UUID is answered with 404 before the database is used, so DB is nil here.
func TestRespondToFriendRequest_InvalidIDIsNotFound(t *testing.T) {
	h := &handler.FriendHandler{}
	actions := map[string]func(echo.Context) error{
		"accept":  h.HandleAcceptRequest,
		"decline": h.HandleDeclineRequest,
	}
	for name, action := range actions {
		for _, id := range []string{"xyz", "not-a-uuid", "11111111-1111-1111-1111-11111111111"} {
			t.Run(name+" "+id, func(t *testing.T) {
				code, msg := asUser(t, func(ctx echo.Context) error {
					ctx.SetParamNames("id")
					ctx.SetParamValues(id)
					return action(ctx)
				}, "")
				if code != http.StatusNotFound || !strings.Contains(msg, "friend request not found") {
					t.Errorf("got %d %s, want 404 friend request not found", code, msg)
				}
			})
		}
	}
}
