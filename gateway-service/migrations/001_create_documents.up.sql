CREATE TABLE IF NOT EXISTS documents (
     id                UUID PRIMARY KEY,
     original_filename TEXT        NOT NULL,
     stored_filename   TEXT        NOT NULL,
     status            TEXT        NOT NULL DEFAULT 'pending',
     created_at        TIMESTAMPTZ NOT NULL DEFAULT NOW()
);