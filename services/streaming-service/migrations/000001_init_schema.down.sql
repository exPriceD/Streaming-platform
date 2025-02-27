-- +migrate Down
DROP TABLE IF EXISTS streams;
DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS user_profile;
DROP TABLE IF EXISTS stream_keys;
DROP TABLE IF EXISTS stream_history;
