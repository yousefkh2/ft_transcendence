-- Profile pictures, stored in the database so no extra volume is needed. One picture per user,
-- already re-encoded by the backend (small, no metadata). Deleted together with the account.
CREATE TABLE IF NOT EXISTS user_avatars (
    user_id UUID PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
    image BYTEA NOT NULL,
    content_type VARCHAR(20) NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
