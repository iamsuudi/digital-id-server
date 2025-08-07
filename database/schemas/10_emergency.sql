CREATE TABLE emergency (
    id           UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    resident_id  UUID NOT NULL UNIQUE REFERENCES resident(id) ON DELETE CASCADE,
    
    name         TEXT NOT NULL,
    relation     TEXT NOT NULL,
    phone        VARCHAR(20) NOT NULL,
    email        TEXT,

    created_at   TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at   TIMESTAMP(3)
);
