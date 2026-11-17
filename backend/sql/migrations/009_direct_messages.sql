-- Direct chat messages between two friends. Deleted together with either account.
-- read_at stays NULL until the recipient opens the conversation (used for the unread counter).
CREATE TABLE IF NOT EXISTS direct_messages (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    sender_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    recipient_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    body TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    read_at TIMESTAMPTZ,
    CONSTRAINT check_message_not_self CHECK (sender_id != recipient_id),
    CONSTRAINT check_message_length CHECK (char_length(body) BETWEEN 1 AND 1000)
);

-- One conversation is all messages between the same two users, in either direction.
CREATE INDEX IF NOT EXISTS idx_direct_messages_conversation ON direct_messages (
    LEAST(sender_id, recipient_id),
    GREATEST(sender_id, recipient_id),
    created_at DESC
);

CREATE INDEX IF NOT EXISTS idx_direct_messages_unread ON direct_messages (recipient_id) WHERE read_at IS NULL;
