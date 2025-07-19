CREATE TABLE actor (
    id            UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    first_name    TEXT NOT NULL,
    second_name   TEXT NOT NULL,
    last_name     TEXT NOT NULL,
    email         TEXT NOT NULL UNIQUE CHECK (email ~* '^[^@]+@[^@]+\.[^@]+$'),
    phone         TEXT NOT NULL CHECK (phone ~ '^\+?[1-9]\d{1,14}$'),
    password_hash TEXT NOT NULL,
    role          TEXT NOT NULL CHECK (role IN ('SUPERADMIN','ADMIN','MANAGER','EXECUTIVE','ENCODER','CASHIER')),
    search_vector tsvector GENERATED ALWAYS AS (to_tsvector('english',
        first_name||' '||second_name||' '||last_name)) STORED,
    created_at    TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at    TIMESTAMP(3)
);
