package db

import (
	"context"
	"crypto/rand"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type LobbyPlayer struct {
	UserID   string
	Username string
	Role     string
}

const LobbyCapacity = 2

const lobbyCodeAlphabet = "ABCDEFGHJKLMNPQRSTUVWXYZ23456789"

const (
	RoleMissionControl = "mission_control"
	RoleOnSite         = "on_site"
)

const DefaultLobbyLang = "en"

var SupportedLobbyLangs = map[string]bool{
	"en": true,
	"de": true,
	"tr": true,
	"pl": true,
}

var (
	ErrLobbyNotFound         = errors.New("lobby not found")
	ErrLobbyNotJoinable      = errors.New("lobby is not joinable")
	ErrLobbyFull             = errors.New("lobby is full")
	ErrAlreadyJoined         = errors.New("already joined this lobby")
	ErrNotInLobby            = errors.New("user is not part of this lobby")
	ErrLobbyNotStarted       = errors.New("lobby has not started yet")
	ErrNotHost               = errors.New("only the host can do this")
	ErrUnsupportedLang       = errors.New("unsupported language")
	ErrLobbyLanguageMismatch = errors.New("lobby language does not match your current language")
	ErrInvalidRole           = errors.New("invalid role")
	ErrRoleTaken             = errors.New("role is already taken")
)

type Lobby struct {
	ID          string `json:"id"`
	Code        string `json:"code"`
	GameMode    string `json:"gameMode"`
	Status      string `json:"status"`
	HostUserID  string `json:"hostUserID"`
	PlayerCount int    `json:"playerCount"`
	CurrLang    string `json:"currLang"`
}

func generateLobbyCode() (string, error) {
	raw := make([]byte, 4)
	if _, err := rand.Read(raw); err != nil {
		return "", err
	}
	code := make([]byte, 4)
	for i, v := range raw {
		code[i] = lobbyCodeAlphabet[int(v)%len(lobbyCodeAlphabet)]
	}

	return string(code), nil
}

func CreateLobby(ctx context.Context, pool *pgxpool.Pool, hostUserID, gameMode, lang string) (Lobby, error) {
	if lang == "" {
		lang = DefaultLobbyLang
	}
	if !SupportedLobbyLangs[lang] {
		return Lobby{}, ErrUnsupportedLang
	}

	code, err := generateLobbyCode()
	if err != nil {
		return Lobby{}, err
	}

	tx, err := pool.Begin(ctx)
	if err != nil {
		return Lobby{}, err
	}
	defer tx.Rollback(ctx)

	var sessionID string
	err = tx.QueryRow(ctx,
		`INSERT INTO game_sessions (game_mode, status, code, host_user_id, lang)
		VALUES ($1, 'waiting', $2, $3, $4)
		RETURNING id`,
		gameMode, code, hostUserID, lang,
	).Scan(&sessionID)
	if err != nil {
		return Lobby{}, err
	}

	if _, err := tx.Exec(ctx,
		`INSERT INTO session_participants (session_id, user_id) VALUES ($1, $2)`,
		sessionID, hostUserID,
	); err != nil {
		return Lobby{}, err
	}

	if err := tx.Commit(ctx); err != nil {
		return Lobby{}, err
	}

	return Lobby{
		ID:          sessionID,
		Code:        code,
		GameMode:    gameMode,
		Status:      "waiting",
		HostUserID:  hostUserID,
		PlayerCount: 1,
		CurrLang:    lang,
	}, nil
}

func ListOpenLobbies(ctx context.Context, pool *pgxpool.Pool) ([]Lobby, error) {
	rows, err := pool.Query(ctx,
		`SELECT gs.id, gs.code, gs.game_mode, gs.status, gs.host_user_id, gs.lang, COUNT(sp.id)
		FROM game_sessions gs
		JOIN session_participants sp ON sp.session_id = gs.id
		WHERE gs.status = 'waiting'
		GROUP BY gs.id, gs.code, gs.game_mode, gs.status, gs.host_user_id, gs.lang
		HAVING COUNT(sp.id) < $1
		ORDER BY gs.started_at DESC`,
		LobbyCapacity,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	lobbies := make([]Lobby, 0)
	for rows.Next() {
		var l Lobby
		if err := rows.Scan(&l.ID, &l.Code, &l.GameMode, &l.Status, &l.HostUserID, &l.CurrLang, &l.PlayerCount); err != nil {
			return nil, err
		}
		lobbies = append(lobbies, l)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return lobbies, nil
}

func ListLobbyPlayers(ctx context.Context, pool *pgxpool.Pool, code string) ([]LobbyPlayer, error) {
	rows, err := pool.Query(ctx,
		`SELECT u.id, u.username, sp.role
		 FROM game_sessions gs
		 JOIN session_participants sp ON sp.session_id = gs.id
		 JOIN users u ON u.id = sp.user_id
		 WHERE gs.code = $1
		 ORDER BY sp.id`,
		code,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	players := make([]LobbyPlayer, 0)
	for rows.Next() {
		var player LobbyPlayer
		var role *string
		if err := rows.Scan(&player.UserID, &player.Username, &role); err != nil {
			return nil, err
		}
		if role != nil {
			player.Role = *role
		}
		players = append(players, player)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return players, nil
}

func JoinLobby(ctx context.Context, pool *pgxpool.Pool, code, userID, requestedLang string) (Lobby, error) {
	tx, err := pool.Begin(ctx)
	if err != nil {
		return Lobby{}, err
	}
	defer tx.Rollback(ctx)

	var sessionID, gameMode, status, hostUserID, lang string
	err = tx.QueryRow(ctx,
		`SELECT id, game_mode, status, host_user_id, lang FROM game_sessions
		WHERE code = $1 FOR UPDATE`,
		code,
	).Scan(&sessionID, &gameMode, &status, &hostUserID, &lang)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return Lobby{}, ErrLobbyNotFound
		}
		return Lobby{}, err
	}
	if status != "waiting" {
		return Lobby{}, ErrLobbyNotJoinable
	}
	if requestedLang != "" {
		if !SupportedLobbyLangs[requestedLang] {
			return Lobby{}, ErrUnsupportedLang
		}
		if requestedLang != lang {
			return Lobby{}, ErrLobbyLanguageMismatch
		}
	}

	var playerCount int
	if err := tx.QueryRow(ctx,
		`SELECT COUNT(*) FROM session_participants WHERE session_id = $1`,
		sessionID,
	).Scan(&playerCount); err != nil {
		return Lobby{}, err
	}
	if playerCount >= LobbyCapacity {
		return Lobby{}, ErrLobbyFull
	}

	if _, err := tx.Exec(ctx,
		`INSERT INTO session_participants (session_id, user_id) VALUES ($1, $2)`,
		sessionID, userID,
	); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return Lobby{}, ErrAlreadyJoined
		}
		return Lobby{}, err
	}

	newPlayerCount := playerCount + 1

	if newPlayerCount == LobbyCapacity {
		if _, err := tx.Exec(ctx,
			`UPDATE game_sessions SET status = 'active' WHERE id = $1`,
			sessionID,
		); err != nil {
			return Lobby{}, err
		}
		status = "active"
	}

	if err := tx.Commit(ctx); err != nil {
		return Lobby{}, err
	}

	return Lobby{
		ID:          sessionID,
		Code:        code,
		GameMode:    gameMode,
		Status:      status,
		HostUserID:  hostUserID,
		PlayerCount: newPlayerCount,
		CurrLang:    lang,
	}, nil
}

func SelectLobbyRole(ctx context.Context, pool *pgxpool.Pool, code, userID, role string) error {
	if role != "" && role != RoleMissionControl && role != RoleOnSite {
		return ErrInvalidRole
	}

	tx, err := pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	var sessionID, participantID string
	if err := tx.QueryRow(ctx,
		`SELECT gs.id, sp.id
		 FROM game_sessions gs
		 JOIN session_participants sp ON sp.session_id = gs.id
		 WHERE gs.code = $1 AND sp.user_id = $2
		 FOR UPDATE OF gs`,
		code, userID,
	).Scan(&sessionID, &participantID); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrNotInLobby
		}
		return err
	}

	if role != "" {
		var roleOwner string
		err = tx.QueryRow(ctx,
			`SELECT user_id FROM session_participants
			 WHERE session_id = $1 AND role = $2 AND user_id <> $3`,
			sessionID, role, userID,
		).Scan(&roleOwner)
		if err == nil {
			return ErrRoleTaken
		}
		if !errors.Is(err, pgx.ErrNoRows) {
			return err
		}
	}

	if role == "" {
		_, err = tx.Exec(ctx, `UPDATE session_participants SET role = NULL WHERE id = $1`, participantID)
	} else {
		_, err = tx.Exec(ctx, `UPDATE session_participants SET role = $1 WHERE id = $2`, role, participantID)
	}
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func UpdateLobbyLanguage(ctx context.Context, pool *pgxpool.Pool, code, userID, lang string) (Lobby, error) {
	if !SupportedLobbyLangs[lang] {
		return Lobby{}, ErrUnsupportedLang
	}

	tx, err := pool.Begin(ctx)
	if err != nil {
		return Lobby{}, err
	}
	defer tx.Rollback(ctx)

	var sessionID, gameMode, status, hostUserID string
	err = tx.QueryRow(ctx,
		`SELECT id, game_mode, status, host_user_id FROM game_sessions
		WHERE code = $1 FOR UPDATE`,
		code,
	).Scan(&sessionID, &gameMode, &status, &hostUserID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return Lobby{}, ErrLobbyNotFound
		}
		return Lobby{}, err
	}
	if hostUserID != userID {
		return Lobby{}, ErrNotHost
	}

	if _, err := tx.Exec(ctx,
		`UPDATE game_sessions SET lang = $1 WHERE id = $2`,
		lang, sessionID,
	); err != nil {
		return Lobby{}, err
	}

	var playerCount int
	if err := tx.QueryRow(ctx,
		`SELECT COUNT(*) FROM session_participants WHERE session_id = $1`,
		sessionID,
	).Scan(&playerCount); err != nil {
		return Lobby{}, err
	}

	if err := tx.Commit(ctx); err != nil {
		return Lobby{}, err
	}

	return Lobby{
		ID:          sessionID,
		Code:        code,
		GameMode:    gameMode,
		Status:      status,
		HostUserID:  hostUserID,
		PlayerCount: playerCount,
		CurrLang:    lang,
	}, nil
}

func LeaveLobby(ctx context.Context, pool *pgxpool.Pool, code, userID string) error {
	tx, err := pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	var sessionID, hostUserID string
	err = tx.QueryRow(ctx,
		`SELECT id, host_user_id FROM game_sessions WHERE code = $1 FOR UPDATE`,
		code,
	).Scan(&sessionID, &hostUserID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrLobbyNotFound
		}
		return err
	}

	if userID == hostUserID {
		if _, err := tx.Exec(ctx, `DELETE FROM game_sessions WHERE id = $1`, sessionID); err != nil {
			return err
		}
	} else {
		if _, err := tx.Exec(ctx,
			`DELETE FROM session_participants WHERE session_id = $1 AND user_id = $2`,
			sessionID, userID,
		); err != nil {
			return err
		}
	}

	return tx.Commit(ctx)
}

func RemoveDisconnectedParticipant(ctx context.Context, pool *pgxpool.Pool, code, userID string) error {
	tx, err := pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	var sessionID, hostUserID string
	if err := tx.QueryRow(ctx,
		`SELECT id, host_user_id FROM game_sessions WHERE code = $1 FOR UPDATE`,
		code,
	).Scan(&sessionID, &hostUserID); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrLobbyNotFound
		}
		return err
	}

	result, err := tx.Exec(ctx,
		`DELETE FROM session_participants WHERE session_id = $1 AND user_id = $2`,
		sessionID, userID,
	)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return ErrNotInLobby
	}

	var playerCount int
	if err := tx.QueryRow(ctx,
		`SELECT COUNT(*) FROM session_participants WHERE session_id = $1`,
		sessionID,
	).Scan(&playerCount); err != nil {
		return err
	}

	if playerCount == 0 {
		if _, err := tx.Exec(ctx, `DELETE FROM game_sessions WHERE id = $1`, sessionID); err != nil {
			return err
		}
	} else {
		if _, err := tx.Exec(ctx,
			`UPDATE game_sessions SET status = 'waiting' WHERE id = $1`,
			sessionID,
		); err != nil {
			return err
		}

		if hostUserID == userID {
			if _, err := tx.Exec(ctx,
				`UPDATE game_sessions
				 SET host_user_id = (SELECT user_id FROM session_participants WHERE session_id = $1 LIMIT 1)
				 WHERE id = $1`,
				sessionID,
			); err != nil {
				return err
			}
		}
	}

	return tx.Commit(ctx)
}

func GetParticipantRole(ctx context.Context, pool *pgxpool.Pool, code, userID string) (sessionID, role string, err error) {
	var roleNullable *string
	err = pool.QueryRow(ctx,
		`SELECT gs.id, sp.role
		FROM game_sessions gs
		JOIN session_participants sp ON sp.session_id = gs.id
		WHERE gs.code = $1 AND sp.user_id = $2`,
		code, userID,
	).Scan(&sessionID, &roleNullable)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", "", ErrNotInLobby
		}
		return "", "", err
	}
	if roleNullable == nil {
		return "", "", ErrLobbyNotStarted
	}
	return sessionID, *roleNullable, nil
}
