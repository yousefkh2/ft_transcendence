package main

import (
	"log"
	"net/http"

	"transcendence/backend/internal/handler"
	"transcendence/backend/internal/hub"
	"transcendence/backend/internal/middleware"
)


func main() {
	hub := hub.NewHub()

	mux := http.NewServeMux()
	mux.HandleFunc("/", handler.HandleRoot)
	mux.HandleFunc("/health", handler.HandleHealth)
	mux.HandleFunc("/health/db", handler.HandleDatabaseHealth)
	mux.HandleFunc("/health/openai", handler.HandleOpenAIHealth)
	mux.HandleFunc("/ws", hub.HandleWebSocket)
	mux.HandleFunc("/transcriptions", handler.HandleTranscription)
	mux.HandleFunc("/realtime/transcription-session", handler.HandleRealtimeTranscriptionSession)

	port := handler.Getenv("PORT", "8080")
	addr := ":" + port

	log.Printf("API server listening on http://localhost%s", addr)
	log.Fatal(http.ListenAndServe(addr, middleware.WithCORS(mux)))
}
