package hub

import (
	"context"
	"log"
	"time"

	"transcendence/backend/internal/auth"
	"transcendence/backend/internal/model"

	"github.com/coder/websocket"
	"github.com/coder/websocket/wsjson"
	"github.com/labstack/echo/v4"
)

// Live delivery for the friends chat.
//
// Messages are sent and stored through the REST API (POST /api/chat/:friendId/messages); this socket
// only pushes them to the browser tabs of both users as soon as they are saved. A user can have the
// site open in several tabs, so every tab gets its own entry.
//
//	client -> server  chat.auth {token}
//	server -> client  chat.ready
//	server -> client  chat.message {payload: the stored message}
const chatWriteTimeout = 5 * time.Second

func (h *Hub) HandleChatWebSocket(c echo.Context) error {
	conn, err := websocket.Accept(c.Response().Writer, c.Request(), &websocket.AcceptOptions{OriginPatterns: h.allowedOrigins})
	if err != nil {
		log.Printf("chat websocket upgrade failed: %v", err)
		return nil
	}
	defer conn.CloseNow()

	ctx := context.Background()
	client := &Player{conn: conn}
	userID := ""

	defer func() {
		h.removeChatClient(userID, client)
	}()

	for {
		var message model.ClientMessage
		if err := wsjson.Read(ctx, conn, &message); err != nil {
			return nil
		}

		switch message.Type {
		case "chat.auth":
			if userID != "" {
				_ = client.send(ctx, model.ServerMessage{Type: "error", Message: "already authenticated"})
				continue
			}
			id, err := auth.ParseJWT(message.Token)
			if err != nil {
				_ = client.send(ctx, model.ServerMessage{Type: "error", Message: "invalid or missing token"})
				continue
			}
			userID = id
			h.addChatClient(userID, client)
			_ = client.send(ctx, model.ServerMessage{Type: "chat.ready"})
		default:
			_ = client.send(ctx, model.ServerMessage{Type: "error", Message: "unknown message type"})
		}
	}
}

func (h *Hub) addChatClient(userID string, client *Player) {
	h.chatMu.Lock()
	defer h.chatMu.Unlock()
	if h.chatClients[userID] == nil {
		h.chatClients[userID] = make(map[*Player]struct{})
	}
	h.chatClients[userID][client] = struct{}{}
}

func (h *Hub) removeChatClient(userID string, client *Player) {
	if userID == "" {
		return
	}
	h.chatMu.Lock()
	defer h.chatMu.Unlock()
	delete(h.chatClients[userID], client)
	if len(h.chatClients[userID]) == 0 {
		delete(h.chatClients, userID)
	}
}

// PushChat sends message to every open chat connection of the given users.
// Users without an open connection simply see the message the next time they load the conversation.
func (h *Hub) PushChat(userIDs []string, message model.ServerMessage) {
	h.chatMu.Lock()
	var targets []*Player
	for _, userID := range userIDs {
		for client := range h.chatClients[userID] {
			targets = append(targets, client)
		}
	}
	h.chatMu.Unlock()

	for _, client := range targets {
		ctx, cancel := context.WithTimeout(context.Background(), chatWriteTimeout)
		if err := client.send(ctx, message); err != nil {
			log.Printf("chat push failed: %v", err)
		}
		cancel()
	}
}
