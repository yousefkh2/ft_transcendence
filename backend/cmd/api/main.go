package main

import (
	"context"
	"crypto/rand"
	"encoding/hex"
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
	RoomCode string `json:"roomCode,omitempty"`
	ObjectID string `json:"objectId,omitempty"`
	Relation string `json:"relation,omitempty"`
	TargetID string `json:"targetId,omitempty"`
}

type ServerMessage struct {
	Type                string   `json:"type"`
	RoomCode            string   `json:"roomCode,omitempty"`
	PlayerID            string   `json:"playerId,omitempty"`
	Role                string   `json:"role,omitempty"`
	CompletedObjectives []string `json:"completedObjectives,omitempty"`
	Message             string   `json:"message"`
}

type Player struct {
	id   string
	conn *websocket.Conn
}

type Room struct {
	code                string
	players             map[string]*Player
	completedObjectives map[string]bool
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

	playerID, err := newPlayerID()
	if err != nil {
		log.Printf("player ID generation failed: %v", err)
		return
	}

	player := &Player{
		id:   playerID,
		conn: conn,
	}

	log.Printf("websocket client connected: %s", playerID)

	ctx := context.Background()

	joinedRoom := ""
	joinedRole := ""

	defer func() {
		h.leaveRoom(joinedRoom, joinedRole, playerID)
	}()

	for {
		var message ClientMessage

		if err := wsjson.Read(ctx, conn, &message); err != nil {
			log.Printf("websocket client disconnected: %s", playerID)
			return
		}

		switch message.Type {
		case "room.join":
			h.handleRoomJoin(ctx, conn, player, message, &joinedRoom, &joinedRole)
		case "game.object_moved":
			h.handleObjectMoved(ctx, conn, message, joinedRoom, joinedRole)
		default:
			_ = wsjson.Write(ctx, conn, ServerMessage{
				Type:    "error",
				Message: "unknown message type",
			})
		}

	}
}

func (h *Hub) handleRoomJoin(
	ctx context.Context,
	conn *websocket.Conn,
	player *Player,
	message ClientMessage,
	joinedRoom *string,
	joinedRole *string,
) {
	if strings.TrimSpace(message.RoomCode) == "" {
		_ = wsjson.Write(ctx, conn, ServerMessage{
			Type:    "error",
			Message: "room.join requires roomCode",
		})
		return
	}

	if *joinedRoom != "" {
		_ = wsjson.Write(ctx, conn, ServerMessage{
			Type:    "error",
			Message: "this connection already joined room " + *joinedRoom,
		})
		return
	}

	roomCode := strings.ToUpper(strings.TrimSpace(message.RoomCode))

	h.mu.Lock()

	room, exists := h.rooms[roomCode]
	if !exists {
		room = &Room{
			code:                roomCode,
			players:             make(map[string]*Player),
			completedObjectives: make(map[string]bool),
		}
		h.rooms[roomCode] = room
	}

	role := ""

	if room.players["mission_control"] == nil {
		role = "mission_control"
	} else if room.players["on_site"] == nil {
		role = "on_site"
	} else {
		h.mu.Unlock()

		_ = wsjson.Write(ctx, conn, ServerMessage{
			Type:    "error",
			Message: "room is full",
		})
		return
	}

	room.players[role] = player

	*joinedRoom = room.code
	*joinedRole = role
	playerCount := len(room.players)

	h.mu.Unlock()

	log.Printf("player %s joined room %s as %s (%d/2)", player.id, room.code, role, playerCount)

	if err := wsjson.Write(ctx, conn, ServerMessage{
		Type:     "room.joined",
		RoomCode: room.code,
		PlayerID: player.id,
		Role:     role,
		Message:  "room joined",
	}); err != nil {
		log.Printf("websocket response failed: %v", err)
	}
}

func (h *Hub) handleObjectMoved(
	ctx context.Context,
	conn *websocket.Conn,
	message ClientMessage,
	joinedRoom string,
	joinedRole string,
) {
	if joinedRoom == "" {
		_ = wsjson.Write(ctx, conn, ServerMessage{
			Type:    "error",
			Message: "join a room before sending game events",
		})
		return
	}

	if joinedRole != "on_site" {
		_ = wsjson.Write(ctx, conn, ServerMessage{
			Type:    "error",
			Message: "only on_site can move objects",
		})
		return
	}

	if message.ObjectID != "plant" ||
		message.Relation != "right_of" ||
		message.TargetID != "sofa" {
		_ = wsjson.Write(ctx, conn, ServerMessage{
			Type:    "error",
			Message: "move does not complete an objective",
		})
		return
	}

	h.mu.Lock()

	room := h.rooms[joinedRoom]
	if room == nil {
		h.mu.Unlock()

		_ = wsjson.Write(ctx, conn, ServerMessage{
			Type:    "error",
			Message: "room no longer exists",
		})
		return
	}

	room.completedObjectives["plant_right_of_sofa"] = true
	completedObjectives := []string{"plant_right_of_sofa"}

	h.mu.Unlock()

	log.Printf("objective completed in room %s: plant_right_of_sofa", joinedRoom)

	if err := wsjson.Write(ctx, conn, ServerMessage{
		Type:                "game.state_updated",
		RoomCode:            joinedRoom,
		CompletedObjectives: completedObjectives,
		Message:             "objective completed",
	}); err != nil {
		log.Printf("websocket response failed: %v", err)
	}
}

func (h *Hub) leaveRoom(roomCode string, role string, playerID string) {
	if roomCode == "" || role == "" {
		return
	}

	h.mu.Lock()
	defer h.mu.Unlock()

	room, exists := h.rooms[roomCode]
	if !exists {
		return
	}

	player, exists := room.players[role]
	if !exists || player.id != playerID {
		return
	}

	delete(room.players, role)

	if len(room.players) == 0 {
		delete(h.rooms, roomCode)
		log.Printf("room deleted after final player left: %s", roomCode)
		return
	}

	log.Printf("player %s left room %s as %s", playerID, roomCode, role)
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

func newPlayerID() (string, error) {
	bytes := make([]byte, 4) // Four bytes produce eight hexadecimal characters (player_e4a91c8f)

	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}

	return "player_" + hex.EncodeToString(bytes), nil
}

func getenv(key string, fallback string) string {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return fallback
	}
	return value
}
