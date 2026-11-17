-- How many wrong checks a finished round had (for the "flawless" achievement). NULL for older rounds.
ALTER TABLE game_sessions
    ADD COLUMN IF NOT EXISTS failed_checks INTEGER;

-- Achievements a user has unlocked. The rules live in the backend (internal/game/achievements.go);
-- achievement_id is one of their stable ids. seen_at stays NULL until the user was shown the
-- "achievement unlocked" message.
CREATE TABLE IF NOT EXISTS user_achievements (
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    achievement_id VARCHAR(40) NOT NULL,
    unlocked_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    seen_at TIMESTAMPTZ,
    PRIMARY KEY (user_id, achievement_id)
);

-- The leaderboard sorts all users by XP.
CREATE INDEX IF NOT EXISTS idx_users_xp ON users (xp DESC);
