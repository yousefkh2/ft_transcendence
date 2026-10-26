-- The language a user wants to practise. Uses the same codes as game_sessions.lang (en, de, pl, tr).
ALTER TABLE users
    ADD COLUMN IF NOT EXISTS learning_lang VARCHAR(5) NOT NULL DEFAULT 'de';
