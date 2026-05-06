package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(_ *http.Request) bool {
		return true
	},
}

type ClientMessage struct {
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

type JoinPayload struct {
	RoomCode string `json:"roomCode"`
}

type MovePayload struct {
	MatchID  string `json:"matchId"`
	ObjectID string `json:"objectId"`
	Relation string `json:"relation"`
	TargetID string `json:"targetId"`
}

type ServerMessage struct {
	Type    string `json:"type"`
	Payload any    `json:"payload"`
}

type Client struct {
	id   string
	role string
	room *Room
	conn *websocket.Conn

	sendMu sync.Mutex
}

type Room struct {
	code                string
	matchID             string
	clients             map[*Client]bool
	completedObjectives map[string]bool
}

type Hub struct {
	mu    sync.Mutex
	rooms map[string]*Room
}

func NewHub() *Hub {
	return &Hub{rooms: make(map[string]*Room)}
}

func main() {
	rand.New(rand.NewSource(time.Now().UnixNano()))

	hub := NewHub()
	mux := http.NewServeMux()
	mux.HandleFunc("/health", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok\n"))
	})
	mux.HandleFunc("/ws", hub.handleWebSocket)
	mux.Handle("/", http.FileServer(http.Dir("frontend/sprint1-poc")))

	port := strings.TrimSpace(getenv("PORT", "8080"))
	addr := ":" + port
	log.Printf("Sprint POC server listening on http://127.0.0.1%s", addr)
	log.Fatal(http.ListenAndServe(addr, mux))
}

func getenv(key string, fallback string) string {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return fallback
	}
	return value
}

func (h *Hub) handleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("websocket upgrade failed: %v", err)
		return
	}

	client := &Client{
		id:   fmt.Sprintf("player_%d", rand.Intn(9000)+1000),
		conn: conn,
	}
	h.readLoop(client)
}

func (h *Hub) readLoop(c *Client) {
	defer h.disconnect(c)

	for {
		messageType, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsCloseError(err, websocket.CloseGoingAway, websocket.CloseNormalClosure) {
				log.Printf("client disconnected: %s", c.id)
			} else {
				log.Printf("read failed for %s: %v", c.id, err)
			}
			return
		}
		if messageType != websocket.TextMessage {
			log.Printf("ignored non-text websocket message from %s", c.id)
			continue
		}

		log.Printf("received from %s: %s", c.id, string(message))

		var clientMessage ClientMessage
		if err := json.Unmarshal(message, &clientMessage); err != nil {
			c.send("error", map[string]string{"message": "invalid JSON message"})
			continue
		}

		switch clientMessage.Type {
		case "room.join":
			h.joinRoom(c, clientMessage.Payload)
		case "game.object_moved":
			h.handleObjectMoved(c, clientMessage.Payload)
		default:
			c.send("error", map[string]string{"message": "unknown message type"})
		}
	}
}

func (h *Hub) joinRoom(c *Client, payload json.RawMessage) {
	var join JoinPayload
	if err := json.Unmarshal(payload, &join); err != nil || strings.TrimSpace(join.RoomCode) == "" {
		c.send("error", map[string]string{"message": "room.join requires roomCode"})
		return
	}

	roomCode := strings.ToUpper(strings.TrimSpace(join.RoomCode))

	h.mu.Lock()
	room := h.rooms[roomCode]
	if room == nil {
		room = &Room{
			code:                roomCode,
			matchID:             "match_" + strings.ToLower(roomCode),
			clients:             make(map[*Client]bool),
			completedObjectives: make(map[string]bool),
		}
		h.rooms[roomCode] = room
	}

	c.room = room
	c.role = nextRole(room)
	room.clients[c] = true
	playerCount := len(room.clients)
	h.mu.Unlock()

	log.Printf("client %s joined room %s as %s", c.id, room.code, c.role)

	c.send("room.role_assigned", map[string]string{
		"role":     c.role,
		"playerId": c.id,
	})

	c.send("game.round_started", map[string]any{
		"matchId":         room.matchID,
		"scenarioId":      "apartment_easy_001",
		"durationSeconds": 180,
		"role":            c.role,
	})

	h.broadcast(room, "room.player_count", map[string]any{
		"roomCode":    room.code,
		"playerCount": playerCount,
	})
}

func (h *Hub) handleObjectMoved(c *Client, payload json.RawMessage) {
	if c.room == nil {
		c.send("error", map[string]string{"message": "join a room before sending game events"})
		return
	}
	if c.role != "on_site" {
		log.Printf("rejected move from %s: role %s cannot move", c.id, c.role)
		c.send("error", map[string]string{"message": "only on_site can move objects"})
		return
	}

	var move MovePayload
	if err := json.Unmarshal(payload, &move); err != nil {
		c.send("error", map[string]string{"message": "invalid object move payload"})
		return
	}

	objectiveID, ok := objectiveFor(move)
	if !ok {
		log.Printf("rejected move from %s: %s %s %s", c.id, move.ObjectID, move.Relation, move.TargetID)
		h.mu.Lock()
		done := completed(c.room)
		h.mu.Unlock()

		c.send("game.state_patch", map[string]any{
			"matchId":             c.room.matchID,
			"completedObjectives": done,
			"accepted":            false,
		})
		return
	}

	h.mu.Lock()
	c.room.completedObjectives[objectiveID] = true
	done := completed(c.room)
	success := len(done) == 4
	h.mu.Unlock()

	log.Printf("accepted move from %s: completed %s", c.id, objectiveID)

	h.broadcast(c.room, "game.state_patch", map[string]any{
		"matchId":              c.room.matchID,
		"completedObjectives":  done,
		"timeRemainingSeconds": 180,
		"accepted":             true,
	})

	if success {
		h.broadcast(c.room, "game.round_ended", map[string]any{
			"matchId":             c.room.matchID,
			"result":              "success",
			"completedObjectives": done,
			"score":               1.0,
		})
	}
}

func (h *Hub) broadcast(room *Room, messageType string, payload any) {
	h.mu.Lock()
	clients := make([]*Client, 0, len(room.clients))
	for client := range room.clients {
		clients = append(clients, client)
	}
	h.mu.Unlock()

	for _, client := range clients {
		client.send(messageType, payload)
	}
}

func (h *Hub) disconnect(c *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if c.room != nil {
		delete(c.room.clients, c)
		if len(c.room.clients) == 0 {
			delete(h.rooms, c.room.code)
		}
	}
	_ = c.conn.Close()
}

func (c *Client) send(messageType string, payload any) {
	c.sendMu.Lock()
	defer c.sendMu.Unlock()

	message, err := json.Marshal(ServerMessage{Type: messageType, Payload: payload})
	if err != nil {
		return
	}
	if err := c.conn.WriteMessage(websocket.TextMessage, message); err != nil {
		log.Printf("send failed for %s: %v", c.id, err)
	}
}

func nextRole(room *Room) string {
	for client := range room.clients {
		if client.role == "mission_control" {
			return "on_site"
		}
	}
	return "mission_control"
}

func objectiveFor(move MovePayload) (string, bool) {
	targets := map[MovePayload]string{
		{ObjectID: "plant", Relation: "right_of", TargetID: "sofa"}:         "objective_1",
		{ObjectID: "red_book", Relation: "on", TargetID: "coffee_table"}:    "objective_2",
		{ObjectID: "blue_cup", Relation: "under", TargetID: "coffee_table"}: "objective_3",
		{ObjectID: "picture", Relation: "above", TargetID: "sofa"}:          "objective_4",
	}
	id, ok := targets[MovePayload{ObjectID: move.ObjectID, Relation: move.Relation, TargetID: move.TargetID}]
	return id, ok
}

func completed(room *Room) []string {
	ids := []string{"objective_1", "objective_2", "objective_3", "objective_4"}
	done := make([]string, 0, len(ids))
	for _, id := range ids {
		if room.completedObjectives[id] {
			done = append(done, id)
		}
	}
	return done
}
