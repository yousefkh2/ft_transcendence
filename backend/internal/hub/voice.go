package hub

import (
	"context"
	"log"
	"strings"

	"transcendence/backend/internal/auth"
	"transcendence/backend/internal/model"

	"github.com/coder/websocket"
	"github.com/coder/websocket/wsjson"
	"github.com/labstack/echo/v4"
)

// Voice chat between the players of a lobby.
//
// The audio itself goes directly from browser to browser (WebRTC); the server only passes the
// connection setup messages between the players and remembers whose microphone is connected,
// because a player may only ready up with a working microphone.
//
//	client -> server  voice.join {token, roomCode}
//	server -> client  voice.joined {peers}            users already in the voice chat: send each an offer
//	server -> others  voice.peer_joined {fromUser}    wait for the newcomer's offer
//	client <-> client webrtc.offer / webrtc.answer / webrtc.ice_candidate {targetUser} -> {fromUser}
//	client -> server  voice.status {connected}        microphone live and connected to the other player
//	server -> others  voice.peer_left {fromUser}
//
// The voice connection is separate from the lobby and game connections, so the call keeps going
// while the game switches from the lobby to the match.
func (h *Hub) HandleVoiceWebSocket(c echo.Context) error {
	conn, err := websocket.Accept(c.Response().Writer, c.Request(), &websocket.AcceptOptions{OriginPatterns: h.allowedOrigins})
	if err != nil {
		log.Printf("voice websocket upgrade failed: %v", err)
		return nil
	}
	defer conn.CloseNow()

	ctx := context.Background()
	peer := &Player{conn: conn}
	code := ""

	defer func() {
		h.leaveVoice(code, peer)
	}()

	for {
		var message model.ClientMessage
		if err := wsjson.Read(ctx, conn, &message); err != nil {
			return nil
		}

		switch message.Type {
		case "voice.join":
			h.handleVoiceJoin(ctx, peer, message, &code)
		case "webrtc.offer", "webrtc.answer", "webrtc.ice_candidate":
			h.relayVoiceSignal(ctx, peer, message, code)
		case "voice.status":
			h.setVoiceConnected(code, peer, message.Connected)
		default:
			_ = peer.send(ctx, model.ServerMessage{Type: "error", Message: "unknown message type"})
		}
	}
}

func (h *Hub) handleVoiceJoin(ctx context.Context, peer *Player, message model.ClientMessage, joinedCode *string) {
	userID, err := auth.ParseJWT(message.Token)
	if err != nil {
		_ = peer.send(ctx, model.ServerMessage{Type: "error", Message: "invalid or missing token"})
		return
	}
	if *joinedCode != "" {
		_ = peer.send(ctx, model.ServerMessage{Type: "error", Message: "already in a voice chat"})
		return
	}

	code := strings.ToUpper(strings.TrimSpace(message.RoomCode))
	if code == "" {
		_ = peer.send(ctx, model.ServerMessage{Type: "error", Message: "voice.join requires roomCode"})
		return
	}
	member, err := h.isLobbyMember(ctx, code, userID)
	if err != nil || !member {
		_ = peer.send(ctx, model.ServerMessage{Type: "error", Message: "you have not joined this lobby via the API"})
		return
	}

	h.mu.Lock()
	if h.voicePeers[code] == nil {
		h.voicePeers[code] = make(map[string]*Player)
	}
	// The same user opening the game again (reload, second tab) replaces the old connection.
	previous := h.voicePeers[code][userID]
	peer.id = userID
	h.voicePeers[code][userID] = peer
	delete(h.voiceConnected[code], userID)
	others := h.otherVoicePeersLocked(code, userID)
	h.mu.Unlock()

	if previous != nil {
		_ = previous.conn.Close(websocket.StatusPolicyViolation, "voice chat opened somewhere else")
	}

	*joinedCode = code
	peerIDs := make([]string, 0, len(others))
	for _, other := range others {
		peerIDs = append(peerIDs, other.id)
	}
	_ = peer.send(ctx, model.ServerMessage{Type: "voice.joined", RoomCode: code, PlayerID: userID, Peers: peerIDs, Message: "voice joined"})
	for _, other := range others {
		_ = other.send(ctx, model.ServerMessage{Type: "voice.peer_joined", RoomCode: code, FromUser: userID, Message: "player joined voice"})
	}
	log.Printf("user %s joined voice chat of lobby %s", userID, code)
}

// relayVoiceSignal passes an offer, answer or ICE candidate to the other player.
func (h *Hub) relayVoiceSignal(ctx context.Context, peer *Player, message model.ClientMessage, code string) {
	if code == "" {
		_ = peer.send(ctx, model.ServerMessage{Type: "error", Message: "join the voice chat first"})
		return
	}

	h.mu.Lock()
	var targets []*Player
	if message.TargetUser != "" {
		if target := h.voicePeers[code][message.TargetUser]; target != nil && target != peer {
			targets = append(targets, target)
		}
	} else {
		targets = h.otherVoicePeersLocked(code, peer.id)
	}
	h.mu.Unlock()

	for _, target := range targets {
		if err := target.send(ctx, model.ServerMessage{
			Type:      message.Type,
			RoomCode:  code,
			SDP:       message.SDP,
			Candidate: message.Candidate,
			FromUser:  peer.id,
		}); err != nil {
			log.Printf("voice signal relay failed for user %s: %v", target.id, err)
		}
	}
}

// setVoiceConnected records whether a player's microphone works. Losing it also un-readies the
// player, which stops a running countdown.
func (h *Hub) setVoiceConnected(code string, peer *Player, connected bool) {
	if code == "" {
		return
	}

	h.mu.Lock()
	if h.voicePeers[code][peer.id] != peer {
		h.mu.Unlock()
		return
	}
	changed := h.voiceConnected[code][peer.id] != connected
	if connected {
		if h.voiceConnected[code] == nil {
			h.voiceConnected[code] = make(map[string]bool)
		}
		h.voiceConnected[code][peer.id] = true
	} else {
		delete(h.voiceConnected[code], peer.id)
		delete(h.lobbyReady[code], peer.id)
	}
	h.mu.Unlock()

	if changed {
		h.PublishLobbyUpdate(context.Background(), code)
	}
}

func (h *Hub) leaveVoice(code string, peer *Player) {
	if code == "" {
		return
	}

	h.mu.Lock()
	if h.voicePeers[code][peer.id] != peer {
		// Already replaced by a newer connection of the same user.
		h.mu.Unlock()
		return
	}
	delete(h.voicePeers[code], peer.id)
	if len(h.voicePeers[code]) == 0 {
		delete(h.voicePeers, code)
	}
	delete(h.voiceConnected[code], peer.id)
	if len(h.voiceConnected[code]) == 0 {
		delete(h.voiceConnected, code)
	}
	delete(h.lobbyReady[code], peer.id)
	others := h.otherVoicePeersLocked(code, peer.id)
	h.mu.Unlock()

	for _, other := range others {
		_ = other.send(context.Background(), model.ServerMessage{Type: "voice.peer_left", RoomCode: code, FromUser: peer.id, Message: "player left voice"})
	}
	h.PublishLobbyUpdate(context.Background(), code)
	log.Printf("user %s left voice chat of lobby %s", peer.id, code)
}

func (h *Hub) otherVoicePeersLocked(code, userID string) []*Player {
	others := make([]*Player, 0, len(h.voicePeers[code]))
	for id, other := range h.voicePeers[code] {
		if id != userID {
			others = append(others, other)
		}
	}
	return others
}
