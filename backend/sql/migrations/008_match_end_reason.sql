-- How a finished round ended: 'solved', 'time_up' or 'out_of_checks'. NULL for rounds saved before this column.
ALTER TABLE game_sessions
    ADD COLUMN IF NOT EXISTS end_reason VARCHAR(20);
