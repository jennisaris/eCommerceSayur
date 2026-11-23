-- Server-side refresh sessions for HttpOnly cookie authentication.
CREATE TABLE IF NOT EXISTS auth_sessions (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token_hash VARCHAR(64) NOT NULL UNIQUE,
    expires_at TIMESTAMP NOT NULL,
    revoked_at TIMESTAMP NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_auth_sessions_user_id
    ON auth_sessions(user_id);

CREATE INDEX IF NOT EXISTS idx_auth_sessions_active_token
    ON auth_sessions(token_hash, expires_at)
    WHERE revoked_at IS NULL;
