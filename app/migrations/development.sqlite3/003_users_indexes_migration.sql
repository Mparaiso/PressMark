-- +migrate Up
CREATE UNIQUE INDEX IF NOT EXISTS USERS_EMAIL ON users(email);
CREATE UNIQUE INDEX IF NOT EXISTS USERS_NAME on users(name);

-- +migrate Down
DROP INDEX IF EXISTS USERS_EMAIL;
DROP INDEX IF EXISTS USERS_NAME;