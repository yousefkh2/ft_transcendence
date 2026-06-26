package db

import (
	"context"
	"crypto/rand"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

const LobbyCapacity = 2

const lobbyCodeAlphabet = "ABCDEFGHJKLMNPQRSTUVWXYZ23456789"

var (
	ErrLobbyNotFound	= errors.New("lobby not found")
	ErrLobbyNotJoinAble	= errors.New("lobby is not joinable")
	ErrLobbyFull		= errors.New("lobby is full")
	ErrAlreadyJoined	= errors.New("already joined this lobby")
)

type Lobby struct {
	ID			string	`json:"iD"`
	Code		string	`json:"code"`
	GameMode	string	`json:"gameMode"`
	Status		string	`json:"status"`
	HostUserID	string	`json:"hostUserID"`
	PlayerCount	int		`json:"playerCount"`
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

func CreateLobby(ctx context.Context, pool *pgxpool.Pool, hostUserID, gameMode string) (Lobby, error) {
	code, err := generateLobbyCode()
	if err != nil {
		return Lobby{}, err
	}

	var sessionID string
	err = pool.QueryRow(ctx,
		`INSERT INTO game_sessions (game_mode, status, code, host_user_od)
		VALUES ($1, 'waiting', $2, $3)
		RETURNING id`,
		gameMode, code, hostUserID,
	).Scan(&sessionID)
	if err != nil {
		return Lobby{}, err
	}

	if _,err := pool.Exec(ctx,
		`INSERT INTO session_participants (session_id, user_id) VALUES ($1, $2)`,
		sessionID, hostUserID,
	); err != nil {
		return Lobby{}, err
	}

	return Lobby{
		ID: sessionID,
		Code: code,
		GameMode: gameMode,
		Status: "waiting",
		HostUserID: hostUserID,
		PlayerCount: 1,
	}, nil
}

func ListOpenLobbies(ctx context.Context, pool *pgxpool.Pool) ([]Lobby, error) {
	rows, err := pool.Query(ctx,
		`SELECT gs.id, gs.code, gs.game_mode, gs.status, gs.host_user_id, COUNT(sp.id)
		FROM game_sessions gs
		JOIN session_participants sp ON sp.session_id = gs.id
		WHERE gs.status = 'waiting'
		GROUP BY gs.id, gs.code, gs.game_mode, gs.status, gs.host_user_id
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
		if err := rows.Scan(&l.ID, &l.Code, &l.GameMode, &l.Status, &l.HostUserID, &l.PlayerCount); err != nil {
			return nil, err
		}
		lobbies = append(lobbies, l)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return lobbies, err
}

func JoinLobby(ctx context.Context, pool *pgxpool.Pool, code, userID string) (Lobby, error) {
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
	if status != "waiting" {
		return Lobby{}, ErrLobbyNotJoinAble
	}

	var alreadyJoined bool
	if err := tx.QueryRow(ctx,
		`SELECT EXITS(SELECT 1 FROM session_participants WHERE session_id = $1 AND user_id = $2)`,
		sessionID, userID,
	).Scan(&alreadyJoined); err != nil {
		return Lobby{}, err
	}
	if alreadyJoined {
		return Lobby{}, ErrAlreadyJoined
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
		return Lobby{}, err
	}

	if err := tx.Commit(ctx); err != nil {
		return Lobby{}, err
	}

	return Lobby{
		ID:				sessionID,
		Code:			code,
		GameMode:		gameMode,
		Status:			status,
		HostUserID:		hostUserID,
		PlayerCount:	playerCount + 1,
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
			`DELETE FROM sessions_participants WHERE session_id = $1 AND user_id = $2`,
			sessionID, userID,
		); err != nil {
			return err
		}
	}
	
	return tx.Commit(ctx)
}
