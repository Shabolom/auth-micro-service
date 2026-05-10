CREATE TABLE users (
   account_id TEXT,
   name TEXT,
   age INTEGER CHECK (age >= 0 AND age <= 150),
   created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
   updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
   deleted_at TIMESTAMPTZ NULL,

   CONSTRAINT chk_name_or_age
       CHECK (
           deleted_at IS NOT NULL
           OR name IS NOT NULL
           OR age IS NOT NULL
           ),

   CONSTRAINT fk_users_account
        FOREIGN KEY (account_id)
        REFERENCES accounts(id)
        ON DELETE RESTRICT
);

CREATE UNIQUE INDEX ux_users_account_id
    ON users (account_id)
    WHERE deleted_at IS NULL;