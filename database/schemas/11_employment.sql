CREATE TABLE employment (
    id            UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    resident_id   UUID NOT NULL UNIQUE REFERENCES resident(id) ON DELETE CASCADE,
    status        TEXT NOT NULL,
    occupation    TEXT,
    employer_name TEXT,
    work_address  TEXT,
    search_vector tsvector GENERATED ALWAYS AS (to_tsvector('english',
        status||' '||COALESCE(occupation,'')||' '||
        COALESCE(employer_name,'')||' '||COALESCE(work_address,''))) STORED,
    created_at    TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at    TIMESTAMP(3)
);
