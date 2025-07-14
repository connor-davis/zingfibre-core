-- +goose Up
-- +goose StatementBegin
CREATE TABLE
    users (
        id UUID PRIMARY KEY DEFAULT uuid_generate_v4 (),
        email TEXT NOT NULL UNIQUE,
        password TEXT NOT NULL,
        mfa_secret TEXT,
        mfa_enabled BOOLEAN DEFAULT FALSE,
        mfa_verified BOOLEAN DEFAULT FALSE,
        created_at TIMESTAMPTZ DEFAULT now (),
        updated_at TIMESTAMPTZ DEFAULT now ()
    );

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS users;

-- +goose StatementEnd