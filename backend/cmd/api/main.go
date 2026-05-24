package main

import (
	"context"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/coder/websocket"
	"github.com/coder/websocket/wsjson"
)

type ClientMessage struct {
	Type     string `json:"type"`
	RoomCode string `json:"roomCode"`
}

type ServerMessage struct {
	Type     string `json:"type"`
	RoomCode string `json:"roomCode, omitempty"`
	Role     string `json:"role, omitempty"`
	Message  string `json:"message"`
}

type Room struct {
	code        string
	playerCount int
}

// Hub  = owner of all live rooms.
type Hub struct {
	mu    sync.Mutex
	rooms map[string]*Room
}

func NewHub() *Hub {
	return &Hub{
		rooms: make(map[string]*Room),
	}
}

func main() {
	hub := NewHub()

	mux := http.NewServeMux()
	mux.HandleFunc("/", handleRoot)
	mux.HandleFunc("/health", handleHealth)
	mux.HandleFunc("/health/db", handleDatabaseHealth)
	mux.HandleFunc("/ws", hub.handleWebSocket)

	port := getenv("PORT", "8080")
	addr := ":" + port

	log.Printf("API server listening on http://localhost%s", addr)
	log.Fatal(http.ListenAndServe(addr, mux))
}

func handleRoot(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	_, _ = w.Write([]byte("ft_transcendence API\n"))
}

func handleHealth(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	_, _ = w.Write([]byte("ok\n"))
}

func (h *Hub) handleWebSocket(w http.ResponseWriter, r *http.Request) {
	opts := &websocket.AcceptOptions{
		OriginPatterns: []string{"http://localhost:5173"},
	}

	conn, err := websocket.Accept(w, r, opts)
	if err != nil {
		log.Printf("websocket upgrade failed: %v", err)
		return
	}
	defer conn.CloseNow()

	log.Print("websocket client connected")

	ctx := context.Background()

	joinedRoom := ""

	for {
		var message ClientMessage

		if err := wsjson.Read(ctx, conn, &message); err != nil {
			log.Print("websocket client disconnected")
			return
		}

		if message.Type != "room.join" || strings.TrimSpace(message.RoomCode) == "" {
			_ = wsjson.Write(ctx, conn, ServerMessage{
				Type:    "error",
				Message: "room.join requires roomCode",
			})
			continue
		}

		if joinedRoom != "" {
			_ = wsjson.Write(ctx, conn, ServerMessage{
				Type:    "error",
				Message: "this connection already joined room " + joinedRoom,
			})
			continue
		}

		roomCode := strings.ToUpper(strings.TrimSpace(message.RoomCode))

		h.mu.Lock()

		room, exists := h.rooms[roomCode]
		if !exists {
			room = &Room{
				code: roomCode,
			}
			h.rooms[roomCode] = room
		}

		if room.playerCount >= 2 {
			h.mu.Unlock()

			_ = wsjson.Write(ctx, conn, ServerMessage{
				Type:    "error",
				Message: "room is full",
			})
			continue
		}

		room.playerCount++

		role := "mission_control"
		if room.playerCount == 2 {
			role = "on_site"
		}

		joinedRoom = room.code
		playerCount := room.playerCount

		h.mu.Unlock()

		log.Printf("player joined room %s as %s (%d/2)", room.code, role,
			playerCount)

		if err := wsjson.Write(ctx, conn, ServerMessage{
			Type:     "room.joined",
			RoomCode: room.code,
			Role:     role,
			Message:  "room joined",
		}); err != nil {

		}
	}
}

func handleDatabaseHealth(w http.ResponseWriter, _ *http.Request) {
	host := getenv("DB_HOST", "localhost")
	port := getenv("DB_PORT", "5432")
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

func getenv(key string, fallback string) string {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return fallback
	}
	return value
}
