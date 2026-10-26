-- Declined friend requests used to stay in the table. Because of idx_unique_friendship they blocked
-- every new request between the two users forever. Declining now deletes the request instead,
-- so old "declined" rows are removed here. Safe to run again: it only deletes rows that should not exist.
DELETE FROM friendships WHERE status = 'declined';
