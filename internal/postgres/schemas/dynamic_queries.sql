CREATE TYPE order_type AS ENUM ('ASC', 'DESC');
CREATE TYPE join_type AS ENUM ('INNER', 'LEFT', 'RIGHT', 'FULL');
CREATE TYPE aggregate_type AS ENUM ('COUNT', 'SUM', 'AVG', 'MIN', 'MAX');

CREATE TABLE
    dynamic_queries (
        id UUID PRIMARY KEY DEFAULT uuid_generate_v4 (),
        query_name TEXT NOT NULL,
        is_primary BOOLEAN NOT NULL DEFAULT false,
        created_at TIMESTAMPTZ DEFAULT now (),
        updated_at TIMESTAMPTZ DEFAULT now ()
    );

CREATE TABLE
    dynamic_query_child_queries (
        id UUID PRIMARY KEY DEFAULT uuid_generate_v4 (),
        dynamic_query_id UUID NOT NULL REFERENCES dynamic_queries (id) ON DELETE SET NULL,
        child_query_id UUID REFERENCES dynamic_queries (id) ON DELETE SET NULL,
        created_at TIMESTAMPTZ DEFAULT now (),
        updated_at TIMESTAMPTZ DEFAULT now ()
    );

CREATE TABLE
    dynamic_query_columns (
        id UUID PRIMARY KEY DEFAULT uuid_generate_v4 (),
        dynamic_query_id UUID NOT NULL REFERENCES dynamic_queries (id) ON DELETE CASCADE,
        column_table TEXT NOT NULL,
        column_name TEXT NOT NULL,
        column_label TEXT NOT NULL,
        created_at TIMESTAMPTZ DEFAULT now (),
        updated_at TIMESTAMPTZ DEFAULT now ()
    );

CREATE TABLE
    dynamic_query_joins (
        id UUID PRIMARY KEY DEFAULT uuid_generate_v4 (),
        dynamic_query_id UUID NOT NULL REFERENCES dynamic_queries (id) ON DELETE CASCADE,
        join_left_table TEXT NOT NULL,
        join_left_column TEXT NOT NULL,
        join_right_table TEXT NOT NULL,
        join_right_column TEXT NOT NULL,
        join_type join_type NOT NULL DEFAULT 'INNER',
        created_at TIMESTAMPTZ DEFAULT now (),
        updated_at TIMESTAMPTZ DEFAULT now ()
    );

CREATE TABLE
    dynamic_query_clauses (
        id UUID PRIMARY KEY DEFAULT uuid_generate_v4 (),
        dynamic_query_id UUID NOT NULL REFERENCES dynamic_queries (id) ON DELETE CASCADE,
        clause_column TEXT NOT NULL,
        clause_operator TEXT NOT NULL DEFAULT '=',
        clause_value TEXT NOT NULL,
        created_at TIMESTAMPTZ DEFAULT now (),
        updated_at TIMESTAMPTZ DEFAULT now ()
    );

CREATE TABLE
    dynamic_query_orders (
        id UUID PRIMARY KEY DEFAULT uuid_generate_v4 (),
        dynamic_query_id UUID NOT NULL REFERENCES dynamic_queries (id) ON DELETE CASCADE,
        order_column TEXT NOT NULL,
        order_type order_type NOT NULL DEFAULT 'ASC',
        created_at TIMESTAMPTZ DEFAULT now (),
        updated_at TIMESTAMPTZ DEFAULT now ()
    );

CREATE TABLE
    dynamic_query_aggregates (
        id UUID PRIMARY KEY DEFAULT uuid_generate_v4 (),
        dynamic_query_column_id UUID NOT NULL REFERENCES dynamic_query_columns (id) ON DELETE CASCADE,
        aggregate_column TEXT NOT NULL,
        aggregate_type aggregate_type NOT NULL DEFAULT 'COUNT',
        created_at TIMESTAMPTZ DEFAULT now (),
        updated_at TIMESTAMPTZ DEFAULT now ()
    );