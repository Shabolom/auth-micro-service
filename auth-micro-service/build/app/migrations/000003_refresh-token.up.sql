CREATE TABLE refresh_sessions (
          id UUID PRIMARY KEY,
          user_id TEXT NOT NULL,
          refresh_token_hash TEXT NOT NULL,
          revoked_at TIMESTAMPTZ NULL,
          expires_at TIMESTAMPTZ NOT NULL,
          created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);