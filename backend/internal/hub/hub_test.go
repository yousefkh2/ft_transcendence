package hub_test

import (
	"context"
	"fmt"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/coder/websocket"
	"github.com/coder/websocket/wsjson"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"

	"transcendence/backend/internal/hub"
	"transcendence/backend/internal/model"
)

func newTestServer(t *testing.T) *httptest.Server {
	t.Helper()

	h := hub.NewHubWithOrigins([]string{"*"})

	e := echo.New()
	e.GET("/ws", h.HandleWebSocket)

	server := httptest.NewServer(e)
	t.Cleanup(server.Close)

	return server
}

func dial(t *testing.T, server *httptest.Server) *websocket.Conn {
	t.Helper()

	wsURL := "ws" + strings.TrimPrefix(server.URL, "http") + "/ws"

	conn, _, err := websocket.Dial(context.Background(), wsURL, nil)
	if err != nil {
		t.Fatalf("dial failed: %v", err)
	}
	t.Cleanup(func() {
		_ = conn.Close(websocket.StatusNormalClosure, "")
	})

	return conn
}

func generateTestToken(t *testing.T) string {
	t.Helper()

	claims := jwt.MapClaims{
		"sub": fmt.Sprintf("test-user-%d", time.Now().UnixNano()),
		"exp": time.Now().Add(time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signed, err := token.SignedString([]byte("test_secret_for_hub_tests"))
	if err != nil {
		t.Fatalf("failed to sign test token: %v", err)
	}

	return signed
}

func joinRoom(t *testing.T, conn *websocket.Conn, roomCode string) model.ServerMessage {
	t.Helper()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := wsjson.Write(ctx, conn, model.ClientMessage{
		Type:     "room.join",
		RoomCode: roomCode,
		Token: generateTestToken(t),
	}); err != nil {
		t.Fatalf("write room.join failed: %v", err)
	}

	var response model.ServerMessage
	if err := wsjson.Read(ctx, conn, &response); err != nil {
		t.Fatalf("read room.join response failed: %v", err)
	}

	return response
}

func TestRoomJoin_FirstPlayerGetsMissionControl(t *testing.T) {
	server := newTestServer(t)
	conn := dial(t, server)

	response := joinRoom(t, conn, "ABCD")

	if response.Type != "room.joined" {
		t.Fatalf("expected room.joined, got %s (%s)", response.Type, response.Message)
	}
	if response.Role != "mission_control" {
		t.Fatalf("expected role mission_control, got %s", response.Role)
	}
}

func TestRoomJoin_SecondPlayerGetsOnSite(t *testing.T) {
	server := newTestServer(t)

	conn1 := dial(t, server)
	conn2 := dial(t, server)

	first := joinRoom(t, conn1, "ABCD")
	second := joinRoom(t, conn2, "ABCD")

	if first.Role != "mission_control" {
		t.Fatalf("expected first player role mission_control, got %s", first.Role)
	}
	if second.Role != "on_site" {
		t.Fatalf("expected second player role on_site, got %s", second.Role)
	}
}

func TestRoomJoin_DistinctPlayerIDs(t *testing.T) {
	server := newTestServer(t)

	conn1 := dial(t, server)
	conn2 := dial(t, server)

	first := joinRoom(t, conn1, "ABCD")
	second := joinRoom(t, conn2, "ABCD")

	if first.PlayerID == "" || second.PlayerID == "" {
		t.Fatalf("expected non-empty player IDs, got %q and %q", first.PlayerID, second.PlayerID)
	}
	if first.PlayerID == second.PlayerID {
		t.Fatalf("expected distinct player IDs, got same value %q", first.PlayerID)
	}
}

func TestRoomJoin_ThirdPlayerGetsError(t *testing.T) {
	server := newTestServer(t)

	conn1 := dial(t, server)
	conn2 := dial(t, server)
	conn3 := dial(t, server)

	joinRoom(t, conn1, "ABCD")
	joinRoom(t, conn2, "ABCD")
	third := joinRoom(t, conn3, "ABCD")

	if third.Type != "error" {
		t.Fatalf("expected error for third player, got %s", third.Type)
	}
	if third.Message != "room is full" {
		t.Fatalf("expected 'room is full' message, got %q", third.Message)
	}
}

func TestRoomJoin_CannotJoinSameRoomTwice(t *testing.T) {
	server := newTestServer(t)
	conn := dial(t, server)

	first := joinRoom(t, conn, "ABCD")
	if first.Type != "room.joined" {
		t.Fatalf("expected room.joined, got %s (%s)", first.Type, first.Message)
	}

	second := joinRoom(t, conn, "ABCD")
	if second.Type != "error" {
		t.Fatalf("expected error on second join, got %s", second.Type)
	}
}

func TestRoomJoin_RoleFreedAfterMissionControlDisconnects(t *testing.T) {
	server := newTestServer(t)

	conn1 := dial(t, server)
	conn2 := dial(t, server)

	first := joinRoom(t, conn1, "ABCD")
	if first.Role != "mission_control" {
		t.Fatalf("expected first player role mission_control, got %s", first.Role)
	}
	joinRoom(t, conn2, "ABCD")

	if err := conn1.Close(websocket.StatusNormalClosure, ""); err != nil {
		t.Fatalf("close conn1 failed: %v", err)
	}

	// Give the hub time to process the disconnect.
	time.Sleep(100 * time.Millisecond)

	conn3 := dial(t, server)
	third := joinRoom(t, conn3, "ABCD")

	if third.Type != "room.joined" {
		t.Fatalf("expected room.joined, got %s (%s)", third.Type, third.Message)
	}
	if third.Role != "mission_control" {
		t.Fatalf("expected freed role mission_control, got %s", third.Role)
	}
}

func TestRoomJoin_RoomDeletedAfterLastPlayerLeaves(t *testing.T) {
	server := newTestServer(t)

	conn1 := dial(t, server)

	first := joinRoom(t, conn1, "ABCD")
	if first.Role != "mission_control" {
		t.Fatalf("expected first player role mission_control, got %s", first.Role)
	}

	if err := conn1.Close(websocket.StatusNormalClosure, ""); err != nil {
		t.Fatalf("close conn1 failed: %v", err)
	}

	// Give the hub time to process the disconnect and delete the room.
	time.Sleep(100 * time.Millisecond)

	conn2 := dial(t, server)
	second := joinRoom(t, conn2, "ABCD")

	if second.Type != "room.joined" {
		t.Fatalf("expected room.joined, got %s (%s)", second.Type, second.Message)
	}
	if second.Role != "mission_control" {
		t.Fatalf("expected recreated room to assign mission_control, got %s", second.Role)
	}
}

func TestMain(m *testing.M) {
	os.Setenv("JWT_SECRET", "test_secret_for_hub_tests")
	os.Exit(m.Run())
}
