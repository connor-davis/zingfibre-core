-- +goose Up
-- +goose StatementBegin
CREATE TYPE role_type AS ENUM ('admin', 'staff', 'user');

ALTER TABLE users
ADD COLUMN IF NOT EXISTS role role_type NOT NULL DEFAULT 'user';

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
ALTER TABLE users
DROP COLUMN IF EXISTS role;

DROP TYPE role_type;

-- +goose StatementEnd