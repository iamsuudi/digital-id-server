CREATE TABLE address (
    id            UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    house_number  TEXT NOT NULL,
    kebele_id     UUID NOT NULL REFERENCES kebele(id) ON DELETE CASCADE,
    subcity_id    UUID REFERENCES subcity(id) ON DELETE SET NULL,
    city_id       UUID NOT NULL REFERENCES city(id) ON DELETE CASCADE,
    search_vector tsvector GENERATED ALWAYS AS (to_tsvector('english', house_number)) STORED,
    created_at    TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at    TIMESTAMP(3),
    UNIQUE (house_number, kebele_id, city_id)
);
