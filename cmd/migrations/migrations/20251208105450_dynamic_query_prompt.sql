-- +goose Up
-- +goose StatementBegin
ALTER TABLE dynamic_queries
ADD COLUMN IF NOT EXISTS prompt TEXT NOT NULL;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE dynamic_queries
DROP COLUMN IF EXISTS prompt;

-- +goose StatementEnd
