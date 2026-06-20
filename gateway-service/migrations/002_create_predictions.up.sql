create table if not exists predictions(
    id UUID PRIMARY KEY,
    document_id UUID NOT NULL REFERENCES documents(id),
    label TEXT NOT NULL,
    confidence DOUBLE PRECISION NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS predictions_index ON predictions (document_id, created_at DESC)
