-- +migrate Up
CREATE UNIQUE INDEX IF NOT EXISTS USERS_EMAIL ON USERS(EMAIL);
CREATE UNIQUE INDEX IF NOT EXISTS USERS_NAME on USERS(NAME);

-- +migrate Down
DROP INDEX IF EXISTS USERS_EMAIL;
DROP INDEX IF EXISTS USERS_NAME;