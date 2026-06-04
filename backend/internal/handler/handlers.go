package handler

import (
	"net/http"
	"time"
	"strings"
	"os"
	"net"
)

func Getenv(key string, fallback string) string {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return fallback
	}
	return value
}

func HandleRoot(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	_, _ = w.Write([]byte("ft_transcendence API\n"))
}

func HandleHealth(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	_, _ = w.Write([]byte("ok\n"))
}

func HandleDatabaseHealth(w http.ResponseWriter, _ *http.Request) {
	host := Getenv("DB_HOST", "localhost")
	port := Getenv("DB_PORT", "5432")
	address := net.JoinHostPort(host, port)

	conn, err := net.DialTimeout("tcp", address, 2*time.Second)
	if err != nil {
		http.Error(w, "database unreachable\n", http.StatusServiceUnavailable)
		return
	}
	_ = conn.Close()

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	_, _ = w.Write([]byte("database reachable\n"))
}

func HandleOpenAIHealth(w http.ResponseWriter, _ *http.Request) {
	if strings.TrimSpace(os.Getenv("OPENAI_API_KEY")) == "" {
		http.Error(w, "openai api key missing\n", http.StatusServiceUnavailable)
		return
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	_, _ = w.Write([]byte("openai api key configured\n"))
}
