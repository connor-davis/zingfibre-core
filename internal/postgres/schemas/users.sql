CREATE TABLE
    users (
        id UUID PRIMARY KEY DEFAULT generate_uuid_v4 (),
        email TEXT NOT NULL UNIQUE,
        password TEXT NOT NULL,
        mfa_secret TEXT,
        mfa_enabled BOOLEAN DEFAULT FALSE,
        mfa_verified BOOLEAN DEFAULT FALSE,
        role role_type NOT NULL DEFAULT 'user',
        created_at TIMESTAMPTZ DEFAULT now (),
        updated_at TIMESTAMPTZ DEFAULT now ()
    );