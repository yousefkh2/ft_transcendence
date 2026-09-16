package hub

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"log"
	"strings"
	"sync"
	"time"

	"transcendence/backend/internal/auth"
	"transcendence/backend/internal/db"
	"transcendence/backend/internal/game"
	"transcendence/backend/internal/model"

	"github.com/coder/websocket"
	"github.com/coder/websocket/wsjson"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
)

type RoleResolver interface {
	GetParticipantRole(ctx context.Context, code, userID string) (sessionID, role string, err error)
}

type pooledRoleResolver struct {
	pool *pgxpool.Pool
}

func (r *pooledRoleResolver) GetParticipantRole(ctx context.Context, code, userID string) (string, string, error) {
	return db.GetParticipantRole(ctx, r.pool, code, userID)
}

type Player struct {
	id		string
	conn	*websocket.Conn
	writeMu	sync.Mutex
}

type Room struct {
	code				string
	players				map[string]*Player
	completedObjectives	map[string]bool
	objectPositions		map[string]model.Position
	roundDeadline		time.Time
	startedAt			time.Time
}

// Hub  = owner of all live rooms.
type Hub struct {
	mu				sync.Mutex
	rooms			map[string]*Room
	allowedOrigins	[]string
	db				*pgxpool.Pool
	roles			RoleResolver
}

func NewHub(pool *pgxpool.Pool) *Hub {
	return &Hub{
		rooms: make(map[string]*Room),
		allowedOrigins: []string{"localhost:5173"},
		db: pool,
		roles: &pooledRoleResolver{pool: pool},
	}
}

func NewHubWithOrigins(origins []string, roles RoleResolver) *Hub {
	return &Hub{
		rooms: make(map[string]*Room),
		allowedOrigins: origins,
		roles: roles,
	}
}

func (h *Hub) HandleWebSocket(c echo.Context) error {
	conn, err := websocket.Accept(c.Response().Writer, c.Request(), &websocket.AcceptOptions{OriginPatterns: h.allowedOrigins})
	if err != nil {
		log.Printf("websocket upgrade failed: %v", err)
		return nil
	}
	defer conn.CloseNow()

	playerID, err := newPlayerID()
	if err != nil {
		log.Printf("player ID generation failed: %v", err)
		return nil
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
			return nil
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
	userID, err := auth.ParseJWT(message.Token)
	if err != nil {
		_ = wsjson.Write(ctx, conn, model.ServerMessage{
			Type: "error",
			Message: "invalid or missing token",
		})
		return
	}

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

	_, role, err := h.roles.GetParticipantRole(ctx, roomCode, userID)
	if err != nil {
		msg := "could not join room"
		switch {
		case errors.Is(err, db.ErrNotInLobby):
			msg = "you have not joined this lobby via the API"
		case errors.Is(err, db.ErrLobbyNotStarted):
			msg = "lobby is not full yet"
		}
		_ = wsjson.Write(ctx, conn, model.ServerMessage{
			Type:		"error",
			Message:	msg,
		})
		return
	}
	if role != db.RoleMissionControl && role != db.RoleOnSite {
		_ = wsjson.Write(ctx, conn, model.ServerMessage{
			Type: "error",
			Message: "invalid role assigned",
		})
		return
	}

	h.mu.Lock()

	room, exists := h.rooms[roomCode]
	if !exists {
		room = &Room{
			code:					roomCode,
			players:				make(map[string]*Player),
			completedObjectives:	make(map[string]bool),
			objectPositions:		game.InitialObjectPositions(),
		}
		h.rooms[roomCode] = room
	}

	if room.players[role] != nil {
		h.mu.Unlock()

		_ = wsjson.Write(ctx, conn, model.ServerMessage{
			Type:		"error",
			Message:	"this role is already connected",
		})
		return
	}

	room.players[role] = player

	if len(room.players) == 2 && room.roundDeadline.IsZero() {
		room.roundDeadline = time.Now().Add(game.RoundDuration)
		room.startedAt = time.Now()
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

	if joinedRole != db.RoleOnSite {
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
		startedAt := room.startedAt

		players := make([]*Player, 0, len(room.players))
		for _, roomPlayer := range room.players {
			players = append(players, roomPlayer)
		}

		h.mu.Unlock()

		go func() {
			if err := db.FinishMatch(context.Background(), h.db, joinedRoom, startedAt, time.Now(), false); err != nil {
				log.Printf("failed to save expired match for room %s: %v", joinedRoom, err)
			}
		}()
		
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

	var startedAt time.Time

	if game.AllObjectivesCompleted(room.completedObjectives) {
		messageType = "game.round_completed"
		stateMessage = "round completed"
		startedAt = room.startedAt
	}

	players := make([]*Player, 0, len(room.players))
	for _, roomPlayer := range room.players {
		players = append(players, roomPlayer)
	}

	h.mu.Unlock()

	if messageType == "game.round_completed" {
		go func() {
			if err := db.FinishMatch(context.Background(), h.db, joinedRoom, startedAt, time.Now(), true); err != nil {
				log.Printf("failed to save completed match for room %s: %v", joinedRoom, err)
			}
		}()
	}
	
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
