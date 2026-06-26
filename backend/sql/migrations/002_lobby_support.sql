ALTER TABLE game_sessions
    ADD COLUMN IF NOT EXISTS code VARCHAR(8) UNIQUE,
    ADD COLUMN IF NOT EXISTS host_user_id UUID REFERENCES users(id) ON DELETE CASCADE;

CREATE INDEX IF NOT EXISTS idx_game_sessions_code ON game_sessions(code);
CREATE INDEX IF NOT EXISTS idx_game_sessions_status ON game_sessions(status);

-- Rolle wird erst beim tatsächlichen Match-Start vom Game-Engine-Code zugewiesen
-- (mission_control/on_site), in der Lobby-Phase ist noch keine Rolle bekannt.
ALTER TABLE session_participants
    ALTER COLUMN role DROP NOT NULL;