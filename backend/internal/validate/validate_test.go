package validate_test

import (
	"strings"
	"testing"

	"transcendence/backend/internal/validate"
)

func TestUsername(t *testing.T) {
	cases := []struct {
		in, want string
	}{
		{"", "Please choose a username."},
		{"ab", "Username must be at least 3 characters."},
		{strings.Repeat("x", 21), "Username can be at most 20 characters."},
		{"bad name!", `Use only letters, numbers, "_" and "-".`},
		{"émile", `Use only letters, numbers, "_" and "-".`},
		{"abc", ""},
		{"mira_99-x", ""},
		{strings.Repeat("x", 20), ""},
	}
	for _, c := range cases {
		if got := validate.Username(c.in); got != c.want {
			t.Errorf("Username(%q) = %q, want %q", c.in, got, c.want)
		}
	}
}

func TestEmail(t *testing.T) {
	cases := []struct {
		in, want string
	}{
		{"", "Please enter your email."},
		{"not-an-email", "Please enter a valid email address."},
		{"a@b", "Please enter a valid email address."},
		{"has space@example.com", "Please enter a valid email address."},
		{strings.Repeat("a", 250) + "@x.de", "Please enter a valid email address."},
		{"mira@example.com", ""},
	}
	for _, c := range cases {
		if got := validate.Email(c.in); got != c.want {
			t.Errorf("Email(%q) = %q, want %q", c.in, got, c.want)
		}
	}
}

func TestNewPassword(t *testing.T) {
	cases := []struct {
		in, want string
	}{
		{"", "Please choose a password."},
		{"short1", "Password must be at least 8 characters."},
		{"longpassword", "Use at least one letter and one number."},
		{"12345678", "Use at least one letter and one number."},
		{"a1" + strings.Repeat("x", 71), "Password is too long."},
		{"GoodPass123", ""},
		{"a1" + strings.Repeat("x", 70), ""},
	}
	for _, c := range cases {
		if got := validate.NewPassword(c.in); got != c.want {
			t.Errorf("NewPassword(%q) = %q, want %q", c.in, got, c.want)
		}
	}
}

func TestFirstError(t *testing.T) {
	if got := validate.FirstError("", "second", "third"); got != "second" {
		t.Errorf("FirstError = %q, want %q", got, "second")
	}
	if got := validate.FirstError("", ""); got != "" {
		t.Errorf("FirstError of valid values = %q, want empty", got)
	}
}
