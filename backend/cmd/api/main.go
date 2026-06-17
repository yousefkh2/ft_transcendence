package main

import (
	"log"
	"net/http"

	"transcendence/backend/internal/db"
	"transcendence/backend/internal/handler"
	"transcendence/backend/internal/hub"
	"transcendence/backend/internal/middleware"
)


func main() {
	pool, err := db.Connect()
	if err != nil {
		log.Fatalf("database connection failed: %v", err)
	}
	defer pool.Close()
	
	hub := hub.NewHub(pool)

	authHandler := &handler.AuthHandler{DB: pool}

	mux := http.NewServeMux()
	mux.HandleFunc("/", handler.HandleRoot)
	mux.HandleFunc("/health", handler.HandleHealth)
	mux.HandleFunc("/health/db", handler.HandleDatabaseHealth)
	mux.HandleFunc("/health/openai", handler.HandleOpenAIHealth)
	mux.HandleFunc("/ws", hub.HandleWebSocket)
	mux.HandleFunc("/transcriptions", handler.HandleTranscription)
	mux.HandleFunc("/realtime/transcription-session", handler.HandleRealtimeTranscriptionSession)

	mux.HandleFunc("POST /api/auth/register", authHandler.HandleRegister)
	mux.HandleFunc("POST /api/auth/login", authHandler.HandleLogin)
	
	port := handler.Getenv("PORT", "8080")
	addr := ":" + port

	log.Printf("API server listening on http://localhost%s", addr)
	log.Fatal(http.ListenAndServe(addr, middleware.WithCORS(mux)))
}
