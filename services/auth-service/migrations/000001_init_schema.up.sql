-- +migrate Up
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TABLE users
(
    id                         UUID PRIMARY KEY,
    username                   VARCHAR(255) UNIQUE NOT NULL,
    email                      VARCHAR(255) UNIQUE NOT NULL,
    password_hash              VARCHAR(255)        NOT NULL,
    consent_to_data_processing BOOLEAN             NOT NULL DEFAULT FALSE,
    created_at                 TIMESTAMP                    DEFAULT CURRENT_TIMESTAMP,
    updated_at                 TIMESTAMP                    DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE refresh_tokens
(
    id         SERIAL PRIMARY KEY,
    user_id    UUID REFERENCES users (id) ON DELETE CASCADE,
    token      VARCHAR(512) NOT NULL,
    expires_at TIMESTAMP    NOT NULL,
    revoked    BOOLEAN   DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_users_email ON users (email);
CREATE INDEX idx_refresh_tokens_user_id ON refresh_tokens (user_id);
