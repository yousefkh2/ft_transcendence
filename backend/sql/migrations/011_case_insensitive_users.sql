-- Usernames and emails are unique regardless of case: "Anna" and "anna" are the same user, and so
-- are "Anna@Mail.de" and "anna@mail.de". The stored value keeps the case the user typed.
-- The backend looks users up with LOWER(...) too, so these indexes also make those lookups fast.
CREATE UNIQUE INDEX IF NOT EXISTS users_username_lower_key ON users (LOWER(username));
CREATE UNIQUE INDEX IF NOT EXISTS users_email_lower_key ON users (LOWER(email));
