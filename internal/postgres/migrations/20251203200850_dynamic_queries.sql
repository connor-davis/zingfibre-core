-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS
    dynamic_queries (
        id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
        name TEXT NOT NULL,
        query TEXT,
        response_id TEXT,
        created_at TIMESTAMP DEFAULT NOW(),
        updated_at TIMESTAMP DEFAULT NOW()
    );

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS dynamic_queries;

-- +goose StatementEnd
