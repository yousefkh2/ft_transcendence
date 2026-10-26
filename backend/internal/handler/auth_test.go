package handler_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"

	"transcendence/backend/internal/handler"
)

// Invalid input is rejected before the database is used, so no DB is needed (DB is nil).
func postJSON(t *testing.T, fn func(echo.Context) error, body string) (int, string) {
	t.Helper()
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	err := fn(e.NewContext(req, rec))
	if err == nil {
		return rec.Code, rec.Body.String()
	}
	httpErr, ok := err.(*echo.HTTPError)
	if !ok {
		t.Fatalf("expected *echo.HTTPError, got %T: %v", err, err)
	}
	msg, _ := json.Marshal(httpErr.Message)
	return httpErr.Code, string(msg)
}

func TestRegister_RejectsInvalidInput(t *testing.T) {
	h := &handler.AuthHandler{}
	cases := []struct {
		name, body, wantMsg string
	}{
		{"missing fields", `{}`, "Please choose a username."},
		{"short username", `{"username":"ab","email":"a@example.com","password":"GoodPass123"}`, "Username must be at least 3 characters."},
		{"username too long for the DB", `{"username":"` + strings.Repeat("x", 60) + `","email":"a@example.com","password":"GoodPass123"}`, "Username can be at most 20 characters."},
		{"bad email", `{"username":"mira","email":"nope","password":"GoodPass123"}`, "Please enter a valid email address."},
		{"weak password", `{"username":"mira","email":"a@example.com","password":"password"}`, "Use at least one letter and one number."},
		{"password over bcrypt limit", `{"username":"mira","email":"a@example.com","password":"a1` + strings.Repeat("x", 80) + `"}`, "Password is too long."},
		{"only spaces count as empty", `{"username":"   ","email":"a@example.com","password":"GoodPass123"}`, "Please choose a username."},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			code, msg := postJSON(t, h.HandleRegister, c.body)
			if code != http.StatusBadRequest || !strings.Contains(msg, c.wantMsg) {
				t.Errorf("got %d %s, want 400 %q", code, msg, c.wantMsg)
			}
		})
	}
}

func TestLogin_RequiresEmailAndPassword(t *testing.T) {
	h := &handler.AuthHandler{}
	for _, fn := range []func(echo.Context) error{h.HandleLogin, h.HandleRealtimeLogin} {
		for _, body := range []string{`{}`, `{"email":"  ","password":"x"}`, `{"email":"a@example.com"}`} {
			code, msg := postJSON(t, fn, body)
			if code != http.StatusBadRequest || !strings.Contains(msg, "email and password are required") {
				t.Errorf("body %s: got %d %s, want 400", body, code, msg)
			}
		}
	}
}
