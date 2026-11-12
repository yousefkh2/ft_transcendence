package hub_test

import (
	"context"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/coder/websocket"
	"github.com/coder/websocket/wsjson"
	"github.com/labstack/echo/v4"

	"transcendence/backend/internal/hub"
	"transcendence/backend/internal/model"
)

func newVoiceServer(t *testing.T, roles hub.RoleResolver) string {
	t.Helper()
	h := hub.NewHubWithOrigins([]string{"*"}, roles)
	e := echo.New()
	e.GET("/ws/voice", h.HandleVoiceWebSocket)
	server := httptest.NewServer(e)
	t.Cleanup(server.Close)
	return "ws" + strings.TrimPrefix(server.URL, "http") + "/ws/voice"
}

func dialVoice(t *testing.T, url string) *websocket.Conn {
	t.Helper()
	conn, _, err := websocket.Dial(context.Background(), url, nil)
	if err != nil {
		t.Fatalf("dial failed: %v", err)
	}
	t.Cleanup(func() { _ = conn.Close(websocket.StatusNormalClosure, "") })
	return conn
}

func voiceSend(t *testing.T, conn *websocket.Conn, message model.ClientMessage) {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := wsjson.Write(ctx, conn, message); err != nil {
		t.Fatalf("write %s failed: %v", message.Type, err)
	}
}

func voiceRead(t *testing.T, conn *websocket.Conn) model.ServerMessage {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var message model.ServerMessage
	if err := wsjson.Read(ctx, conn, &message); err != nil {
		t.Fatalf("read failed: %v", err)
	}
	return message
}

func TestVoice_JoinRelayAndLeave(t *testing.T) {
	roles := newFakeRoleResolver()
	alice := uniqueUserID(t, "alice")
	bob := uniqueUserID(t, "bob")
	roles.assign("VOICE", alice, "mission_control")
	roles.assign("VOICE", bob, "on_site")
	url := newVoiceServer(t, roles)

	aliceConn := dialVoice(t, url)
	voiceSend(t, aliceConn, model.ClientMessage{Type: "voice.join", RoomCode: "voice", Token: tokenFor(t, alice)})
	if got := voiceRead(t, aliceConn); got.Type != "voice.joined" || len(got.Peers) != 0 {
		t.Fatalf("expected voice.joined without peers, got %s %v", got.Type, got.Peers)
	}

	bobConn := dialVoice(t, url)
	voiceSend(t, bobConn, model.ClientMessage{Type: "voice.join", RoomCode: "VOICE", Token: tokenFor(t, bob)})
	if got := voiceRead(t, bobConn); got.Type != "voice.joined" || len(got.Peers) != 1 || got.Peers[0] != alice {
		t.Fatalf("expected bob to see alice as peer, got %s %v", got.Type, got.Peers)
	}
	if got := voiceRead(t, aliceConn); got.Type != "voice.peer_joined" || got.FromUser != bob {
		t.Fatalf("expected alice to hear bob joined, got %s from %s", got.Type, got.FromUser)
	}

	// Bob (the newcomer) offers, alice answers.
	voiceSend(t, bobConn, model.ClientMessage{Type: "webrtc.offer", TargetUser: alice, SDP: "offer-sdp"})
	if got := voiceRead(t, aliceConn); got.Type != "webrtc.offer" || got.FromUser != bob || got.SDP != "offer-sdp" {
		t.Fatalf("unexpected relayed offer: %+v", got)
	}
	voiceSend(t, aliceConn, model.ClientMessage{Type: "webrtc.answer", TargetUser: bob, SDP: "answer-sdp"})
	if got := voiceRead(t, bobConn); got.Type != "webrtc.answer" || got.FromUser != alice || got.SDP != "answer-sdp" {
		t.Fatalf("unexpected relayed answer: %+v", got)
	}
	voiceSend(t, aliceConn, model.ClientMessage{Type: "webrtc.ice_candidate", TargetUser: bob, Candidate: `{"candidate":"x"}`})
	if got := voiceRead(t, bobConn); got.Type != "webrtc.ice_candidate" || got.Candidate != `{"candidate":"x"}` {
		t.Fatalf("unexpected relayed candidate: %+v", got)
	}

	_ = bobConn.Close(websocket.StatusNormalClosure, "")
	if got := voiceRead(t, aliceConn); got.Type != "voice.peer_left" || got.FromUser != bob {
		t.Fatalf("expected bob to leave, got %s from %s", got.Type, got.FromUser)
	}
}

func TestVoice_NonMemberIsRejected(t *testing.T) {
	roles := newFakeRoleResolver()
	url := newVoiceServer(t, roles)

	conn := dialVoice(t, url)
	voiceSend(t, conn, model.ClientMessage{Type: "voice.join", RoomCode: "NOPE", Token: tokenFor(t, uniqueUserID(t, "eve"))})
	if got := voiceRead(t, conn); got.Type != "error" {
		t.Fatalf("expected error for a user outside the lobby, got %s", got.Type)
	}

	// Without joining, signals are refused instead of relayed.
	voiceSend(t, conn, model.ClientMessage{Type: "webrtc.offer", SDP: "x"})
	if got := voiceRead(t, conn); got.Type != "error" {
		t.Fatalf("expected error for a signal before joining, got %s", got.Type)
	}
}

func TestVoice_InvalidTokenIsRejected(t *testing.T) {
	roles := newFakeRoleResolver()
	roles.assign("VOICE", "alice", "on_site")
	url := newVoiceServer(t, roles)

	conn := dialVoice(t, url)
	voiceSend(t, conn, model.ClientMessage{Type: "voice.join", RoomCode: "VOICE", Token: "not-a-jwt"})
	if got := voiceRead(t, conn); got.Type != "error" {
		t.Fatalf("expected error for an invalid token, got %s", got.Type)
	}
}
