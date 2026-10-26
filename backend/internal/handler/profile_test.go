package handler_test

import (
	"net/http"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"

	"transcendence/backend/internal/handler"
)

// asUser runs fn with a JSON body as if the auth middleware had already identified the user.
func asUser(t *testing.T, fn func(echo.Context) error, body string) (int, string) {
	t.Helper()
	return postJSON(t, func(c echo.Context) error {
		c.Set("userID", "00000000-0000-0000-0000-000000000000")
		return fn(c)
	}, body)
}

// Invalid input is rejected before the database is used, so DB is nil here.
func TestUpdateProfile_RejectsInvalidInput(t *testing.T) {
	h := &handler.ProfileHandler{}
	cases := []struct {
		name, body, wantMsg string
	}{
		{"nothing to update", `{}`, "nothing to update"},
		{"short username", `{"username":"ab"}`, "Username must be at least 3 characters."},
		{"username with spaces only", `{"username":"   "}`, "Please choose a username."},
		{"bad email", `{"email":"nope"}`, "Please enter a valid email address."},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			code, msg := asUser(t, h.HandleUpdateProfile, c.body)
			if code != http.StatusBadRequest || !strings.Contains(msg, c.wantMsg) {
				t.Errorf("got %d %s, want 400 %q", code, msg, c.wantMsg)
			}
		})
	}
}

func TestChangePassword_RejectsWeakNewPassword(t *testing.T) {
	h := &handler.ProfileHandler{}
	code, msg := asUser(t, h.HandleChangePassword, `{"currentPassword":"x","newPassword":"password"}`)
	if code != http.StatusBadRequest || !strings.Contains(msg, "Use at least one letter and one number.") {
		t.Errorf("got %d %s, want 400", code, msg)
	}
}

func TestChangePassword_RequiresCurrentPassword(t *testing.T) {
	h := &handler.ProfileHandler{}
	code, msg := asUser(t, h.HandleChangePassword, `{"newPassword":"GoodPass123"}`)
	if code != http.StatusBadRequest || !strings.Contains(msg, "current password is required") {
		t.Errorf("got %d %s, want 400", code, msg)
	}
}
