CREATE TABLE document (
    id           UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    type         document_type   NOT NULL,
    resident_id  UUID NOT NULL UNIQUE REFERENCES resident(id) ON DELETE CASCADE,
    url          TEXT NOT NULL,
    status       document_status NOT NULL,
    number       TEXT NOT NULL,
    
    created_at   TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at   TIMESTAMP(3)
);
