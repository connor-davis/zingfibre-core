-- +goose Up
-- +goose StatementBegin
CREATE TABLE
    points_of_interest (
        id UUID PRIMARY KEY DEFAULT uuid_generate_v4 (),
        name TEXT NOT NULL,
        key TEXT NOT NULL UNIQUE,
        created_at TIMESTAMPTZ DEFAULT now (),
        updated_at TIMESTAMPTZ DEFAULT now ()
    );

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS points_of_interest;

-- +goose StatementEnd