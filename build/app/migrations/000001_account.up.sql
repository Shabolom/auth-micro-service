CREATE TABLE accounts (
                          id UUID PRIMARY KEY,

                          email TEXT NOT NULL,
                          password_hash TEXT NOT NULL,

                          name TEXT NOT NULL,

                          age INTEGER NOT NULL
        CHECK (age >= 0 AND age <= 150),

                          created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
                          updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
                          deleted_at TIMESTAMPTZ NULL
);

CREATE UNIQUE INDEX ux_accounts_email
    ON accounts (email)
    WHERE deleted_at IS NULL;