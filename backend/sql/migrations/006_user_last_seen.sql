-- When the user's browser last reported that the app is open (heartbeat). Used for the online status.
ALTER TABLE users
    ADD COLUMN IF NOT EXISTS last_seen_at TIMESTAMPTZ;
