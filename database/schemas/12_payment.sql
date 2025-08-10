CREATE TABLE payment (
    id           UUID  PRIMARY KEY DEFAULT uuid_generate_v4(),
    resident_id  UUID           NOT NULL UNIQUE REFERENCES resident(id) ON DELETE CASCADE,
    amount       NUMERIC        CHECK (amount >= 0),
    description  TEXT,
    status       VARCHAR(20)    NOT NULL,
    reference    TEXT,
    method       VARCHAR(20),
    
    created_at   TIMESTAMP(3)   NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at   TIMESTAMP(3)
);
