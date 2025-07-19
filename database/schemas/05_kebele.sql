CREATE TABLE kebele (
    id            UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name          TEXT NOT NULL,
    subcity_id    UUID REFERENCES subcity(id) ON DELETE SET NULL,
    city_id       UUID NOT NULL REFERENCES city(id) ON DELETE CASCADE,
    executive_id  UUID REFERENCES actor(id) ON DELETE SET NULL,
    search_vector tsvector GENERATED ALWAYS AS (to_tsvector('english', name)) STORED,
    created_at    TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at    TIMESTAMP(3),
    UNIQUE (name, city_id)
);
