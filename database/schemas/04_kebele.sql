CREATE TABLE kebele (
    id            UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name          TEXT NOT NULL,
    lat           NUMERIC,
    lon           NUMERIC,
    subcity_id    UUID NOT NULL REFERENCES subcity(id) ON DELETE SET NULL,
    city_id       UUID NOT NULL REFERENCES city(id) ON DELETE CASCADE,
    
    created_at    TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at    TIMESTAMP(3),
    UNIQUE (name, subcity_id, city_id)
);
