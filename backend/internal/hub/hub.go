package hub

import (
	"log"
	"net/http"
	"context"
	"crypto/rand"
	"encoding/hex"
	"strings"
	"time"
	"sync"

	"github.com/coder/websocket"
	"github.com/coder/websocket/wsjson"
	"transcendence/backend/internal/model"
	"transcendence/backend/internal/game"
)

type Player struct {
	id      string
	conn    *websocket.Conn
	writeMu sync.Mutex
}

type Room struct {
	code                string
	players             map[string]*Player
	completedObjectives map[string]bool
	objectPositions     map[string]model.Position
	roundDeadline       time.Time
}

// Hub  = owner of all live rooms.
type Hub struct {
	mu		sync.Mutex
	rooms	map[string]*Room
	
}

func NewHub() *Hub {
	return &Hub{
		rooms: make(map[string]*Room),
	}
}

func (h *Hub) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
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
		var message model.ClientMessage

		if err := wsjson.Read(ctx, conn, &message); err != nil {
			log.Printf("websocket client disconnected: %s", playerID)
			return
		}

		switch message.Type {
		case "room.join":
			h.handleRoomJoin(ctx, conn, player, message, &joinedRoom, &joinedRole)
		case "game.object_moved":
			h.handleObjectMoved(ctx, player, message, joinedRoom, joinedRole)
		case "voice.transcript":
			h.handleVoiceTranscript(ctx, player, message, joinedRoom, joinedRole)
		default:
			_ = wsjson.Write(ctx, conn, model.ServerMessage{
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
	message model.ClientMessage,
	joinedRoom *string,
	joinedRole *string,
) {
	if strings.TrimSpace(message.RoomCode) == "" {
		_ = wsjson.Write(ctx, conn, model.ServerMessage{
			Type:    "error",
			Message: "room.join requires roomCode",
		})
		return
	}

	if *joinedRoom != "" {
		_ = wsjson.Write(ctx, conn, model.ServerMessage{
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
			objectPositions:     game.InitialObjectPositions(),
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

		_ = wsjson.Write(ctx, conn, model.ServerMessage{
			Type:    "error",
			Message: "room is full",
		})
		return
	}

	room.players[role] = player

	if len(room.players) == 2 && room.roundDeadline.IsZero() {
		room.roundDeadline = time.Now().Add(game.RoundDuration)
	}

	*joinedRoom = room.code
	*joinedRole = role
	playerCount := len(room.players)

	completedObjectives := game.CompletedObjectiveIDs(room.completedObjectives)
	objectPositions := game.CopyObjectPositions(room.objectPositions)
	remaining := game.RemainingSeconds(room.roundDeadline)

	h.mu.Unlock()

	log.Printf("player %s joined room %s as %s (%d/2)", player.id, room.code, role, playerCount)

	if err := wsjson.Write(ctx, conn, model.ServerMessage{
		Type:                "room.joined",
		RoomCode:            room.code,
		PlayerID:            player.id,
		Role:                role,
		CompletedObjectives: completedObjectives,
		ObjectPositions:     objectPositions,
		Message:             "room joined",
		RemainingSeconds:    remaining,
	}); err != nil {
		log.Printf("websocket response failed: %v", err)
	}
}

func (h *Hub) handleObjectMoved(
	ctx context.Context,
	player *Player,
	message model.ClientMessage,
	joinedRoom string,
	joinedRole string,
) {
	if joinedRoom == "" {
		_ = player.send(ctx, model.ServerMessage{
			Type:    "error",
			Message: "join a room before sending game events",
		})
		return
	}

	if joinedRole != "on_site" {
		_ = player.send(ctx, model.ServerMessage{
			Type:    "error",
			Message: "only on_site can move objects",
		})
		return
	}

	h.mu.Lock()

	room := h.rooms[joinedRoom]
	if room == nil {
		h.mu.Unlock()

		_ = player.send(ctx, model.ServerMessage{
			Type:    "error",
			Message: "room no longer exists",
		})
		return
	}

	if game.RoundExpired(room.roundDeadline) {
		completedObjectives := game.CompletedObjectiveIDs(room.completedObjectives)
		objectPositions := game.CopyObjectPositions(room.objectPositions)
		remaining := game.RemainingSeconds(room.roundDeadline)

		players := make([]*Player, 0, len(room.players))
		for _, roomPlayer := range room.players {
			players = append(players, roomPlayer)
		}

		h.mu.Unlock()

		expiredMessage := model.ServerMessage{
			Type:                "game.round_expired",
			RoomCode:            joinedRoom,
			CompletedObjectives: completedObjectives,
			ObjectPositions:     objectPositions,
			RemainingSeconds:    remaining,
			Message:             "round expired",
		}

		for _, roomPlayer := range players {
			if err := roomPlayer.send(ctx, expiredMessage); err != nil {
				log.Printf("round expiry update failed for player %s: %v", roomPlayer.id,
					err)
			}
		}

		return
	}

	room.objectPositions[message.ObjectID] = model.Position{
		X: message.X,
		Y: message.Y,
	}

	room.completedObjectives =
		game.CompletedObjectivesFromPositions(room.objectPositions)

	completedObjectives := game.CompletedObjectiveIDs(room.completedObjectives)
	objectPositions := game.CopyObjectPositions(room.objectPositions)
	remaining := game.RemainingSeconds(room.roundDeadline)

	messageType := "game.state_updated"
	stateMessage := "object moved"

	if game.AllObjectivesCompleted(room.completedObjectives) {
		messageType = "game.round_completed"
		stateMessage = "round completed"
	}

	players := make([]*Player, 0, len(room.players))
	for _, roomPlayer := range room.players {
		players = append(players, roomPlayer)
	}

	h.mu.Unlock()

	log.Printf("object moved in room %s: %s to (%d,%d)", joinedRoom,
		message.ObjectID, message.X, message.Y)

	stateUpdate := model.ServerMessage{
		Type:                messageType,
		RoomCode:            joinedRoom,
		CompletedObjectives: completedObjectives,
		Message:             stateMessage,
		ObjectPositions:     objectPositions,
		RemainingSeconds:    remaining,
	}

	for _, roomPlayer := range players {
		if err := roomPlayer.send(ctx, stateUpdate); err != nil {
			log.Printf("state update failed for player %s: %v", roomPlayer.id, err)
		}
	}
}

func (h *Hub) handleVoiceTranscript(
	ctx context.Context,
	player *Player,
	message model.ClientMessage,
	joinedRoom string,
	joinedRole string,
) {
	if joinedRoom == "" {
		_ = player.send(ctx, model.ServerMessage{
			Type:    "error",
			Message: "join a room before sending transcript events",
		})
		return
	}
	text := strings.TrimSpace(message.Text)
	if text == "" {
		return
	}

	h.mu.Lock()

	room := h.rooms[joinedRoom]
	if room == nil {
		h.mu.Unlock()

		_ = player.send(ctx, model.ServerMessage{
			Type:    "error",
			Message: "room no longer exists",
		})
		return
	}

	players := make([]*Player, 0, len(room.players))
	for _, roomPlayer := range room.players {
		players = append(players, roomPlayer)
	}

	h.mu.Unlock()

	transcriptMessage := model.ServerMessage{
		Type:     "voice.transcript",
		RoomCode: joinedRoom,
		PlayerID: player.id,
		Role:     joinedRole,
		Text:     text,
		IsFinal:  message.IsFinal,
		Message:  "transcript received",
	}

	for _, roomPlayer := range players {
		if err := roomPlayer.send(ctx, transcriptMessage); err != nil {
			log.Printf("transcript broadcast failed for player %s: %v", roomPlayer.id, err)
		}
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

func (p *Player) send(ctx context.Context, message model.ServerMessage) error {
	p.writeMu.Lock()
	defer p.writeMu.Unlock()

	return wsjson.Write(ctx, p.conn, message)
}

func newPlayerID() (string, error) {
	bytes := make([]byte, 4) // Four bytes produce eight hexadecimal characters (player_e4a91c8f)

	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}

	return "player_" + hex.EncodeToString(bytes), nil
}
