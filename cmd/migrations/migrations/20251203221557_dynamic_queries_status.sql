-- +goose Up
-- +goose StatementBegin
CREATE TYPE dynamic_query_status AS ENUM ('complete', 'in_progress', 'error');

ALTER TABLE dynamic_queries
ADD COLUMN IF NOT EXISTS status dynamic_query_status;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE dynamic_queries
DROP COLUMN IF EXISTS status;

DROP TYPE dynamic_query_status;

-- +goose StatementEnd
