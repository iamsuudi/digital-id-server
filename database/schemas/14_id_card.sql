CREATE TABLE idcard (
    id           UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    resident_id  UUID NOT NULL UNIQUE REFERENCES resident(id) ON DELETE CASCADE,
    number       TEXT NOT NULL UNIQUE,
    issue_date   TIMESTAMP(3) NOT NULL,
    expiry_date  TIMESTAMP(3) NOT NULL,
    issue_place  TEXT NOT NULL,
    
    created_at   TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at   TIMESTAMP(3)
);
