CREATE TYPE dynamic_query_status AS ENUM ('complete', 'in_progress', 'error');

CREATE TABLE IF NOT EXISTS dynamic_queries (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT NOT NULL,
    query TEXT,
    response_id TEXT,
    status dynamic_query_status NOT NULL DEFAULT 'in_progress',
    prompt TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
)
