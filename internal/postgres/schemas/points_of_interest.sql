CREATE TABLE
    points_of_interest (
        id UUID PRIMARY KEY DEFAULT generate_uuid_v4 (),
        name TEXT NOT NULL,
        key TEXT NOT NULL UNIQUE,
        created_at TIMESTAMPTZ DEFAULT now (),
        updated_at TIMESTAMPTZ DEFAULT now ()
    );