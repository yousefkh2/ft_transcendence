package handler

import (
	"net/http"
	"time"
	"strings"
	"os"
	"net"

	"github.com/labstack/echo/v4"
)

func Getenv(key string, fallback string) string {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return fallback
	}
	return value
}

func HandleRoot(c echo.Context) error {
	return c.String(http.StatusOK, "ft_transcendence API\n")
}

func HandleHealth(c echo.Context) error {
	return c.String(http.StatusOK, "ok\n")
}

func HandleDatabaseHealth(c echo.Context) error {
	host := Getenv("DB_HOST", "localhost")
	port := Getenv("DB_PORT", "5432")
	address := net.JoinHostPort(host, port)

	conn, err := net.DialTimeout("tcp", address, 2*time.Second)
	if err != nil {
		return c.String(http.StatusServiceUnavailable, "database unreachable \n")
	}
	_ = conn.Close()

	return c.String(http.StatusOK, "database healthy\n")
}

func HandleOpenAIHealth(c echo.Context) error {
	if strings.TrimSpace(os.Getenv("OPENAI_API_KEY")) == "" {
		return c.String(http.StatusServiceUnavailable, "openai api key missing\n")
	}

	return c.String(http.StatusOK, "openai api key configured\n")
}
