// Package validate holds the input rules for user data. They mirror
// frontend/src/validation/rules.ts, so both sides reject the same input with the same message.
// Each function returns an error message, or "" when the value is valid.
package validate

import (
	"fmt"
	"regexp"
	"unicode/utf8"
)

const (
	usernameMin = 3
	usernameMax = 20
	emailMax    = 254
	passwordMin = 8
	// bcrypt only accepts up to 72 bytes.
	passwordMaxBytes = 72
)

var (
	usernamePattern = regexp.MustCompile(`^[A-Za-z0-9_-]+$`)
	emailPattern    = regexp.MustCompile(`^[^\s@]+@[^\s@]+\.[^\s@]+$`)
	letterPattern   = regexp.MustCompile(`[A-Za-z]`)
	digitPattern    = regexp.MustCompile(`[0-9]`)
)

func Username(username string) string {
	length := utf8.RuneCountInString(username)
	switch {
	case username == "":
		return "Please choose a username."
	case length < usernameMin:
		return fmt.Sprintf("Username must be at least %d characters.", usernameMin)
	case length > usernameMax:
		return fmt.Sprintf("Username can be at most %d characters.", usernameMax)
	case !usernamePattern.MatchString(username):
		return `Use only letters, numbers, "_" and "-".`
	}
	return ""
}

func Email(email string) string {
	switch {
	case email == "":
		return "Please enter your email."
	case utf8.RuneCountInString(email) > emailMax || !emailPattern.MatchString(email):
		return "Please enter a valid email address."
	}
	return ""
}

func NewPassword(password string) string {
	switch {
	case password == "":
		return "Please choose a password."
	case utf8.RuneCountInString(password) < passwordMin:
		return fmt.Sprintf("Password must be at least %d characters.", passwordMin)
	case len(password) > passwordMaxBytes:
		return "Password is too long."
	case !letterPattern.MatchString(password) || !digitPattern.MatchString(password):
		return "Use at least one letter and one number."
	}
	return ""
}

// FirstError returns the first non-empty message, so a handler can report one problem at a time.
func FirstError(messages ...string) string {
	for _, m := range messages {
		if m != "" {
			return m
		}
	}
	return ""
}
