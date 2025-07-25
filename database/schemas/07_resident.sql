CREATE TABLE resident (
    id               UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    email            TEXT NOT NULL UNIQUE CHECK (email ~* '^[^@]+@[^@]+\.[^@]+$'),
    first_name       TEXT NOT NULL,
    second_name      TEXT NOT NULL,
    last_name        TEXT NOT NULL,
    birth_date       TIMESTAMP(3) NOT NULL,
    gender           gender          NOT NULL,
    phone            VARCHAR(20)     NOT NULL CHECK (phone ~ '^\+?[1-9]\d{1,14}$'),
    marital_status   marital_status  NOT NULL,
    religion         religion        NOT NULL,
    ethnicity        TEXT,
    disability_status TEXT,
    education_level  TEXT,
    languages_spoken TEXT            NOT NULL,
    address_id       UUID REFERENCES address(id) ON DELETE SET NULL,
    
    created_at       TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at       TIMESTAMP(3)
);
