package handler_test

import (
	"bytes"
	"encoding/json"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"

	"transcendence/backend/internal/handler"
)

// uploadAvatar sends file as the "avatar" field of a multipart form, as a logged-in user.
func uploadAvatar(t *testing.T, field string, file []byte) (int, string) {
	t.Helper()
	var body bytes.Buffer
	form := multipart.NewWriter(&body)
	part, _ := form.CreateFormFile(field, "picture.png")
	part.Write(file)
	form.Close()

	e := echo.New()
	req := httptest.NewRequest(http.MethodPut, "/", &body)
	req.Header.Set(echo.HeaderContentType, form.FormDataContentType())
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("userID", "00000000-0000-0000-0000-000000000000")

	err := (&handler.ProfileHandler{}).HandleUploadAvatar(c)
	httpErr, ok := err.(*echo.HTTPError)
	if !ok {
		t.Fatalf("expected *echo.HTTPError, got %T: %v", err, err)
	}
	msg, _ := json.Marshal(httpErr.Message)
	return httpErr.Code, string(msg)
}

// Bad uploads are rejected before the database is used, so DB is nil here.
func TestUploadAvatar_RejectsBadFiles(t *testing.T) {
	cases := []struct {
		name, field string
		file        []byte
		wantCode    int
		wantMsg     string
	}{
		{"missing field", "picture", []byte("x"), http.StatusBadRequest, "avatar file is required"},
		{"not an image", "avatar", []byte("<svg onload=alert(1)>"), http.StatusBadRequest, "Please upload a PNG or JPEG image."},
		{"too large", "avatar", make([]byte, 3<<20), http.StatusRequestEntityTooLarge, "2 MB or smaller"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			code, msg := uploadAvatar(t, c.field, c.file)
			if code != c.wantCode || !strings.Contains(msg, c.wantMsg) {
				t.Errorf("got %d %s, want %d %q", code, msg, c.wantCode, c.wantMsg)
			}
		})
	}
}
