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

func TestRoomLeave_NotifiesRemainingPlayer(t *testing.T) {
	roles := newFakeRoleResolver()
	alice := uniqueUserID(t, "alice")
	bob := uniqueUserID(t, "bob")
	roles.assign("ABCD", alice, "mission_control")
	roles.assign("ABCD", bob, "on_site")

	server := newTestServer(t, roles)
	aliceConn := dial(t, server)
	bobConn := dial(t, server)
	joinRoom(t, aliceConn, "ABCD", alice)
	joinRoom(t, bobConn, "ABCD", bob)

	if err := aliceConn.Close(websocket.StatusNormalClosure, ""); err != nil {
		t.Fatalf("close alice connection failed: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var response model.ServerMessage
	if err := wsjson.Read(ctx, bobConn, &response); err != nil {
		t.Fatalf("read disconnect notification failed: %v", err)
	}
	if response.Type != "player.disconnected" {
		t.Fatalf("expected player.disconnected, got %s", response.Type)
	}
	if response.Message != "A player left the game" {
		t.Fatalf("expected leave message, got %q", response.Message)
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

func TestGameRelay_ForwardsPayloadToOtherRole(t *testing.T) {
	roles := newFakeRoleResolver()
	alice := uniqueUserID(t, "alice")
	bob := uniqueUserID(t, "bob")
	roles.assign("RELAY", alice, "mission_control")
	roles.assign("RELAY", bob, "on_site")

	server := newTestServer(t, roles)
	aliceConn := dial(t, server)
	bobConn := dial(t, server)
	joinRoom(t, aliceConn, "RELAY", alice)
	joinRoom(t, bobConn, "RELAY", bob)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// on_site -> mission_control
	layout := []byte(`{"layout":["","lamp","","plant","couch","","table",""]}`)
	if err := wsjson.Write(ctx, bobConn, model.ClientMessage{Type: "game.layout", Payload: layout}); err != nil {
		t.Fatalf("write game.layout failed: %v", err)
	}
	var got model.ServerMessage
	if err := wsjson.Read(ctx, aliceConn, &got); err != nil {
		t.Fatalf("read relayed layout failed: %v", err)
	}
	if got.Type != "game.layout" || got.FromRole != "on_site" || string(got.Payload) != string(layout) {
		t.Fatalf("unexpected relay: type=%s from=%s payload=%s", got.Type, got.FromRole, got.Payload)
	}

	// mission_control -> on_site
	result := []byte(`{"solved":false,"correct":2}`)
	if err := wsjson.Write(ctx, aliceConn, model.ClientMessage{Type: "game.check_result", Payload: result}); err != nil {
		t.Fatalf("write game.check_result failed: %v", err)
	}
	if err := wsjson.Read(ctx, bobConn, &got); err != nil {
		t.Fatalf("read relayed check result failed: %v", err)
	}
	if got.Type != "game.check_result" || got.FromRole != "mission_control" || string(got.Payload) != string(result) {
		t.Fatalf("unexpected relay: type=%s from=%s payload=%s", got.Type, got.FromRole, got.Payload)
	}

	// mission_control may not send a layout
	if err := wsjson.Write(ctx, aliceConn, model.ClientMessage{Type: "game.layout", Payload: layout}); err != nil {
		t.Fatalf("write game.layout failed: %v", err)
	}
	if err := wsjson.Read(ctx, aliceConn, &got); err != nil {
		t.Fatalf("read error response failed: %v", err)
	}
	if got.Type != "error" {
		t.Fatalf("expected error for wrong sender, got %s", got.Type)
	}
}

// savedMatches records what the hub would write to the match history.
type savedMatches struct {
	mu      sync.Mutex
	reasons []string
}

func (s *savedMatches) save(code string, startedAt, endedAt time.Time, endReason string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.reasons = append(s.reasons, endReason)
	return nil
}

func (s *savedMatches) waitFor(t *testing.T, count int) []string {
	t.Helper()
	deadline := time.Now().Add(2 * time.Second)
	for time.Now().Before(deadline) {
		s.mu.Lock()
		n := len(s.reasons)
		s.mu.Unlock()
		if n >= count {
			break
		}
		time.Sleep(10 * time.Millisecond)
	}
	// Give a wrong second save a moment to show up.
	time.Sleep(50 * time.Millisecond)
	s.mu.Lock()
	defer s.mu.Unlock()
	return append([]string(nil), s.reasons...)
}

// startMatch joins mission control and on site to the room and returns the mission control connection.
func startMatch(t *testing.T, code string) (*websocket.Conn, *savedMatches) {
	t.Helper()
	roles := newFakeRoleResolver()
	alice := uniqueUserID(t, "alice")
	bob := uniqueUserID(t, "bob")
	roles.assign(code, alice, "mission_control")
	roles.assign(code, bob, "on_site")

	saved := &savedMatches{}
	h := hub.NewHubWithOrigins([]string{"*"}, roles)
	h.SaveMatch = saved.save
	e := echo.New()
	e.GET("/ws", h.HandleWebSocket)
	server := httptest.NewServer(e)
	t.Cleanup(server.Close)

	aliceConn := dial(t, server)
	bobConn := dial(t, server)
	joinRoom(t, aliceConn, code, alice)
	joinRoom(t, bobConn, code, bob)
	return aliceConn, saved
}

func sendCheck(t *testing.T, conn *websocket.Conn, payload string) {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := wsjson.Write(ctx, conn, model.ClientMessage{Type: "game.check_result", Payload: []byte(payload)}); err != nil {
		t.Fatalf("write game.check_result failed: %v", err)
	}
}

func TestMatchHistory_SolvedCheckSavesWinOnce(t *testing.T) {
	missionControl, saved := startMatch(t, "WINNER")
	sendCheck(t, missionControl, `{"solved":false,"correct":3,"attemptsLeft":4}`)
	sendCheck(t, missionControl, `{"solved":true,"correct":8,"attemptsLeft":4}`)
	sendCheck(t, missionControl, `{"solved":true,"correct":8,"attemptsLeft":4}`)

	reasons := saved.waitFor(t, 1)
	if len(reasons) != 1 || reasons[0] != db.EndSolved {
		t.Fatalf("expected exactly one solved match, got %v", reasons)
	}
}

func TestMatchHistory_OutOfChecksSavesLoss(t *testing.T) {
	missionControl, saved := startMatch(t, "LOSER")
	// The server counts wrong checks itself, whatever attemptsLeft the client claims.
	for i := 0; i < 5; i++ {
		sendCheck(t, missionControl, `{"solved":false,"correct":2,"attemptsLeft":99}`)
	}

	reasons := saved.waitFor(t, 1)
	if len(reasons) != 1 || reasons[0] != db.EndOutOfChecks {
		t.Fatalf("expected exactly one out-of-checks match, got %v", reasons)
	}
}

func TestMatchHistory_FewWrongChecksSaveNothing(t *testing.T) {
	missionControl, saved := startMatch(t, "PENDING")
	for i := 0; i < 4; i++ {
		sendCheck(t, missionControl, `{"solved":false,"correct":2}`)
	}

	if reasons := saved.waitFor(t, 1); len(reasons) != 0 {
		t.Fatalf("expected no saved match yet, got %v", reasons)
	}
}

func TestMain(m *testing.M) {
	os.Setenv("JWT_SECRET", "test_secret_for_hub_tests")
	os.Exit(m.Run())
}
