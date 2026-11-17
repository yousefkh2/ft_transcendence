package main

import (
	"log"

	"transcendence/backend/internal/db"
	"transcendence/backend/internal/handler"
	"transcendence/backend/internal/hub"
	"transcendence/backend/internal/middleware"

	"github.com/labstack/echo/v4"
	echomw "github.com/labstack/echo/v4/middleware"
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
	lobbyHandler := &handler.LobbyHandler{DB: pool}
	realtimeLobbyHandler := &handler.RealtimeLobbyHandler{DB: pool, Hub: hub}
	progressHandler := &handler.ProgressHandler{DB: pool}
	chatHandler := &handler.ChatHandler{DB: pool, Push: hub.PushChat}

	e := echo.New()
	e.Use(echomw.Recover())
	e.Use(middleware.WithCors)
	e.GET("/", handler.HandleRoot)
	e.GET("/health", handler.HandleHealth)
	e.GET("/health/db", handler.HandleDatabaseHealth)
	e.GET("/health/openai", handler.HandleOpenAIHealth)
	e.GET("/ws", hub.HandleWebSocket)
	e.GET("/ws/lobby", hub.HandleLobbyWebSocket)
	e.GET("/ws/voice", hub.HandleVoiceWebSocket)
	e.GET("/ws/chat", hub.HandleChatWebSocket)
	e.Any("/transcriptions", handler.HandleTranscription)
	e.Any("/realtime/transcription-session", handler.HandleRealtimeTranscriptionSession)

	e.POST("/api/auth/register", authHandler.HandleRegister)
	e.POST("/api/auth/login", authHandler.HandleLogin)
	e.POST("/api/v2/auth/login", authHandler.HandleRealtimeLogin)

	users := e.Group("/api/users", middleware.EchoWithAuth)
	users.GET("/me", profileHandler.HandleMe, middleware.EchoWithAuth)
	users.GET("/me/matches", profileHandler.HandleMatchHistory, middleware.EchoWithAuth)
	users.POST("/me/heartbeat", profileHandler.HandleHeartbeat)
	users.PATCH("/me", profileHandler.HandleUpdateProfile)
	users.PATCH("/me/password", profileHandler.HandleChangePassword)
	users.DELETE("/me", profileHandler.HandleDeleteAccount)
	users.PUT("/me/avatar", profileHandler.HandleUploadAvatar)
	users.DELETE("/me/avatar", profileHandler.HandleDeleteAvatar)
	users.GET("/me/achievements", progressHandler.HandleMyAchievements)
	users.POST("/me/achievements/seen", progressHandler.HandleAchievementsSeen)
	users.GET("/:username", profileHandler.HandleUserProfile)
	users.GET("/:username/achievements", progressHandler.HandleUserAchievements)

	// Public on purpose: <img> tags cannot send the login token.
	e.GET("/api/avatars/:userId", profileHandler.HandleAvatar)

	friends := e.Group("/api/friends", middleware.EchoWithAuth)
	friends.POST("/requests", friendHandler.HandleSendRequest)
	friends.GET("/requests", friendHandler.HandleListRequests)
	friends.POST("/requests/:id/accept", friendHandler.HandleAcceptRequest)
	friends.POST("/requests/:id/decline", friendHandler.HandleDeclineRequest)
	friends.GET("", friendHandler.HandleListFriends)
	friends.DELETE("/:friendId", friendHandler.HandleRemoveFriend)

	e.GET("/api/leaderboard", progressHandler.HandleLeaderboard, middleware.EchoWithAuth)

	chat := e.Group("/api/chat", middleware.EchoWithAuth)
	chat.GET("/unread", chatHandler.HandleUnread)
	chat.GET("/:friendId/messages", chatHandler.HandleListMessages)
	chat.POST("/:friendId/messages", chatHandler.HandleSendMessage)
	chat.POST("/:friendId/read", chatHandler.HandleMarkRead)

	lobbies := e.Group("/api/lobbies", middleware.EchoWithAuth)
	lobbies.POST("", lobbyHandler.HandleCreateLobby)
	lobbies.GET("", lobbyHandler.HandleListLobbies)
	lobbies.POST("/:code/join", lobbyHandler.HandleJoinLobby)
	lobbies.POST("/:code/leave", lobbyHandler.HandleLeaveLobby)
	lobbies.PATCH("/:code/language", lobbyHandler.HandleUpdateLobbyLanguage)

	realtimeLobbies := e.Group("/api/v2/lobbies", middleware.EchoWithAuth)
	realtimeLobbies.POST("", realtimeLobbyHandler.HandleCreateLobby)
	realtimeLobbies.POST("/:code/join", realtimeLobbyHandler.HandleJoinLobby)

	port := handler.Getenv("PORT", "8080")
	addr := ":" + port

	log.Printf("API server listening on http://localhost%s", addr)
	log.Fatal(e.Start(addr))
}
