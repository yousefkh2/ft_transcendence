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
	X        int    `json:"x,omitempty"`
	Y        int    `json:"y,omitempty"`
}

type ServerMessage struct {
	Type                string              `json:"type"`
	RoomCode            string              `json:"roomCode,omitempty"`
	PlayerID            string              `json:"playerId,omitempty"`
	Role                string              `json:"role,omitempty"`
	CompletedObjectives []string            `json:"completedObjectives,omitempty"`
	Message             string              `json:"message"`
	ObjectPositions     map[string]Position `json:"objectPositions,omitempty"`
}

type Objective struct {
	ID       string
	ObjectID string
	Relation string
	TargetID string
}

type Player struct {
	id      string
	conn    *websocket.Conn
	writeMu sync.Mutex
}

type Room struct {
	code                string
	players             map[string]*Player
	completedObjectives map[string]bool
	objectPositions     map[string]Position
}

// Hub  = owner of all live rooms.
type Hub struct {
	mu    sync.Mutex
	rooms map[string]*Room
}

type Position struct {
	X int `json:"x"`
	Y int `json:"y"`
}

var apartmentObjectives = []Objective{
	{
		ID:       "plant_right_of_sofa",
		ObjectID: "plant",
		Relation: "right_of",
		TargetID: "sofa",
	},
	{
		ID:       "lamp_left_of_sofa",
		ObjectID: "lamp",
		Relation: "left_of",
		TargetID: "sofa",
	},
	{
		ID:       "table_in_front_of_sofa",
		ObjectID: "table",
		Relation: "in_front_of",
		TargetID: "sofa",
	},
}

func (p *Player) send(ctx context.Context, message ServerMessage) error {
	p.writeMu.Lock()
	defer p.writeMu.Unlock()

	return wsjson.Write(ctx, p.conn, message)
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
			h.handleObjectMoved(ctx, player, message, joinedRoom, joinedRole)
		default:
			_ = wsjson.Write(ctx, conn, ServerMessage{
				Type:    "error",
				Message: "unknown message type",
			})
		}

	}
}

func initialObjectPositions() map[string]Position {
	return map[string]Position{
		"sofa":  {X: 0, Y: 1},
		"plant": {X: 3, Y: 1},
		"lamp":  {X: 2, Y: 3},
		"table": {X: 4, Y: 4},
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
			objectPositions:     initialObjectPositions(),
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

	completedObjectives := completedObjectiveIDs(room.completedObjectives)
	objectPositions := copyObjectPositions(room.objectPositions)

	h.mu.Unlock()

	log.Printf("player %s joined room %s as %s (%d/2)", player.id, room.code, role, playerCount)

	if err := wsjson.Write(ctx, conn, ServerMessage{
		Type:                "room.joined",
		RoomCode:            room.code,
		PlayerID:            player.id,
		Role:                role,
		CompletedObjectives: completedObjectives,
		ObjectPositions:     objectPositions,
		Message:             "room joined",
	}); err != nil {
		log.Printf("websocket response failed: %v", err)
	}
}

func (h *Hub) handleObjectMoved(
	ctx context.Context,
	player *Player,
	message ClientMessage,
	joinedRoom string,
	joinedRole string,
) {
	if joinedRoom == "" {
		_ = player.send(ctx, ServerMessage{
			Type:    "error",
			Message: "join a room before sending game events",
		})
		return
	}

	if joinedRole != "on_site" {
		_ = player.send(ctx, ServerMessage{
			Type:    "error",
			Message: "only on_site can move objects",
		})
		return
	}

	h.mu.Lock()

	room := h.rooms[joinedRoom]
	if room == nil {
		h.mu.Unlock()

		_ = player.send(ctx, ServerMessage{
			Type:    "error",
			Message: "room no longer exists",
		})
		return
	}

	room.objectPositions[message.ObjectID] = Position{
		X: message.X,
		Y: message.Y,
	}

	room.completedObjectives =
		completedObjectivesFromPositions(room.objectPositions)

	completedObjectives := completedObjectiveIDs(room.completedObjectives)

	objectPositions := copyObjectPositions(room.objectPositions)

	players := make([]*Player, 0, len(room.players))
	for _, roomPlayer := range room.players {
		players = append(players, roomPlayer)
	}

	h.mu.Unlock()

	stateMessage := "object moved"

	log.Printf("object moved in room %s: %s to (%d,%d)", joinedRoom,
		message.ObjectID, message.X, message.Y)

	stateUpdate := ServerMessage{
		Type:                "game.state_updated",
		RoomCode:            joinedRoom,
		CompletedObjectives: completedObjectives,
		Message:             stateMessage,
		ObjectPositions:     objectPositions,
	}

	for _, roomPlayer := range players {
		if err := roomPlayer.send(ctx, stateUpdate); err != nil {
			log.Printf("state update failed for player %s: %v", roomPlayer.id, err)
		}
	}
}

func copyObjectPositions(positions map[string]Position) map[string]Position {
	copied := make(map[string]Position, len(positions))
	for objectID, position := range positions {
		copied[objectID] = position
	}

	return copied
}

func matchingObjective(message ClientMessage) (Objective, bool) {
	for _, objective := range apartmentObjectives {
		if message.ObjectID == objective.ObjectID &&
			message.Relation == objective.Relation &&
			message.TargetID == objective.TargetID {
			return objective, true
		}
	}

	return Objective{}, false
}

func relationMatches(object Position, relation string, target Position) bool {
	switch relation {
	case "right_of":
		return object.Y == target.Y && object.X == target.X+1
	case "left_of":
		return object.Y == target.Y && object.X == target.X-1
	case "in_front_of":
		return object.X == target.X && object.Y == target.Y+1
	default:
		return false
	}
}

func completedObjectivesFromPositions(positions map[string]Position) map[string]bool {
	completed := make(map[string]bool)

	for _, objective := range apartmentObjectives {
		objectPosition, objectExists := positions[objective.ObjectID]
		targetPosition, targetExists := positions[objective.TargetID]
		if !objectExists || !targetExists {
			continue
		}
		if relationMatches(objectPosition, objective.Relation, targetPosition) {
			completed[objective.ID] = true
		}
	}

	return completed
}

func completedObjectiveIDs(completed map[string]bool) []string {
	ids := make([]string, 0, len(completed))
	for _, objective := range apartmentObjectives {
		if completed[objective.ID] {
			ids = append(ids, objective.ID)
		}
	}

	return ids
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
