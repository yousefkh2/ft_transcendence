package hub_test

import (
	"context"
	"encoding/json"
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

func newChatServer(t *testing.T) (*hub.Hub, *httptest.Server) {
	t.Helper()
	h := hub.NewHubWithOrigins([]string{"*"}, newFakeRoleResolver())
	e := echo.New()
	e.GET("/ws/chat", h.HandleChatWebSocket)
	server := httptest.NewServer(e)
	t.Cleanup(server.Close)
	return h, server
}

func dialChat(t *testing.T, server *httptest.Server, userID string) *websocket.Conn {
	t.Helper()
	url := "ws" + strings.TrimPrefix(server.URL, "http") + "/ws/chat"
	conn, _, err := websocket.Dial(context.Background(), url, nil)
	if err != nil {
		t.Fatalf("dial failed: %v", err)
	}
	t.Cleanup(func() { _ = conn.Close(websocket.StatusNormalClosure, "") })
	if userID == "" {
		return conn
	}
	if err := wsjson.Write(context.Background(), conn, model.ClientMessage{Type: "chat.auth", Token: tokenFor(t, userID)}); err != nil {
		t.Fatalf("write auth: %v", err)
	}
	if got := readChat(t, conn); got.Type != "chat.ready" {
		t.Fatalf("expected chat.ready, got %+v", got)
	}
	return conn
}

func readChat(t *testing.T, conn *websocket.Conn) model.ServerMessage {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	var message model.ServerMessage
	if err := wsjson.Read(ctx, conn, &message); err != nil {
		t.Fatalf("read: %v", err)
	}
	return message
}

func TestChat_PushReachesEveryTabOfTheTargetUsersOnly(t *testing.T) {
	h, server := newChatServer(t)
	aliceTab1 := dialChat(t, server, "alice")
	aliceTab2 := dialChat(t, server, "alice")
	bob := dialChat(t, server, "bob")
	carol := dialChat(t, server, "carol")

	payload, _ := json.Marshal(map[string]string{"body": "hallo"})
	h.PushChat([]string{"alice", "bob"}, model.ServerMessage{Type: "chat.message", Payload: payload})

	for name, conn := range map[string]*websocket.Conn{"alice tab 1": aliceTab1, "alice tab 2": aliceTab2, "bob": bob} {
		got := readChat(t, conn)
		if got.Type != "chat.message" || !strings.Contains(string(got.Payload), "hallo") {
			t.Errorf("%s: got %+v", name, got)
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()
	var unexpected model.ServerMessage
	if err := wsjson.Read(ctx, carol, &unexpected); err == nil {
		t.Errorf("carol should not receive the message, got %+v", unexpected)
	}
}

func TestChat_RejectsInvalidToken(t *testing.T) {
	_, server := newChatServer(t)
	conn := dialChat(t, server, "")
	if err := wsjson.Write(context.Background(), conn, model.ClientMessage{Type: "chat.auth", Token: "garbage"}); err != nil {
		t.Fatalf("write: %v", err)
	}
	if got := readChat(t, conn); got.Type != "error" {
		t.Fatalf("expected error, got %+v", got)
	}
}

func TestChat_ClosedTabIsForgotten(t *testing.T) {
	h, server := newChatServer(t)
	closed := dialChat(t, server, "dave")
	open := dialChat(t, server, "dave")
	_ = closed.Close(websocket.StatusNormalClosure, "")
	time.Sleep(100 * time.Millisecond)

	// Pushing to a user whose other tab is gone must still reach the open one.
	h.PushChat([]string{"dave"}, model.ServerMessage{Type: "chat.message"})
	if got := readChat(t, open); got.Type != "chat.message" {
		t.Fatalf("got %+v", got)
	}
}
