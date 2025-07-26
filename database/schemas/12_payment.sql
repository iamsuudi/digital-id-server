CREATE TABLE payment (
    id           UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    resident_id  UUID NOT NULL UNIQUE REFERENCES resident(id) ON DELETE CASCADE,
    amount       NUMERIC NOT NULL CHECK (amount >= 0),
    description  TEXT NOT NULL,
    status       payment_status NOT NULL,
    reference    TEXT NOT NULL,
    method       payment_method NOT NULL,
    
    created_at   TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at   TIMESTAMP(3)
);
