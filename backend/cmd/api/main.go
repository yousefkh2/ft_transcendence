package main

import (
	"log"

	"transcendence/backend/internal/db"
	"transcendence/backend/internal/handler"
	"transcendence/backend/internal/hub"
	"transcendence/backend/internal/middleware"

	"github.com/labstack/echo/v4"
)


func main() {
	pool, err := db.Connect()
	if err != nil {
		log.Fatalf("database connection failed: %v", err)
	}
	defer pool.Close()
	
	hub := hub.NewHub(pool)

	authHandler := &handler.AuthHandler{DB: pool}
	profileHandler := &handler.ProfileHandler{DB: pool}
	friendHandler := &handler.FriendHandler{DB: pool}
	
	e := echo.New()
	e.Use(middleware.WithCors)
	e.Any("/", handler.HandleRoot)
	e.Any("/health", handler.HandleHealth)
	e.Any("/health/db", handler.HandleDatabaseHealth)
	e.Any("/health/openai", handler.HandleOpenAIHealth)
	e.Any("/ws", hub.HandleWebSocket)
	e.Any("/transcriptions", handler.HandleTranscription)
	e.Any("/realtime/transcription-session", handler.HandleRealtimeTranscriptionSession)

	e.POST("/api/auth/register", authHandler.HandleRegister)
	e.POST("/api/auth/login", authHandler.HandleLogin)

	e.GET("/api/users/me", profileHandler.HandleMe, middleware.EchoWithAuth)
	e.GET("/api/users/me/matches", profileHandler.HandleMatchHistory, middleware.EchoWithAuth)

	friends := e.Group("/api/friends", middleware.EchoWithAuth)
	friends.POST("/requests", friendHandler.HandleSendRequest)
	friends.GET("/requests", friendHandler.HandleListRequests)
	friends.POST("/requests/:id/accept", friendHandler.HandleAcceptRequest)
	friends.POST("/requests/:id/decline", friendHandler.HandleDeclineRequest)
	friends.GET("", friendHandler.HandleListFriends)

	port := handler.Getenv("PORT", "8080")
	addr := ":" + port

	log.Printf("API server listening on http://localhost%s", addr)
	log.Fatal(e.Start(addr))
}
