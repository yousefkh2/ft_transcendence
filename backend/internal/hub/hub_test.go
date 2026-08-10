package hub_test

import (
	"context"
	"fmt"
	"net/http/httptest"
	"os"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/coder/websocket"
	"github.com/coder/websocket/wsjson"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"

	"transcendence/backend/internal/db"
	"transcendence/backend/internal/hub"
	"transcendence/backend/internal/model"
)

// fakeRoleResolver simuliert die Rollen, die in echt schon beim REST-Lobby-Join
// in der DB festgelegt wurden -- ohne dafür eine echte Postgres-Instanz zu brauchen.
type fakeRoleResolver struct {
	mu    sync.Mutex
	roles map[string]map[string]string // roomCode -> userID -> role
}

func newFakeRoleResolver() *fakeRoleResolver {
	return &fakeRoleResolver{roles: make(map[string]map[string]string)}
}

func (f *fakeRoleResolver) assign(roomCode, userID, role string) {
	f.mu.Lock()
	defer f.mu.Unlock()
	if f.roles[roomCode] == nil {
		f.roles[roomCode] = make(map[string]string)
	}
	f.roles[roomCode][userID] = role
}

func (f *fakeRoleResolver) GetParticipantRole(ctx context.Context, code, userID string) (string, string, error) {
	f.mu.Lock()
	defer f.mu.Unlock()

	participants, ok := f.roles[code]
	if !ok {
		return "", "", db.ErrNotInLobby
	}
	role, ok := participants[userID]
	if !ok {
		return "", "", db.ErrNotInLobby
	}
	return "session-" + code, role, nil
}

func newTestServer(t *testing.T, roles hub.RoleResolver) *httptest.Server {
	t.Helper()

	h := hub.NewHubWithOrigins([]string{"*"}, roles)

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

func tokenFor(t *testing.T, userID string) string {
	t.Helper()

	claims := jwt.MapClaims{
		"sub": userID,
		"exp": time.Now().Add(time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signed, err := token.SignedString([]byte("test_secret_for_hub_tests"))
	if err != nil {
		t.Fatalf("failed to sign test token: %v", err)
	}

	return signed
}

func joinRoom(t *testing.T, conn *websocket.Conn, roomCode, userID string) model.ServerMessage {
	t.Helper()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := wsjson.Write(ctx, conn, model.ClientMessage{
		Type:     "room.join",
		RoomCode: roomCode,
		Token:    tokenFor(t, userID),
	}); err != nil {
		t.Fatalf("write room.join failed: %v", err)
	}

	var response model.ServerMessage
	if err := wsjson.Read(ctx, conn, &response); err != nil {
		t.Fatalf("read room.join response failed: %v", err)
	}

	return response
}

func uniqueUserID(t *testing.T, label string) string {
	t.Helper()
	return fmt.Sprintf("%s-%d", label, time.Now().UnixNano())
}

func TestRoomJoin_RoleComesFromResolverNotConnectionOrder(t *testing.T) {
	roles := newFakeRoleResolver()
	alice := uniqueUserID(t, "alice")
	bob := uniqueUserID(t, "bob")
	roles.assign("ABCD", alice, "mission_control")
	roles.assign("ABCD", bob, "on_site")

	server := newTestServer(t, roles)

	// Bob verbindet sich zuerst -- im alten Modell hätte er fälschlich
	// mission_control bekommen. Jetzt bekommt er trotzdem on_site, weil
	// das schon vorher in der Lobby so festgelegt wurde.
	bobConn := dial(t, server)
	bobResp := joinRoom(t, bobConn, "ABCD", bob)
	if bobResp.Role != "on_site" {
		t.Fatalf("expected bob to get on_site regardless of connection order, got %s", bobResp.Role)
	}

	aliceConn := dial(t, server)
	aliceResp := joinRoom(t, aliceConn, "ABCD", alice)
	if aliceResp.Role != "mission_control" {
		t.Fatalf("expected alice to get mission_control, got %s", aliceResp.Role)
	}
}

func TestRoomJoin_UnknownUserGetsError(t *testing.T) {
	roles := newFakeRoleResolver()
	alice := uniqueUserID(t, "alice")
	roles.assign("ABCD", alice, "mission_control")

	server := newTestServer(t, roles)
	conn := dial(t, server)

	stranger := uniqueUserID(t, "stranger")
	resp := joinRoom(t, conn, "ABCD", stranger)

	if resp.Type != "error" {
		t.Fatalf("expected error for user not in lobby, got %s", resp.Type)
	}
}

func TestRoomJoin_DistinctPlayerIDs(t *testing.T) {
	roles := newFakeRoleResolver()
	alice := uniqueUserID(t, "alice")
	bob := uniqueUserID(t, "bob")
	roles.assign("ABCD", alice, "mission_control")
	roles.assign("ABCD", bob, "on_site")

	server := newTestServer(t, roles)

	conn1 := dial(t, server)
	conn2 := dial(t, server)

	first := joinRoom(t, conn1, "ABCD", alice)
	second := joinRoom(t, conn2, "ABCD", bob)

	if first.PlayerID == "" || second.PlayerID == "" {
		t.Fatalf("expected non-empty player IDs, got %q and %q", first.PlayerID, second.PlayerID)
	}
	if first.PlayerID == second.PlayerID {
		t.Fatalf("expected distinct player IDs, got same value %q", first.PlayerID)
	}
}

func TestRoomJoin_SameRoleTwiceGetsError(t *testing.T) {
	roles := newFakeRoleResolver()
	alice := uniqueUserID(t, "alice")
	roles.assign("ABCD", alice, "mission_control")

	server := newTestServer(t, roles)
	conn := dial(t, server)

	first := joinRoom(t, conn, "ABCD", alice)
	if first.Type != "room.joined" {
		t.Fatalf("expected room.joined, got %s (%s)", first.Type, first.Message)
	}

	conn2 := dial(t, server)
	second := joinRoom(t, conn2, "ABCD", alice)
	if second.Type != "error" {
		t.Fatalf("expected error when the same role tries to connect twice, got %s", second.Type)
	}
}

func TestRoomJoin_CannotJoinSameRoomTwiceOnSameConnection(t *testing.T) {
	roles := newFakeRoleResolver()
	alice := uniqueUserID(t, "alice")
	roles.assign("ABCD", alice, "mission_control")

	server := newTestServer(t, roles)
	conn := dial(t, server)

	first := joinRoom(t, conn, "ABCD", alice)
	if first.Type != "room.joined" {
		t.Fatalf("expected room.joined, got %s (%s)", first.Type, first.Message)
	}

	second := joinRoom(t, conn, "ABCD", alice)
	if second.Type != "error" {
		t.Fatalf("expected error on second join, got %s", second.Type)
	}
}

func TestRoomJoin_RoleFreedAfterDisconnect(t *testing.T) {
	roles := newFakeRoleResolver()
	alice := uniqueUserID(t, "alice")
	bob := uniqueUserID(t, "bob")
	roles.assign("ABCD", alice, "mission_control")
	roles.assign("ABCD", bob, "on_site")

	server := newTestServer(t, roles)

	conn1 := dial(t, server)
	conn2 := dial(t, server)

	first := joinRoom(t, conn1, "ABCD", alice)
	if first.Role != "mission_control" {
		t.Fatalf("expected alice to get mission_control, got %s", first.Role)
	}
	joinRoom(t, conn2, "ABCD", bob)

	if err := conn1.Close(websocket.StatusNormalClosure, ""); err != nil {
		t.Fatalf("close conn1 failed: %v", err)
	}

	time.Sleep(100 * time.Millisecond)

	conn3 := dial(t, server)
	third := joinRoom(t, conn3, "ABCD", alice)

	if third.Type != "room.joined" {
		t.Fatalf("expected room.joined, got %s (%s)", third.Type, third.Message)
	}
	if third.Role != "mission_control" {
		t.Fatalf("expected freed role mission_control, got %s", third.Role)
	}
}

func TestRoomJoin_RoomDeletedAfterLastPlayerLeaves(t *testing.T) {
	roles := newFakeRoleResolver()
	alice := uniqueUserID(t, "alice")
	roles.assign("ABCD", alice, "mission_control")

	server := newTestServer(t, roles)

	conn1 := dial(t, server)

	first := joinRoom(t, conn1, "ABCD", alice)
	if first.Role != "mission_control" {
		t.Fatalf("expected alice to get mission_control, got %s", first.Role)
	}

	if err := conn1.Close(websocket.StatusNormalClosure, ""); err != nil {
		t.Fatalf("close conn1 failed: %v", err)
	}

	time.Sleep(100 * time.Millisecond)

	conn2 := dial(t, server)
	second := joinRoom(t, conn2, "ABCD", alice)

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