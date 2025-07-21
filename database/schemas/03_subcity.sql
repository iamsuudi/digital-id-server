CREATE TABLE subcity (
    id            UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name          TEXT NOT NULL,
    city_id       UUID NOT NULL REFERENCES city(id) ON DELETE CASCADE,
    search_vector tsvector GENERATED ALWAYS AS (to_tsvector('english', name)) STORED,
    created_at    TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at    TIMESTAMP(3),
    UNIQUE (name, city_id)
);
