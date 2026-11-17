package handler_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"

	"transcendence/backend/internal/handler"
)

// Bad query parameters are rejected before the database is used, so DB is nil here.
func TestLeaderboard_RejectsInvalidQuery(t *testing.T) {
	h := &handler.ProgressHandler{}
	for _, query := range []string{"?scope=everyone", "?limit=0", "?limit=101", "?limit=abc"} {
		t.Run(query, func(t *testing.T) {
			e := echo.New()
			req := httptest.NewRequest(http.MethodGet, "/api/leaderboard"+query, nil)
			c := e.NewContext(req, httptest.NewRecorder())
			c.Set("userID", "00000000-0000-0000-0000-000000000000")
			err := h.HandleLeaderboard(c)
			httpErr, ok := err.(*echo.HTTPError)
			if !ok || httpErr.Code != http.StatusBadRequest {
				t.Errorf("got %v, want 400", err)
			}
		})
	}
}
